package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"text/template"
	"time"

	"real-time-forum/internal/models"
	"real-time-forum/internal/repository"
	"real-time-forum/internal/service"
	utils "real-time-forum/pkg"
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

func (H *Handler) InsertPostsHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil || !utils.CheckExpiredCookie(cookie.Value, time.Now(), H.Service.Database.Db) {
		utils.WriteJson(w, http.StatusUnauthorized, struct {
			Error string `json:"error"`
		}{Error: "Unauthorized"})
		return
	}
	var post models.Post
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		utils.WriteJson(w, 500, "internal server err")
		return
	}

	err = H.Service.CreatePost(post, cookie.Value)
	if err != nil {
		switch err.Error() {
		case models.PostErrors.ContentLength:
			utils.WriteJson(w, http.StatusBadRequest, struct {
				Error string `json:"error"`
			}{Error: err.Error()})
			return
		case models.PostErrors.TitleLength:
			utils.WriteJson(w, http.StatusBadRequest, struct {
				Error string `json:"error"`
			}{Error: err.Error()})
			return
		case models.PostErrors.CategoryDoesntExist:
			utils.WriteJson(w, http.StatusBadRequest, struct {
				Error string `json:"error"`
			}{Error: err.Error()})
			return
		case models.UserErrors.UserNotExist:
			utils.WriteJson(w, http.StatusBadRequest, struct {
				Error string `json:"error"`
			}{Error: err.Error()})
			return
		case sql.ErrNoRows.Error():
			utils.WriteJson(w, http.StatusBadRequest, struct {
				Error string `json:"error"`
			}{Error: err.Error()})
			return
		}

		utils.WriteJson(w, http.StatusInternalServerError, struct {
			Error string `json:"error"`
		}{Error: err.Error()})
		return
	}
	w.WriteHeader(http.StatusCreated)
}