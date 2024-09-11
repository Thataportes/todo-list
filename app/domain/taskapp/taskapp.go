package taskapp

import (
	"TODO-list/business/domain/taskbus"
	"context"
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
func (a *App) Create(ctx context.Context, nt NewTask) (Task, error) {
	taskBus, err := a.taskBus.Create(ctx, toBusNewTask(nt))
	if err != nil {
		return Task{}, err
	}

	return toAppTask(taskBus), nil
}

// Query retrieves all tasks.
func (a *App) Query(ctx context.Context) ([]Task, error) {
	tasksBus, err := a.taskBus.Query(ctx)
	if err != nil {
		return nil, err
	}
	return toAppTasks(tasksBus), nil
}

// QueryByID retrieves a task by its ID.
func (a *App) QueryByID(ctx context.Context, id int) (Task, error) {
	taskBus, err := a.taskBus.QueryByID(ctx, id)
	if err != nil {
		return Task{}, err
	}

	return toAppTask(taskBus), nil
}

// Update modifies an existing task and returns the updated task.
func (a *App) Update(ctx context.Context, ut UpdateTask) (Task, error) {
	taskBus, err := a.taskBus.Update(ctx, toBusUpdateTask(ut))
	if err != nil {
		return Task{}, err
	}
	return toAppTask(taskBus), nil
}

// Delete removes a task by its ID.
func (a *App) Delete(ctx context.Context, id int) error {
	err := a.taskBus.Delete(ctx, id)
	return err
}

// Finish marks a task as completed.
func (a *App) Finish(ctx context.Context, id int) error {
	err := a.taskBus.Finish(ctx, id)
	return err
}
