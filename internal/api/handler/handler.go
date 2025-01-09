package handler

import (
	"database/sql"
	"net/http"
	"text/template"

	"real-time-forum/internal/repository"
	"real-time-forum/internal/service"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(db *sql.DB) *Handler {
	userData := repository.Database{
		Db: db,
	}

	userService := service.Service{
		Database: &userData,
	}

	return &Handler{
		Service: &userService,
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("Assets/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}
