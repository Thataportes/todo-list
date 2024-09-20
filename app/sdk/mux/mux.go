package mux

import (
	"TODO-list/app/domain/taskapp"
	"TODO-list/business/domain/taskbus"
	"TODO-list/foundation/web"
	"database/sql"
	"net/http"
)

// Config holds the dependencies required for initializing the web API.
type Config struct {
	DB *sql.DB
}

// WebAPI initializes the web application with the given configuration.
func WebAPI(cfg Config) (http.Handler, error) {
	app := web.NewApp(nil)

	taskBus := taskbus.NewBusiness(cfg.DB)

	taskapp.Routes(app, taskapp.Config{
		TaskBus: taskBus,
	})
	return app, nil
}
