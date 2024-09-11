package taskbus

import (
	"database/sql"
	"time"
)

// Task represents a task in the system.
type Task struct {
	ID          int
	Title       string
	Description string
	FinishedAt  time.Time
	CreatedAt   time.Time
}

// NewTask represents a new task to be created.
type NewTask struct {
	Title       string
	Description string
}

// toBusinessTask converts data into a business Task model.
func toBusinessTask(id int, title, description string, createdAt, finishedAt sql.NullTime) Task {
	return Task{
		ID:          id,
		Title:       title,
		Description: description,
		CreatedAt:   createdAt.Time,
		FinishedAt:  finishedAt.Time,
	}
}

// UpdateTask represents a task with updates to be applied.
type UpdateTask struct {
	ID          int
	Title       string
	Description string
}
