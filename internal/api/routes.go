package api

import (
	"database/sql"
	"net/http"

	"real-time-forum/internal/api/handler"
)

func Routes(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.HomeHandler)
	return mux
}
