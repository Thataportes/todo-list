package userapp

import (
	"TODO-list/business/domain/userbus"
	"TODO-list/foundation/logger"
	"TODO-list/foundation/web"
	"net/http"
)

// Config contains the dependencies required for initializing the user application.
type Config struct {
	UserBus *userbus.Business
	Logger  *logger.Logger
}

// Routes sets up the HTTP routes for the user-related API endpoints.
func Routes(app *web.App, cfg Config) {
	appUser := newApp(cfg.UserBus)

	app.HandlerFunc(http.MethodPost, "", "/api/users", appUser.Create, nil)
	app.HandlerFunc(http.MethodGet, "", "/api/users", appUser.Query, nil)
	app.HandlerFunc(http.MethodGet, "", "/api/users/{id}", appUser.QueryById, nil)
	app.HandlerFunc(http.MethodGet, "", "/api/users/email/{email}", appUser.QueryByEmail, nil)
	app.HandlerFunc(http.MethodPut, "", "/api/users/{id}", appUser.Update, nil)
	app.HandlerFunc(http.MethodDelete, "", "/api/users/{id}", appUser.Delete, nil)
}
