package taskbus

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Business handles business logic and persistence of tasks.
type Business struct {
	db *sql.DB
}

// NewBusiness creates a new instance of Business.
func NewBusiness(db *sql.DB) *Business {
	return &Business{db: db}
}

// Create adds a new task to the database and returns the created task.
func (s *Business) Create(ctx context.Context, nt NewTask) (Task, error) {
	createdAt := sql.NullTime{Time: time.Now(), Valid: true}
	finishedAt := sql.NullTime{Valid: false}

	query := "INSERT INTO task (title, description, created_by, assigned_to, created_at, finished_at) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := s.db.ExecContext(ctx, query, nt.Title, nt.Description, nt.CreatedBy, nt.AssignedTo, createdAt, finishedAt)
	if err != nil {
		return Task{}, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return Task{}, err
	}

	createdByName, err := s.getUserNameByID(ctx, nt.CreatedBy)
	if err != nil {
		return Task{}, err
	}

	assignedToName := ""
	if nt.AssignedTo.Valid {
		assignedToName, err = s.getUserNameByID(ctx, int(nt.AssignedTo.Int32))
		if err != nil {
			return Task{}, err
		}
	}

	return toBusinessTask(int(lastInsertID), nt.Title, nt.Description, createdAt, finishedAt, nt.CreatedBy, createdByName, nt.AssignedTo, assignedToName), nil
}

// Query retrieves all tasks from the database.
func (s *Business) Query(ctx context.Context) ([]Task, error) {
	query := `SELECT 
				t.id, t.title, t.description, t.created_at, t.finished_at, 
				t.created_by, cu.name AS created_by_name, 
				t.assigned_to, au.name AS assigned_to_name 
			  FROM 
				task t 
			  LEFT JOIN 
				users cu ON t.created_by = cu.id 
			  LEFT JOIN 
				users au ON t.assigned_to = au.id`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var id, createdBy int
		var assignedTo sql.NullInt32
		var title, description string
		var createdByName, assignedToName sql.NullString
		var createdAt, finishedAt sql.NullTime

		err := rows.Scan(&id, &title, &description, &createdAt, &finishedAt, &createdBy, &createdByName, &assignedTo, &assignedToName)
		if err != nil {
			return nil, err
		}

		task := toBusinessTask(
			id,
			title,
			description,
			createdAt,
			finishedAt,
			createdBy,
			createdByName.String,
			assignedTo,
			assignedToName.String,
		)
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// QueryByID retrieves a task by its ID.
func (s *Business) QueryByID(ctx context.Context, id int) (Task, error) {
	query := `SELECT 
				t.id, t.title, t.description, t.created_at, t.finished_at, 
				t.created_by, cu.name AS created_by_name, 
				t.assigned_to, au.name AS assigned_to_name 
			  FROM task t 
			  LEFT JOIN users cu ON t.created_by = cu.id 
			  LEFT JOIN users au ON t.assigned_to = au.id 
			  WHERE t.id = ?`
	row := s.db.QueryRowContext(ctx, query, id)

	var busTask Task
	var createdByName, assignedToName sql.NullString
	err := row.Scan(&busTask.ID, &busTask.Title, &busTask.Description, &busTask.CreatedAt, &busTask.FinishedAt, &busTask.CreatedBy, &createdByName, &busTask.AssignedTo, &assignedToName)
	if err != nil {
		return Task{}, err
	}

	busTask.CreatedByName = createdByName.String
	busTask.AssignedToName = assignedToName.String

	return busTask, nil
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
