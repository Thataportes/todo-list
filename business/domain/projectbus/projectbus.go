package projectbus

import (
	"TODO-list/business/domain/userbus"
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Business handles business logic and persistence for project-related operations.
type Business struct {
	db      *sql.DB
	userBus *userbus.Business
}

// NewBusiness creates a new instance of Business with the provided database connection and user operations.
func NewBusiness(db *sql.DB, userBus *userbus.Business) *Business {
	return &Business{
		db:      db,
		userBus: userBus,
	}
}

// Create inserts a new project into the database and returns the created project.
func (s *Business) Create(ctx context.Context, np NewProject) (Project, error) {
	creator, err := s.userBus.QueryById(ctx, np.CreatedBy)
	if err != nil {
		return Project{}, fmt.Errorf("failed to retrieve creator user with ID %d: %v", np.CreatedBy, err)
	}
	if !creator.Active {
		return Project{}, fmt.Errorf("creator user with ID %d is not active", np.CreatedBy)
	}

	createdAt := time.Now()
	query := "INSERT INTO project (name, active, created_at, created_by) VALUES (?, ?, ?, ?) "
	result, err := s.db.ExecContext(ctx, query, np.Name, true, createdAt, np.CreatedBy)
	if err != nil {
		return Project{}, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return Project{}, err
	}

	return Project{
		ID:        int(lastInsertID),
		Name:      np.Name,
		Active:    true,
		CreatedAt: createdAt,
		CreatedBy: np.CreatedBy,
	}, nil
}

// Query retrieves all projects from the database.
func (s *Business) Query(ctx context.Context) ([]Project, error) {
	query := "SELECT id, name, active, created_at, created_by FROM project"
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var project Project
		err := rows.Scan(&project.ID, &project.Name, &project.Active, &project.CreatedAt, &project.CreatedBy)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

// QueryById retrieves a specific project by its ID from the database.
func (s *Business) QueryById(ctx context.Context, id int) (Project, error) {
	query := "SELECT id, name, active, created_at, created_by FROM project WHERE id = ?"
	row := s.db.QueryRowContext(ctx, query, id)

	var project Project
	err := row.Scan(&project.ID, &project.Name, &project.Active, &project.CreatedAt, &project.CreatedBy)
	if err != nil {
		return Project{}, err
	}

	return project, nil
}

// Update modifies an existing project's information in the database.
func (s *Business) Update(ctx context.Context, id int, up UpdateProject) error {
	query := "UPDATE project SET name = ? WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, up.Name, id)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a project from the database by its ID.
func (s *Business) Delete(ctx context.Context, id int) error {
	checkQuery := "SELECT EXISTS(SELECT 1 FROM task WHERE project_id = ?)"
	var hasTasks bool
	err := s.db.QueryRowContext(ctx, checkQuery, id).Scan(&hasTasks)
	if err != nil {
		return fmt.Errorf("failed to check if project has associated tasks with ID %d: %w", id, err)
	}
	if hasTasks {
		return s.Deactivate(ctx, id)
	}
	Query := "DELETE FROM project WHERE id = ?"
	_, err = s.db.ExecContext(ctx, Query, id)
	if err != nil {
		return fmt.Errorf("failed to delete project with ID %d: %w", id, err)
	}

	return nil
}

// Deactivate sets a project's status to inactive (false) in the database.
func (s *Business) Deactivate(ctx context.Context, id int) error {
	query := "UPDATE project SET active = false WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to deactivate project with ID %d: %v", id, err)
	}
	return nil
}
