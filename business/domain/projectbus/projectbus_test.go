package projectbus_test

import (
	"TODO-list/business/domain/projectbus"
	"TODO-list/business/domain/userbus"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var (
	db       *sql.DB
	mock     sqlmock.Sqlmock
	business *projectbus.Business
)

func setupMockDB(t *testing.T) {
	var err error
	db, mock, err = sqlmock.New()
	assert.NoError(t, err)

	userBus := userbus.NewBusiness(db)
	business = projectbus.NewBusiness(db, userBus)
}

func mockProjectRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "name", "active", "created_at", "created_by"}).
		AddRow(1, "Project 1", true, time.Now(), 1).
		AddRow(2, "Project 2", false, time.Now(), 1)
}

func assertMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestCreate(t *testing.T) {
	setupMockDB(t)
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, email, active, created_at, updated_at FROM users WHERE id = ?").
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "active", "created_at", "updated_at"}).
			AddRow(1, "Creator Name", "creator@example.com", true, time.Now(), time.Now()))

	mock.ExpectExec("INSERT INTO project \\(name, active, created_at, created_by\\) VALUES \\(\\?, \\?, \\?, \\?\\)").
		WithArgs("New Project", true, sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	newProject := projectbus.NewProject{
		Name:      "New Project",
		CreatedBy: 1,
	}
	project, err := business.Create(ctx, newProject)

	assert.NoError(t, err)
	assert.Equal(t, 1, project.ID)
	assert.Equal(t, "New Project", project.Name)
	assert.True(t, project.Active)
	assert.NotEmpty(t, project.CreatedAt)
	assert.Equal(t, int64(1), int64(project.ID))

	assertMockExpectations(t, mock)
}

func TestQuery(t *testing.T) {
	setupMockDB(t)
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, active, created_at, created_by FROM project").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "active", "created_at", "created_by"}).
			AddRow(1, "Project 1", true, time.Now(), 1).
			AddRow(2, "Project 2", true, time.Now(), 1))

	ctx := context.Background()
	projects, err := business.Query(ctx)

	assert.NoError(t, err)
	assert.Len(t, projects, 2)
	assert.Equal(t, "Project 1", projects[0].Name)
	assert.True(t, projects[0].Active)
	assert.NotEmpty(t, projects[0].CreatedAt)
	assert.True(t, projects[0].CreatedAt.After(time.Now().Add(-time.Hour)))
	assert.Equal(t, 1, projects[0].CreatedBy)
	assertMockExpectations(t, mock)
}

func TestQueryById(t *testing.T) {
	setupMockDB(t)
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, active, created_at, created_by FROM project WHERE id = ?").
		WithArgs(3).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "active", "created_at", "created_by"}).
			AddRow(3, "Project 3", true, time.Now(), 3))

	ctx := context.Background()
	project, err := business.QueryById(ctx, 3)

	assert.NoError(t, err)
	assert.Equal(t, "Project 3", project.Name)
	assert.True(t, project.Active)
	assert.NotEmpty(t, project.CreatedAt)

	assertMockExpectations(t, mock)
}

func TestUpdate(t *testing.T) {
	setupMockDB(t)
	defer db.Close()

	mock.ExpectExec("^UPDATE project SET name = \\? WHERE id = \\?$").
		WithArgs("Updated Project", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := context.Background()

	update := projectbus.UpdateProject{
		Name: "Updated Project",
	}
	err := business.Update(ctx, 1, update)

	assert.NoError(t, err)
	assertMockExpectations(t, mock)
}

func TestDeactivate(t *testing.T) {
	setupMockDB(t)
	defer db.Close()

	mock.ExpectExec("^UPDATE project SET active = false WHERE id = \\?$").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := context.Background()
	err := business.Deactivate(ctx, 1)

	assert.NoError(t, err)
	assertMockExpectations(t, mock)
}

func TestDelete(t *testing.T) {
	setupMockDB(t)
	defer db.Close()

	mock.ExpectExec("^DELETE FROM project WHERE id = \\?$").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := context.Background()
	err := business.Delete(ctx, 1)

	assert.NoError(t, err)
	assertMockExpectations(t, mock)
}
