package api

import (
	"database/sql"
	"net/http"
	handler "real-time-forum/internal/api/handler"
)

func Routes(db *sql.DB) *http.ServeMux {
	d := handler.NewHandler(db)
	mux := http.NewServeMux()
	FileServer := http.FileServer(http.Dir("./Assets/"))
	mux.Handle("/Assets/", http.StripPrefix("/Assets/", FileServer))
	mux.HandleFunc("/", handler.HomeHandler)
	mux.HandleFunc("/Signin", d.Signin)
	mux.HandleFunc("/Signup", d.Signup)
	mux.HandleFunc("/post", d.InsertPostsHandler)
	//mux.HandleFunc("GET /post/{id}", d.GetPostByIdHandler)
	return mux
}
