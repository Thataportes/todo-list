package mux

import (
	"TODO-list/app/domain/taskapp"
	"TODO-list/business/domain/taskbus"
	"TODO-list/foundation/logger"
	"TODO-list/foundation/web"
	"TODO-list/user/domain/usersapp"
	"TODO-list/userbusiness/domain/userbus"
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

	taskBus := taskbus.NewBusiness(cfg.DB)

	taskapp.Routes(app, taskapp.Config{
		TaskBus: taskBus,
		Logger:  cfg.Log,
	})

	userBus := userbus.NewBusiness(cfg.DB)
	usersapp.Routes(app, usersapp.Config{
		UserBus: userBus,
	})
	return app, nil
}
