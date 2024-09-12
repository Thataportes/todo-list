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
	query := "INSERT INTO task (title, description, created_at, finished_at) VALUES (?, ?, ?, ?)"
	result, err := s.db.ExecContext(ctx, query, nt.Title, nt.Description, createdAt, finishedAt)
	if err != nil {
		return Task{}, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return Task{}, err
	}

	return toBusinessTask(int(lastInsertID), nt.Title, nt.Description, createdAt, finishedAt), nil
}

// Query retrieves all tasks from the database.
func (s *Business) Query(ctx context.Context) ([]Task, error) {
	query := "SELECT id, title, description, created_at, finished_at FROM task"
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var id int
		var title, description string
		var createdAt, finishedAt sql.NullTime

		err := rows.Scan(&id, &title, &description, &createdAt, &finishedAt)
		if err != nil {
			return nil, err
		}

		task := toBusinessTask(id, title, description, createdAt, finishedAt)
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// QueryByID retrieves a task by its ID.
func (s *Business) QueryByID(ctx context.Context, id int) (Task, error) {
	query := "SELECT id, title, description, created_at, finished_at FROM task WHERE id = ?"
	row := s.db.QueryRowContext(ctx, query, id)

	var ID int
	var title, description string
	var createdAt, finishedAt sql.NullTime

	err := row.Scan(&ID, &title, &description, &createdAt, &finishedAt)
	if err != nil {
		return Task{}, err
	}

	return toBusinessTask(ID, title, description, createdAt, finishedAt), nil
}

// Update modifies task information in the database and returns the updated task.
func (s *Business) Update(ctx context.Context, ut UpdateTask) error {
	query := "UPDATE task SET title = ?, description = ? WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, ut.Title, ut.Description, ut.ID)
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
