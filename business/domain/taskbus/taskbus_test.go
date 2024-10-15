package taskbus_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"TODO-list/business/domain/projectbus"
	"TODO-list/business/domain/taskbus"
	"TODO-list/business/domain/userbus"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var (
	db       *sql.DB
	mock     sqlmock.Sqlmock
	business *taskbus.Business
)

func setupMockDB(t *testing.T) {
	var err error
	db, mock, err = sqlmock.New()
	assert.NoError(t, err)

	userBus := userbus.NewBusiness(db)
	projectBus := projectbus.NewBusiness(db, userBus)
	business = taskbus.NewBusiness(db, userBus, projectBus)
}

func mockTaskRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "title", "description", "created_at", "finished_at", "created_by", "assigned_to", "project_id"}).
		AddRow(1, "Task 1", "Description 1", time.Now(), sql.NullTime{Valid: false}, 1, sql.NullInt32{}, 3).
		AddRow(2, "Task 2", "Description 2", time.Now(), sql.NullTime{Valid: false}, 1, sql.NullInt32{}, 3)
}

func assertMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate(t *testing.T) {
	setupMockDB(t)
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, active, created_at, created_by FROM project WHERE id = ?").
		WithArgs(3).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "active", "created_at", "created_by"}).
			AddRow(3, "Project Name", true, time.Now(), 1))

	mock.ExpectQuery("SELECT id, name, email, active, created_at, updated_at FROM users WHERE id = ?").
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "active", "created_at", "updated_at"}).
			AddRow(1, "Creator Name", "creator@example.com", true, time.Now(), time.Now()))

	mock.ExpectQuery("SELECT id, name, email, active, created_at, updated_at FROM users WHERE id = ?").
		WithArgs(int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "active", "created_at", "updated_at"}).
			AddRow(2, "Assigned Name", "assigned@example.com", true, time.Now(), time.Now()))

	mock.ExpectExec("INSERT INTO task \\(title, description, created_by, assigned_to, project_id, created_at, finished_at\\) VALUES \\(\\?, \\?, \\?, \\?, \\?, \\?, \\?\\)").
		WithArgs("New Task", "This is a new task", 1, sql.NullInt32{Int32: 2, Valid: true}, 3, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	newTask := taskbus.NewTask{
		Title:       "New Task",
		Description: "This is a new task",
		ProjectID:   3,
		CreatedBy:   1,
		AssignedTo:  sql.NullInt32{Int32: 2, Valid: true},
	}
	task, err := business.Create(ctx, newTask)

	assert.NoError(t, err)
	assert.Equal(t, 1, task.ID)
	assert.Equal(t, "New Task", task.Title)
	assert.Equal(t, "This is a new task", task.Description)
	assert.Equal(t, 3, task.ProjectID)
	assert.NotEmpty(t, task.CreatedAt)
	assert.True(t, task.CreatedAt.After(time.Now().Add(-time.Hour)))
	assert.False(t, task.FinishedAt.Valid)

	assertMockExpectations(t, mock)
}

func TestQuery(t *testing.T) {
	setupMockDB(t)
	defer db.Close()

	mock.ExpectQuery("SELECT id, title, description, created_at, finished_at, created_by, assigned_to, project_id FROM task").
		WillReturnRows(mockTaskRows())

	ctx := context.Background()
	tasks, err := business.Query(ctx)

	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
	assert.Equal(t, "Task 1", tasks[0].Title)
	assert.Equal(t, "Description 1", tasks[0].Description)
	assert.NotEmpty(t, tasks[0].CreatedAt)
	assert.True(t, tasks[0].CreatedAt.After(time.Now().Add(-time.Hour)))
	assert.WithinDuration(t, time.Now(), tasks[0].CreatedAt, time.Minute)
	assert.False(t, tasks[0].FinishedAt.Valid)
	assertMockExpectations(t, mock)
}

func TestQueryByID(t *testing.T) {
	setupMockDB(t)
	defer db.Close()

	mock.ExpectQuery("SELECT id, title, description, project_id, created_at, finished_at, created_by, assigned_to FROM task WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "project_id", "created_at", "finished_at", "created_by", "assigned_to"}).
			AddRow(1, "Task 1", "Description 1", 3, time.Now(), sql.NullTime{}, 1, sql.NullInt32{}))

	ctx := context.Background()
	task, err := business.QueryByID(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, "Task 1", task.Title)
	assert.Equal(t, "Description 1", task.Description)
	assert.Equal(t, 3, task.ProjectID)
	assert.NotEmpty(t, task.CreatedAt)
	assert.True(t, task.CreatedAt.After(time.Now().Add(-time.Hour)))
	assert.False(t, task.FinishedAt.Valid)
	assertMockExpectations(t, mock)
}

func TestUpdate(t *testing.T) {
	setupMockDB(t)
	defer db.Close()

	mock.ExpectExec("^UPDATE task SET title = \\?, description = \\?, assigned_to = \\? WHERE id = \\?$").
		WithArgs("Update Title", "Update Description", sql.NullInt32{}, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := context.Background()
	updateTask := taskbus.UpdateTask{Title: "Update Title", Description: "Update Description"}
	err := business.Update(ctx, 1, updateTask)

	assert.NoError(t, err)
	assertMockExpectations(t, mock)
}

func TestFinish(t *testing.T) {
	setupMockDB(t)
	defer db.Close()

	now := time.Now()
	finishedAt := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.Local)
	nullFinishedAt := sql.NullTime{Time: finishedAt, Valid: true}

	mock.ExpectExec("^UPDATE task SET finished_at = \\? WHERE id = \\?$").
		WithArgs(nullFinishedAt, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := context.Background()
	err := business.Finish(ctx, 1)

	assert.True(t, finishedAt.After(time.Now().Add(-time.Hour)))
	assert.NoError(t, err)
	assertMockExpectations(t, mock)
}

func TestDelete(t *testing.T) {
	setupMockDB(t)
	defer db.Close()

	mock.ExpectExec("DELETE FROM task WHERE id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := context.Background()
	err := business.Delete(ctx, 1)

	assert.NoError(t, err)
	assertMockExpectations(t, mock)
}
