package taskbus

import "time"

// Task represents a task in the system.
type Task struct {
	ID          int
	Title       string
	Description string
	FinishedAt  *time.Time
	CreatedAt   time.Time
}

// NewTask represents a new task to be created.
type NewTask struct {
	Title       string
	Description string
}

// UpdateTask represents a task with updates to be applied.
type UpdateTask struct {
	ID          int
	Title       string
	Description string
}
