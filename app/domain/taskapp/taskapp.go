package taskapp

import (
	"TODO-list/app/sdk/errs"
	"TODO-list/business/domain/taskbus"
	"TODO-list/foundation/web"
	"context"
	"net/http"
	"strconv"
)

// App handles the application layer of tasks.
type App struct {
	taskBus *taskbus.Business
}

// NewApp creates a new instance of App.
func NewApp(taskBus *taskbus.Business) *App {
	return &App{
		taskBus: taskBus,
	}
}

// Create adds a new task and returns the created task.
func (a *App) Create(ctx context.Context, r *http.Request) web.Encoder {
	var app NewTask
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	taskBus, err := a.taskBus.Create(ctx, toBusNewTask(app))
	if err != nil {
		return errs.New(errs.InternalOnlyLog, err)
	}

	return toAppTask(taskBus)
}

// Query retrieves all tasks.
func (a *App) Query(ctx context.Context, r *http.Request) web.Encoder {
	tasksBus, err := a.taskBus.Query(ctx)
	if err != nil {
		return errs.New(errs.InternalOnlyLog, err)
	}
	return TaskList{Tasks: toAppTasks(tasksBus)}
}

// QueryByID retrieves a task by its ID.
func (a *App) QueryByID(ctx context.Context, r *http.Request) web.Encoder {
	idStr := web.DecodeParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	taskBus, err := a.taskBus.QueryByID(ctx, id)
	if err != nil {
		return errs.New(errs.InternalOnlyLog, err)
	}

	return toAppTask(taskBus)
}

// Update modifies an existing task and returns the updated task.
func (a *App) Update(ctx context.Context, r *http.Request) web.Encoder {
	var ut UpdateTask
	if err := web.Decode(r, &ut); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	err := a.taskBus.Update(ctx, toBusUpdateTask(ut))
	if err != nil {
		return errs.New(errs.InternalOnlyLog, err)
	}

	return nil
}

// Delete removes a task by its ID.
func (a *App) Delete(ctx context.Context, r *http.Request) web.Encoder {
	idStr := web.DecodeParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	err = a.taskBus.Delete(ctx, id)
	if err != nil {
		return errs.New(errs.InternalOnlyLog, err)
	}

	return nil
}

// Finish marks a task as completed.
func (a *App) Finish(ctx context.Context, r *http.Request) web.Encoder {
	idStr := web.DecodeParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	err = a.taskBus.Finish(ctx, id)
	if err != nil {
		return errs.New(errs.InternalOnlyLog, err)
	}

	return nil
}
