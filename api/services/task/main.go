package main

import (
	"TODO-list/app/sdk/mux"
	"TODO-list/foundation/logger"
	"TODO-list/foundation/otel"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var log *logger.Logger
	traceIDFn := func(ctx context.Context) string {
		return otel.GetTraceID(ctx)
	}

	log = logger.New(os.Stdout, logger.LevelInfo, "TASKS", traceIDFn)

	// Create a root context for managing application lifecycle.
	ctx := context.Background()

	// Run the main logic of the application and handle any errors.
	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "err", err)
		os.Exit(1)
	}
}

// run sets up the application, including database connection, server initialization,
// and graceful shutdown logic.
func run(ctx context.Context, log *logger.Logger) error {
	log.BuildInfo(ctx)
	log.Info(ctx, "startup", "status", "initializing database support")

	// Open a connection to the MySQL database using a DSN (Data Source Name).
	// The `parseTime=true` parameter ensures that MySQL DATETIME and TIMESTAMP
	// fields are automatically parsed into Go's time.Time type.
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/todolist?parseTime=true")
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	// Ensure that the database connection is closed when the function exits.
	defer db.Close()

	// Create a channel to listen for interrupt signals (e.g., SIGINT, SIGTERM).
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	log.Info(ctx, "startup", "status", "initializing V1 API support")

	// cfgMux defines the configuration for the mux-based web API, which includes
	// the database connection.
	cfgMux := mux.Config{
		DB:  db,
		Log: log,
	}

	// webAPI initializes a new WebAPI instance with the provided configuration.
	webAPI, err := mux.WebAPI(cfgMux)
	if err != nil {
		return fmt.Errorf("webapi: %w", err)
	}

	// Create an HTTP server with the API handler and specify the server address.
	api := http.Server{
		Addr:     "0.0.0.0:8080",
		Handler:  webAPI,
		ErrorLog: logger.NewStdLogger(log, logger.LevelError),
	}

	// Create a channel to capture any errors from the HTTP server.
	serverErrors := make(chan error, 1)
	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)
		// Start the HTTP server and capture any errors.
		serverErrors <- api.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Graceful Shutdown

	// Wait for a termination signal or a server error to gracefully shutdown the application.
	select {
	case err := <-serverErrors:
		// If there is an error from the HTTP server, return it.
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		// If a shutdown signal is received, log the shutdown process.
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		// Create a context with a timeout to ensure the server shuts down gracefully.
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		// Attempt to gracefully shut down the HTTP server.
		if err := api.Shutdown(ctx); err != nil {
			// If there is an error shutting down, forcibly close the server.
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	// Return nil to indicate a successful shutdown.
	return nil
}
