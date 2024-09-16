package taskbus_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"TODO-list/business/domain/taskbus"

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
	business = taskbus.NewBusiness(db)
}

func mockTaskRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "title", "description", "created_at", "finished_at"}).
		AddRow(1, "Task 1", "Description 1", time.Now(), sql.NullTime{}).
		AddRow(2, "Task 2", "Description 2", time.Now(), sql.NullTime{})

}

func assertMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate(t *testing.T) {
	setupMockDB(t)
	defer db.Close()

	mock.ExpectExec("INSERT INTO task").
		WithArgs("New Task", "This is a new task", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	newTask := taskbus.NewTask{Title: "New Task", Description: "This is a new task"}
	task, err := business.Create(ctx, newTask)

	assert.NoError(t, err)
	assert.Equal(t, 1, task.ID)
	assert.Equal(t, "New Task", task.Title)
	assert.Equal(t, "This is a new task", task.Description)
	assert.NotEmpty(t, task.CreatedAt)
	assert.True(t, task.CreatedAt.After(time.Now().Add(-time.Hour)))
	assert.True(t, task.FinishedAt.IsZero())
	assertMockExpectations(t, mock)
}

func TestQuery(t *testing.T) {
	setupMockDB(t)
	defer db.Close()

	mock.ExpectQuery("SELECT id, title, description, created_at, finished_at FROM task").
		WillReturnRows(mockTaskRows())

	ctx := context.Background()
	tasks, err := business.Query(ctx)

	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
	assert.Equal(t, "Task 1", tasks[0].Title)
	assert.Equal(t, "Description 1", tasks[0].Description)
	assert.NotEmpty(t, tasks[0].Title)
	assert.True(t, tasks[0].CreatedAt.After(time.Now().Add(-time.Hour)))
	assert.True(t, tasks[0].FinishedAt.IsZero())
	assertMockExpectations(t, mock)
}

func TestQueryByID(t *testing.T) {
	setupMockDB(t)
	defer db.Close()

	row := sqlmock.NewRows([]string{"id", "title", "description", "created_at", "finished_at"}).
		AddRow(1, "Task 1", "Description 1", time.Now(), sql.NullTime{})

	mock.ExpectQuery("SELECT id, title, description, created_at, finished_at FROM task WHERE id = ?").
		WithArgs(1).
		WillReturnRows(row)

	ctx := context.Background()
	task, err := business.QueryByID(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, "Task 1", task.Title)
	assert.Equal(t, "Description 1", task.Description)
	assert.NotEmpty(t, task.CreatedAt)
	assert.True(t, task.CreatedAt.After(time.Now().Add(-time.Hour)))
	assert.True(t, task.FinishedAt.IsZero())
	assertMockExpectations(t, mock)
}

func TestUpdate(t *testing.T) {
	setupMockDB(t)
	defer db.Close()

	mock.ExpectExec("^UPDATE task SET title = \\?, description = \\? WHERE id = \\?$").
		WithArgs("Update Title", "Update Description", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := context.Background()
	updateTask := taskbus.UpdateTask{ID: 1, Title: "Update Title", Description: "Update Description"}
	err := business.Update(ctx, updateTask)

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
