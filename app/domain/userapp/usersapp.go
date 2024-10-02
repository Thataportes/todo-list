package userapp

import (
	"TODO-list/app/sdk/errs"
	"TODO-list/business/domain/userbus"
	"TODO-list/foundation/web"
	"context"
	"net/http"
	"strconv"
)

// App represents the application layer for handling user-related requests.
type App struct {
	userBus *userbus.Business
}

// newApp creates a new instance of the App, initializing it with the business layer (userBus).
func newApp(userBus *userbus.Business) *App {
	return &App{
		userBus: userBus,
	}
}

// Create handles the creation of a new user.
func (a *App) Create(ctx context.Context, r *http.Request) web.Encoder {
	var app NewUser
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	userBus, err := a.userBus.Create(ctx, toBusNewUser(app))
	if err != nil {
		return errs.New(errs.InternalOnlyLog, err)
	}

	return toAppUser(userBus)
}

// Query retrieves a list of all users from the business layer.
func (a *App) Query(ctx context.Context, r *http.Request) web.Encoder {
	usersBus, err := a.userBus.Query(ctx)
	if err != nil {
		return errs.New(errs.InternalOnlyLog, err)
	}

	return toAppUsers(usersBus)
}

// QueryById retrieves a specific user by its ID from the business layer.
func (a *App) QueryById(ctx context.Context, r *http.Request) web.Encoder {
	id, err := strconv.Atoi(web.Param(r, "id"))
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	userBus, err := a.userBus.QueryById(ctx, id)
	if err != nil {
		return errs.New(errs.InternalOnlyLog, err)
	}

	return toAppUser(userBus)
}

// QueryByEmail retrieves a specific user by their email from the business layer.
func (a *App) QueryByEmail(ctx context.Context, r *http.Request) web.Encoder {
	email := web.Param(r, "email")

	userBus, err := a.userBus.QueryByEmail(ctx, email)
	if err != nil {
		return errs.New(errs.InternalOnlyLog, err)
	}

	return toAppUser(userBus)
}

// Update modifies an existing user's information.
func (a *App) Update(ctx context.Context, r *http.Request) web.Encoder {
	id, err := strconv.Atoi(web.Param(r, "id"))
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	var uu UpdateUser
	if err := web.Decode(r, &uu); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	err = a.userBus.Update(ctx, id, toBusUpdateUser(uu))
	if err != nil {
		return errs.New(errs.InternalOnlyLog, err)
	}

	return nil
}

// Delete deactivates a user by their ID.
func (a *App) Delete(ctx context.Context, r *http.Request) web.Encoder {
	id, err := strconv.Atoi(web.Param(r, "id"))
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	err = a.userBus.Deactivate(ctx, id)
	if err != nil {
		return errs.New(errs.InternalOnlyLog, err)
	}

	return nil
}
