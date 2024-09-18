package mux

import (
	"TODO-list/foundation/web"
	"database/sql"
	"net/http"
)

type Config struct {
	DB *sql.DB
}

func WebAPI(cfg Config) http.Handler {
	app := web.NewApp(nil)
}
