package taskbus

import (
	"database/sql"
	"time"
)

// Task represents a task in the system.
type Task struct {
	ID          int           `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	CreatedAt   time.Time     `json:"created_at"`
	FinishedAt  sql.NullTime  `json:"finished_at"`
	CreatedBy   int           `json:"created_by"`
	AssignedTo  sql.NullInt32 `json:"assigned_to"`
}

// NewTask represents a new task to be created.
type NewTask struct {
	Title       string
	Description string
	CreatedBy   int
	AssignedTo  sql.NullInt32
}

// UpdateTask represents a task with updates to be applied.
type UpdateTask struct {
	Title       string
	Description string
	AssignedTo  sql.NullInt32
}
