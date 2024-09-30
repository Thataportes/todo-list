package taskbus

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Task represents a task in the system.
type Task struct {
	ID             int           `json:"id"`
	Title          string        `json:"title"`
	Description    string        `json:"description"`
	CreatedAt      time.Time     `json:"created_at"`
	FinishedAt     sql.NullTime  `json:"finished_at"`
	CreatedBy      int           `json:"created_by"`
	CreatedByName  string        `json:"created_by_name"`
	AssignedTo     sql.NullInt32 `json:"assigned_to"`
	AssignedToName string        `json:"assigned_to_name"`
}

// NewTask represents a new task to be created.
type NewTask struct {
	Title       string
	Description string
	CreatedBy   int
	AssignedTo  sql.NullInt32
}

// toBusinessTask converts data into a business Task model.
func toBusinessTask(id int, title, description string, createdAt, finishedAt sql.NullTime, createdBy int, createdByName string, assignedTo sql.NullInt32, assignedToName string) Task {

	return Task{
		ID:             id,
		Title:          title,
		Description:    description,
		CreatedAt:      createdAt.Time,
		FinishedAt:     finishedAt,
		CreatedBy:      createdBy,
		CreatedByName:  createdByName,
		AssignedTo:     assignedTo,
		AssignedToName: assignedToName,
	}
}

// UpdateTask represents a task with updates to be applied.
type UpdateTask struct {
	Title       string
	Description string
	AssignedTo  sql.NullInt32
}

// getUserNameByID retrieves the name of a user by their ID from the database.
func (s *Business) getUserNameByID(ctx context.Context, userID int) (string, error) {
	var name string
	query := "SELECT name FROM users WHERE id = ?"
	err := s.db.QueryRowContext(ctx, query, userID).Scan(&name)
	if err != nil {
		return "", fmt.Errorf("failed to get user name for ID %d: %w", userID, err)
	}
	return name, nil
}
