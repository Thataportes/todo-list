package taskapp

import (
	"TODO-list/business/domain/taskbus"
	"TODO-list/foundation/web"
	"net/http"
)

// Config holds the configuration dependencies for the application.
type Config struct {
	TaskBus *taskbus.Business
}

// Routes sets up the HTTP routes for the task-related API endpoints.
func Routes(web *web.App, cfg Config) {
	app := newApp(cfg.TaskBus)

	web.HandlerFunc(http.MethodPost, "", "/api/tasks", app.Create, nil)
	web.HandlerFunc(http.MethodGet, "", "/api/tasks", app.Query, nil)
	web.HandlerFunc(http.MethodGet, "", "/api/tasks", app.QueryByID, nil)
	web.HandlerFunc(http.MethodPut, "", "/api/tasks", app.Update, nil)
	web.HandlerFunc(http.MethodDelete, "", "/api/tasks", app.Delete, nil)
	web.HandlerFunc(http.MethodPost, "", "/api/tasks", app.Finish, nil)

}
