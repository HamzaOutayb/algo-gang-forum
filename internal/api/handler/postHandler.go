package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"real-time-forum/internal/models"
	utils "real-time-forum/pkg"
)

func (H *Handler) InsertPostsHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil || !H.Service.Database.CheckExpiredCookie(cookie.Value, time.Now()) {
		utils.WriteJson(w, http.StatusUnauthorized, struct {
			Error string `json:"error"`
		}{Error: "Unauthorized"})
		return
	}
	var post models.Post
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		utils.WriteJson(w, 400, "bad request")
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

func (H *Handler) GetPostByIdHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	cookie, err := r.Cookie("session_token")
	userid := 0
	if err != http.ErrNoCookie && H.Service.Database.CheckExpiredCookie(cookie.Value, time.Now()) {
		userid, err = H.Service.Database.GetUser(cookie.Value)
		if err != nil {
			utils.WriteJson(w, http.StatusBadRequest, "bad request")
			return
		}
	}

	posts, err := H.Service.GetPostbyid(idString, userid)
	if err != nil {
		if err.Error() == "not found" {
			utils.WriteJson(w, http.StatusNotFound, "not found")
			return
		}
		if err == sql.ErrNoRows {
			utils.WriteJson(w, http.StatusLocked, "database locked")
			return
		}
		utils.WriteJson(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.WriteJson(w, http.StatusOK, posts)
}
