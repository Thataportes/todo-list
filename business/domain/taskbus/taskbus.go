package taskbus

import (
	"context"
	"database/sql"
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
	createdAt := time.Now()
	query := "INSERT INTO tasks (title, description, created_at) VALUES (?, ?, ?)"
	result, err := s.db.ExecContext(ctx, query, nt.Title, nt.Description, createdAt)
	if err != nil {
		return Task{}, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return Task{}, err
	}

	t := Task{
		ID:          int(lastInsertID),
		Title:       nt.Title,
		Description: nt.Description,
		CreatedAt:   createdAt,
	}

	return t, nil
}

// Query retrieves all tasks from the database.
func (s *Business) Query(ctx context.Context) ([]Task, error) {
	query := "SELECT id, title, description, created_at, finished_at FROM tasks"
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var finishedAt sql.NullTime
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.CreatedAt, &finishedAt)
		if err != nil {
			return nil, err
		}
		if finishedAt.Valid {
			task.FinishedAt = &finishedAt.Time
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
	query := "SELECT id, title, description, created_at, finished_at FROM tasks WHERE id =?"
	row := s.db.QueryRowContext(ctx, query, id)

	var task Task
	var finishedAt sql.NullTime
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.CreatedAt, &finishedAt)
	if err != nil {
		return Task{}, err
	}
	if finishedAt.Valid {
		task.FinishedAt = &finishedAt.Time
	}
	return task, nil
}

// Update modifies task information in the database and returns the updated task.ções de uma tarefa no banco de dados e retorna a tarefa atualizada.
func (s *Business) Update(ctx context.Context, ut UpdateTask) (Task, error) {
	query := "UPDATE tasks SET title=?, description=? WHERE id=? "
	_, err := s.db.ExecContext(ctx, query, ut.Title, ut.Description, ut.ID)
	if err != nil {
		return Task{}, err
	}

	updatedTask, err := s.QueryByID(ctx, ut.ID)
	if err != nil {
		return Task{}, err
	}

	return updatedTask, nil
}

// Delete removes a task from the database by its ID.
func (s *Business) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM tasks WHERE id=?"
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

// Finish updates the finishedAt timestamp for a task.
func (s *Business) Finish(ctx context.Context, id int) error {
	finishedAt := time.Now()
	query := "UPDATE tasks SET finished_at = ? WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, finishedAt, id)
	return err
}
