package userbus

import (
	"database/sql"
)

// User represents a user entity in the business layer.
type User struct {
	ID            int          `json:"id"`
	Name          string       `json:"name"`
	Email         string       `json:"email"`
	Status        bool         `json:"status"`
	CreatedAt     sql.NullTime `json:"created_at"`
	LastUpdatedAt sql.NullTime `json:"last_updated_at"`
}

// NewUser represents the input data required to create a new user.
type NewUser struct {
	Name  string
	Email string
}

// UpdateUser represents the input data required to update an existing user.
type UpdateUser struct {
	Name   string
	Email  string
	Status bool
}

// toBusinessUser converts the input data into a User struct in the business layer.
func toBusinessUser(id int, name, email string, status bool, createdAt, lastUpdatedAt sql.NullTime) User {
	return User{
		ID:            id,
		Name:          name,
		Email:         email,
		Status:        status,
		CreatedAt:     createdAt,
		LastUpdatedAt: lastUpdatedAt,
	}
}
