package usersapp

import (
	"TODO-list/userbusiness/domain/userbus"
	"encoding/json"
	"time"
)

// NewUser represents the input data required to create a new user.
type NewUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Decode decodes a JSON byte slice into a NewUser struct.
func (nu *NewUser) Decode(data []byte) error {
	return json.Unmarshal(data, &nu)
}

// toBusNewUser converts a NewUser struct from the application layer to the business layer representation.
func toBusNewUser(nu NewUser) userbus.NewUser {
	return userbus.NewUser{
		Name:  nu.Name,
		Email: nu.Email,
	}
}

// User represents a user entity in the system with all relevant fields.
type User struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Status        bool      `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
}

// Encode encodes the User struct into a JSON byte slice.
func (u User) Encode() ([]byte, string, error) {
	data, err := json.Marshal(u)
	return data, "application/json", err
}

// toAppUser converts a User struct from the business layer to the application layer representation.
func toAppUser(userBus userbus.User) User {
	return User{
		ID:            userBus.ID,
		Name:          userBus.Name,
		Email:         userBus.Email,
		Status:        userBus.Status,
		CreatedAt:     userBus.CreatedAt.Time,
		LastUpdatedAt: userBus.LastUpdatedAt.Time,
	}
}

// Users represents a collection of User entities.
type Users []User

// Encode encodes the Users slice into a JSON byte slice.
func (us Users) Encode() ([]byte, string, error) {
	data, err := json.Marshal(us)
	return data, "application/json", err
}

// toAppUsers converts a slice of User structs from the business layer to the application layer representation.
func toAppUsers(usersBus []userbus.User) Users {
	usersApp := make(Users, len(usersBus))
	for i, userBus := range usersBus {
		usersApp[i] = toAppUser(userBus)
	}
	return usersApp
}

// UpdateUser represents the input data required to update an existing user.
type UpdateUser struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status bool   `json:"status"`
}

// Decode decodes a JSON byte slice into an UpdateUser struct.
func (uu *UpdateUser) Decode(data []byte) error {
	return json.Unmarshal(data, &uu)
}

// toBusUpdateUser converts an UpdateUser struct from the application layer to the business layer representation.
func toBusUpdateUser(uu UpdateUser) userbus.UpdateUser {
	return userbus.UpdateUser{
		Name:   uu.Name,
		Email:  uu.Email,
		Status: uu.Status,
	}
}
