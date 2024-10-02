package userbus

import (
	"context"
	"database/sql"
	"time"
)

// Business handles business logic and persistence of user-related operations.
type Business struct {
	db *sql.DB
}

// NewBusiness creates a new instance of Business with the provided database connection.
func NewBusiness(db *sql.DB) *Business {
	return &Business{db: db}
}

// Create inserts a new user into the database and returns the created user.
func (s *Business) Create(ctx context.Context, nu NewUser) (User, error) {
	createdAt := sql.NullTime{Time: time.Now(), Valid: true}
	lastUpdatedAt := sql.NullTime{Time: createdAt.Time, Valid: true}
	query := "INSERT INTO users (name, email, status, created_at, last_updated_at) VALUES (?, ?, ?, ?, ?)"
	result, err := s.db.ExecContext(ctx, query, nu.Name, nu.Email, true, createdAt, lastUpdatedAt)
	if err != nil {
		return User{}, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return User{}, err
	}

	return User{
		ID:            int(lastInsertID),
		Name:          nu.Name,
		Email:         nu.Email,
		Status:        true,
		CreatedAt:     createdAt,
		LastUpdatedAt: lastUpdatedAt,
	}, nil
}

// Query retrieves all users from the database and returns them as a slice of User structs.
func (s *Business) Query(ctx context.Context) ([]User, error) {
	query := "SELECT id, name, email, status, created_at, last_updated_at FROM users"
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var busUser User

		err := rows.Scan(&busUser.ID, &busUser.Name, &busUser.Email, &busUser.Status, &busUser.CreatedAt, &busUser.LastUpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, busUser)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// QueryById retrieves a specific user by their ID from the database.
func (s *Business) QueryById(ctx context.Context, id int) (User, error) {
	query := "SELECT id, name, email, status, created_at, last_updated_at FROM users WHERE id = ?"
	row := s.db.QueryRowContext(ctx, query, id)

	var busUser User
	err := row.Scan(&busUser.ID, &busUser.Name, &busUser.Email, &busUser.Status, &busUser.CreatedAt, &busUser.LastUpdatedAt)
	if err != nil {
		return User{}, err
	}

	return busUser, nil
}

// QueryByEmail retrieves a specific user by their email from the database.
func (s *Business) QueryByEmail(ctx context.Context, email string) (User, error) {
	query := "SELECT id, name, email, status, created_at, last_updated_at FROM users WHERE email = ? "
	row := s.db.QueryRowContext(ctx, query, email)

	var busUser User
	err := row.Scan(&busUser.ID, &busUser.Name, &busUser.Email, &busUser.Status, &busUser.CreatedAt, &busUser.LastUpdatedAt)
	if err != nil {
		return User{}, err
	}

	return busUser, nil
}

// Update modifies an existing user's information in the database.
func (s *Business) Update(ctx context.Context, id int, uu UpdateUser) error {
	lastUpdatedAt := sql.NullTime{Time: time.Now(), Valid: true}
	query := "UPDATE users SET name = ?, email = ?, status = ?, last_updated_at = ? WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, uu.Name, uu.Email, uu.Status, lastUpdatedAt, id)
	if err != nil {
		return err
	}

	return nil
}

// Deactivate sets a user's status to inactive (false) in the database.
func (s *Business) Deactivate(ctx context.Context, id int) error {
	lastUpdatedAt := sql.NullTime{Time: time.Now(), Valid: true}
	query := "UPDATE users SET status = false, last_updated_at = ? WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, lastUpdatedAt, id)
	return err
}

// IsUserActive checks if a user is active in the system based on their user ID.
func (s *Business) IsUserActive(ctx context.Context, userID int) (bool, error) {
	var status bool
	query := "SELECT status FROM users WHERE id = ?"
	err := s.db.QueryRowContext(ctx, query, userID).Scan(&status)
	if err != nil {
		return false, err
	}
	return status, nil
}
