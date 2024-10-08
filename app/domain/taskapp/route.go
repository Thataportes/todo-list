package taskapp

import (
	"TODO-list/business/domain/taskbus"
	"TODO-list/foundation/logger"
	"TODO-list/foundation/web"
	"net/http"
)

// Config holds the configuration dependencies for the application.
type Config struct {
	TaskBus *taskbus.Business
	Logger  *logger.Logger
}

// Routes sets up the HTTP routes for the task-related API endpoints.
func Routes(web *web.App, cfg Config) {
	app := newApp(cfg.TaskBus)

	web.HandlerFunc(http.MethodPost, "", "/api/tasks", app.Create, nil)
	web.HandlerFunc(http.MethodGet, "", "/api/tasks", app.Query, nil)
	web.HandlerFunc(http.MethodGet, "", "/api/tasks/{id}", app.QueryByID, nil)
	web.HandlerFunc(http.MethodPut, "", "/api/tasks/{id}", app.Update, nil)
	web.HandlerFunc(http.MethodDelete, "", "/api/tasks/{id}", app.Delete, nil)
	web.HandlerFunc(http.MethodPut, "", "/api/tasks/finish/{id}", app.Finish, nil)

}
