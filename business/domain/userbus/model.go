package userbus

import (
	"database/sql"
)

// User represents a user entity in the business layer.
type User struct {
	ID        int          `json:"id"`
	Name      string       `json:"name"`
	Email     string       `json:"email"`
	Active    bool         `json:"active"`
	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}

// NewUser represents the input data required to create a new user.
type NewUser struct {
	Name  string
	Email string
}

// UpdateUser represents the input data required to update an existing user.
type UpdateUser struct {
	Name  string
	Email string
}
