package projectbus

import "time"

// Project represents a project entity in the system.
type Project struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy int       `json:"created_by"`
}

// NewProject represents the data required to create a new project.
type NewProject struct {
	Name      string
	CreatedBy int
}

// UpdateProject represents the data required to update an existing project.
type UpdateProject struct {
	Name string
}
