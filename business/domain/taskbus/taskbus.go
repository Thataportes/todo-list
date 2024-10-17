package taskbus

import (
	"TODO-list/business/domain/projectbus"
	"TODO-list/business/domain/userbus"
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Business handles business logic and persistence of tasks.
type Business struct {
	db         *sql.DB
	userBus    *userbus.Business
	projectBus *projectbus.Business
}

// NewBusiness initializes a new instance of Business with the given database and user business logic.
func NewBusiness(db *sql.DB, userBus *userbus.Business, projectBus *projectbus.Business) *Business {
	return &Business{
		db:         db,
		userBus:    userBus,
		projectBus: projectBus,
	}
}

// Create adds a new task to the database after validating the creator and assigned users,
func (s *Business) Create(ctx context.Context, nt NewTask) (Task, error) {
	project, err := s.projectBus.QueryById(ctx, nt.ProjectID)
	if err != nil {
		return Task{}, fmt.Errorf("project with ID %d does not exist: %w", nt.ProjectID, err)
	}
	if !project.Active {
		return Task{}, fmt.Errorf("project with ID %d is not active", nt.ProjectID)
	}

	creator, err := s.userBus.QueryById(ctx, nt.CreatedBy)
	if err != nil {
		return Task{}, fmt.Errorf("failed to retrieve creator user with ID %d: %v", nt.CreatedBy, err)
	}
	if !creator.Active {
		return Task{}, fmt.Errorf("creator user with ID %d is not active", nt.CreatedBy)
	}

	if nt.AssignedTo.Valid {
		user, err := s.userBus.QueryById(ctx, int(nt.AssignedTo.Int32))
		if err != nil {
			return Task{}, fmt.Errorf("failed to retrieve assigned user with ID %d: %v", nt.AssignedTo.Int32, err)
		}
		if !user.Active {
			return Task{}, fmt.Errorf("assigned user with ID %d is not active", nt.AssignedTo.Int32)
		}
	}

	createdAt := sql.NullTime{Time: time.Now(), Valid: true}
	finishedAt := sql.NullTime{Valid: false}

	query := "INSERT INTO task (title, description, created_by, assigned_to, project_id, created_at, finished_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, err := s.db.ExecContext(ctx, query, nt.Title, nt.Description, nt.CreatedBy, nt.AssignedTo, nt.ProjectID, createdAt, finishedAt)
	if err != nil {
		return Task{}, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return Task{}, err
	}

	return Task{
		ID:          int(lastInsertID),
		Title:       nt.Title,
		Description: nt.Description,
		ProjectID:   nt.ProjectID,
		CreatedAt:   createdAt.Time,
		FinishedAt:  finishedAt,
		CreatedBy:   nt.CreatedBy,
		AssignedTo:  nt.AssignedTo,
	}, nil
}

// Query retrieves all tasks from the database.
func (s *Business) Query(ctx context.Context) ([]Task, error) {
	query := "SELECT id, title, description, created_at, finished_at, created_by, assigned_to, project_id FROM task"

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.CreatedAt, &task.FinishedAt, &task.CreatedBy, &task.AssignedTo, &task.ProjectID)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// QueryByID retrieves a task by its ID.
func (s *Business) QueryByID(ctx context.Context, id int) (Task, error) {
	query := "SELECT id, title, description, project_id, created_at, finished_at, created_by, assigned_to FROM task WHERE id = ?"
	row := s.db.QueryRowContext(ctx, query, id)

	var task Task
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.ProjectID, &task.CreatedAt, &task.FinishedAt, &task.CreatedBy, &task.AssignedTo)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

// Update modifies task information in the database and returns the updated task.
func (s *Business) Update(ctx context.Context, id int, ut UpdateTask) error {
	query := "UPDATE task SET title = ?, description = ?, assigned_to = ? WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, ut.Title, ut.Description, ut.AssignedTo, id)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a task from the database by its ID.
func (s *Business) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM task WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

// Finish updates the finishedAt timestamp for a task.
func (s *Business) Finish(ctx context.Context, id int) error {
	now := time.Now()

	finishedAt := sql.NullTime{
		Time:  time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.Local),
		Valid: true,
	}

	query := "UPDATE task SET finished_at = ? WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, finishedAt, id)
	if err != nil {
		return fmt.Errorf("failed to finish task with ID %d: %v", id, err)
	}

	return nil
}
