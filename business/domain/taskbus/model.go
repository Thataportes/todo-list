package taskbus

import "time"

// Representa uma tarefa no sistema.
type Task struct {
	ID          int
	Title       string
	Description string
	FinishedAt  time.Time
	CreatedAt   time.Time
}

type NewTask struct {
	Title       string
	Description string
}

type UpdateTask struct {
	ID          int
	Title       string
	Description string
}
