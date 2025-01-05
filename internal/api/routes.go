package api

import (
    "database/sql"
    "net/http"

    "real-time-forum/internal/api/handler"
)

func Routes(db *sql.DB) *http.ServeMux {
    cssFs := http.FileServer(http.Dir("./internal/repository/style/"))
    jsFs := http.FileServer(http.Dir("./internal/repository/js/"))
    fs := http.FileServer(http.Dir("/internal/repository/"))
    mux := http.NewServeMux()
    mux.Handle("/style/", http.StripPrefix("/style/", cssFs))
    mux.Handle("/js/", http.StripPrefix("/js/", jsFs))
    mux.Handle("/internal/repository/", http.StripPrefix("/internal/repository/", fs))
    mux.HandleFunc("/", handler.HomeHandler)
    return mux
}