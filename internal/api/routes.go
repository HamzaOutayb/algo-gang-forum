package api

import (
    "database/sql"
    "net/http"

    "real-time-forum/internal/api/handler"
)

func Routes(db *sql.DB) *http.ServeMux {
    mux := http.NewServeMux()
    FileServer := http.FileServer(http.Dir("./Assets/"))
    mux.Handle("/Assets/", http.StripPrefix("/Assets/", FileServer))
    mux.HandleFunc("/", handler.HomeHandler)
    return mux
}