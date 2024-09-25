package taskbus

import (
	"database/sql"
	"time"
)

// Task represents a task in the system.
type Task struct {
	ID          int          `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	CreatedAt   time.Time    `json:"created_at"`
	FinishedAt  sql.NullTime `json:"finished_at"`
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
		FinishedAt:  finishedAt,
	}
}

// UpdateTask represents a task with updates to be applied.
type UpdateTask struct {
	Title       string
	Description string
}
