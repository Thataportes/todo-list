package projectapp

import (
	"TODO-list/app/sdk/errs"
	"TODO-list/business/domain/projectbus"
	"TODO-list/foundation/web"
	"context"
	"fmt"
	"net/http"
	"strconv"
)

// App handles the application layer for project-related operations.
type App struct {
	projectBus *projectbus.Business
}

// newApp creates a new instance of App with the provided business layer (projectBus).
func newApp(projectBus *projectbus.Business) *App {
	return &App{projectBus: projectBus}
}

// Create handles the creation of a new project.
func (a *App) Create(ctx context.Context, r *http.Request) web.Encoder {
	var app NewProject
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	if app.CreatedBy == 0 {
		return errs.New(errs.InvalidArgument, fmt.Errorf("created_by is required"))
	}

	projectBus, err := a.projectBus.Create(ctx, toBusNewProject(app))
	if err != nil {
		return errs.New(errs.InternalOnlyLog, err)
	}

	return toAppProject(projectBus)
}

// Query retrieves a list of all projects from the business layer.
func (a *App) Query(ctx context.Context, r *http.Request) web.Encoder {
	projectsBus, err := a.projectBus.Query(ctx)
	if err != nil {
		return errs.New(errs.InternalOnlyLog, err)
	}
	return toAppProjects(projectsBus)
}

// QueryByID retrieves a specific project by its ID from the business layer.
func (a *App) QueryByID(ctx context.Context, r *http.Request) web.Encoder {
	id, err := strconv.Atoi(web.Param(r, "id"))
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	projectBus, err := a.projectBus.QueryById(ctx, id)
	if err != nil {
		return errs.New(errs.InternalOnlyLog, err)
	}

	return toAppProject(projectBus)
}

// Update modifies an existing project's information.
func (a *App) Update(ctx context.Context, r *http.Request) web.Encoder {
	id, err := strconv.Atoi(web.Param(r, "id"))
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	var up UpdateProject
	if err := web.Decode(r, &up); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	err = a.projectBus.Update(ctx, id, toBusUpdateProject(up))
	if err != nil {
		return errs.New(errs.InternalOnlyLog, err)
	}

	return nil
}

// Delete removes a project permanently from the database.
func (a *App) Delete(ctx context.Context, r *http.Request) web.Encoder {
	id, err := strconv.Atoi(web.Param(r, "id"))
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	err = a.projectBus.Delete(ctx, id)
	if err != nil {
		return errs.New(errs.InternalOnlyLog, err)
	}

	return nil
}

// Deactivate sets the project's status to inactive without deleting it.
func (a *App) Deactivate(ctx context.Context, r *http.Request) web.Encoder {
	id, err := strconv.Atoi(web.Param(r, "id"))
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	err = a.projectBus.Deactivate(ctx, id)
	if err != nil {
		return errs.New(errs.InternalOnlyLog, err)
	}

	return nil
}
