package projectapp

import (
	"TODO-list/business/domain/projectbus"
	"TODO-list/foundation/logger"
	"TODO-list/foundation/web"
	"net/http"
)

// Config contains the dependencies required for initializing the project application.
type Config struct {
	ProjectBus *projectbus.Business
	Logger     *logger.Logger
}

// Routes sets up the HTTP routes for the project-related API endpoints.
func Routes(web *web.App, cfg Config) {
	app := newApp(cfg.ProjectBus)

	web.HandlerFunc(http.MethodPost, "", "/api/project", app.Create, nil)
	web.HandlerFunc(http.MethodGet, "", "/api/project", app.Query, nil)
	web.HandlerFunc(http.MethodGet, "", "/api/project/{id}", app.QueryByID, nil)
	web.HandlerFunc(http.MethodPut, "", "/api/project/{id}", app.Update, nil)
	web.HandlerFunc(http.MethodDelete, "", "/api/project/{id}", app.Delete, nil)
	web.HandlerFunc(http.MethodDelete, "", "/api/project/{id}/deactivate", app.Deactivate, nil)

}
