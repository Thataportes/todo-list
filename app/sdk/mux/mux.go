package mux

import (
	"TODO-list/app/domain/taskapp"
	"TODO-list/app/domain/userapp"
	"TODO-list/business/domain/taskbus"
	"TODO-list/business/domain/userbus"
	"TODO-list/foundation/logger"
	"TODO-list/foundation/web"
	"context"
	"database/sql"
	"net/http"
)

// Config holds the dependencies required for initializing the web API.
type Config struct {
	Log *logger.Logger
	DB  *sql.DB
}

// WebAPI initializes the web application with the given configuration.
func WebAPI(cfg Config) (http.Handler, error) {
	logger := func(ctx context.Context, msg string, args ...any) {
		cfg.Log.Info(ctx, msg, args...)
	}
	app := web.NewApp(logger)

	userBus := userbus.NewBusiness(cfg.DB)
	taskBus := taskbus.NewBusiness(cfg.DB, userBus)

	userapp.Routes(app, userapp.Config{
		UserBus: userBus,
		Logger:  cfg.Log,
	})

	taskapp.Routes(app, taskapp.Config{
		TaskBus: taskBus,
		Logger:  cfg.Log,
	})
	return app, nil
}
