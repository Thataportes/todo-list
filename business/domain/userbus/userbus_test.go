package userbus_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"TODO-list/business/domain/userbus"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var (
	db       *sql.DB
	mock     sqlmock.Sqlmock
	business *userbus.Business
)

func setupMockDB(t *testing.T) {
	var err error
	db, mock, err = sqlmock.New()
	assert.NoError(t, err)
	business = userbus.NewBusiness(db)
}

func mockUserRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "name", "email", "active", "created_at", "updated_at"}).
		AddRow(1, "User 1", "user1@example.com", true, time.Now(), time.Now()).
		AddRow(2, "User 2", "user2@example.com", true, time.Now(), time.Now())
}

func assertMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate(t *testing.T) {
	setupMockDB(t)
	defer db.Close()

	mock.ExpectExec("INSERT INTO users").
		WithArgs("New User", "newuser@example.com", true, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	newUser := userbus.NewUser{Name: "New User", Email: "newuser@example.com"}
	user, err := business.Create(ctx, newUser)

	assert.NoError(t, err)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "New User", user.Name)
	assert.Equal(t, "newuser@example.com", user.Email)
	assert.True(t, user.Active)
	assert.True(t, user.CreatedAt.Valid)
	assert.True(t, user.UpdatedAt.Valid)
	assert.True(t, user.CreatedAt.Time.After(time.Now().Add(-time.Hour)))
	assertMockExpectations(t, mock)
}

func TestQuery(t *testing.T) {
	setupMockDB(t)
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, email, active, created_at, updated_at FROM users").
		WillReturnRows(mockUserRows())

	ctx := context.Background()
	users, err := business.Query(ctx)

	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "User 1", users[0].Name)
	assert.Equal(t, "user1@example.com", users[0].Email)
	assert.True(t, users[0].Active)
	assert.True(t, users[0].CreatedAt.Valid)
	assert.True(t, users[0].UpdatedAt.Valid)
	assert.True(t, users[0].CreatedAt.Time.After(time.Now().Add(-time.Hour)))
	assertMockExpectations(t, mock)
}

func TestQueryByID(t *testing.T) {
	setupMockDB(t)
	defer db.Close()

	row := sqlmock.NewRows([]string{"id", "name", "email", "active", "created_at", "updated_at"}).
		AddRow(1, "User 1", "user1@example.com", true, time.Now(), time.Now())

	mock.ExpectQuery("SELECT id, name, email, active, created_at, updated_at FROM users WHERE id = ?").
		WithArgs(1).
		WillReturnRows(row)

	ctx := context.Background()
	user, err := business.QueryById(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, "User 1", user.Name)
	assert.Equal(t, "user1@example.com", user.Email)
	assert.True(t, user.Active)
	assert.True(t, user.CreatedAt.Valid)
	assert.True(t, user.UpdatedAt.Valid)
	assert.True(t, user.CreatedAt.Time.After(time.Now().Add(-time.Hour)))
	assertMockExpectations(t, mock)
}

func TestQueryByEmail(t *testing.T) {
	setupMockDB(t)
	defer db.Close()

	row := sqlmock.NewRows([]string{"id", "name", "email", "active", "created_at", "updated_at"}).
		AddRow(1, "User 1", "user1@example.com", true, time.Now(), time.Now())

	mock.ExpectQuery("SELECT id, name, email, active, created_at, updated_at FROM users WHERE email = ?").
		WithArgs("user1@example.com").
		WillReturnRows(row)

	ctx := context.Background()
	user, err := business.QueryByEmail(ctx, "user1@example.com")

	assert.NoError(t, err)
	assert.Equal(t, "User 1", user.Name)
	assert.Equal(t, "user1@example.com", user.Email)
	assert.True(t, user.Active)
	assert.True(t, user.CreatedAt.Valid)
	assert.True(t, user.UpdatedAt.Valid)
	assert.True(t, user.CreatedAt.Time.After(time.Now().Add(-time.Hour)))
	assertMockExpectations(t, mock)
}

func TestUpdate(t *testing.T) {
	setupMockDB(t)
	defer db.Close()

	mock.ExpectExec("^UPDATE users SET name = \\?, email = \\?, updated_at = \\? WHERE id = \\?$").
		WithArgs("Updated Name", "updated@example.com", sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := context.Background()
	updateUser := userbus.UpdateUser{Name: "Updated Name", Email: "updated@example.com"}
	err := business.Update(ctx, 1, updateUser)

	assert.NoError(t, err)
	assertMockExpectations(t, mock)
}

func TestDeactivate(t *testing.T) {
	setupMockDB(t)
	defer db.Close()

	mock.ExpectExec("^UPDATE users SET active = false, updated_at = \\? WHERE id = \\?$").
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := context.Background()
	err := business.Deactivate(ctx, 1)

	assert.NoError(t, err)
	assertMockExpectations(t, mock)
}
