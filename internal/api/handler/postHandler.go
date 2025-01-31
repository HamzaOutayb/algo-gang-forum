package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"real-time-forum/internal/models"
	utils "real-time-forum/pkg"

	"github.com/mattn/go-sqlite3"
)

func (H *Handler) InsertPostsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("test")
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
	fmt.Println(post)

	err = H.Service.CreatePost(post, cookie.Value)
	if err != nil {
		fmt.Println(err)
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

func (H *Handler) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	num, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, nil)
		return
	}

	var Posts []models.Post

	cookie, err := r.Cookie("session_token")
	id := 0
	if err != http.ErrNoCookie && H.Service.Database.CheckExpiredCookie(cookie.Value, time.Now()) {
		id, _ = H.Service.Database.GetUser(cookie.Value)
	}

	Posts, err = H.Service.GetPost(num, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			utils.WriteJson(w, http.StatusOK, []models.Post{})
			return
		case sqlite3.ErrLocked:
			utils.WriteJson(w, http.StatusLocked, struct {
				Error string `json:"error"`
			}{Error: "Database Locked"})
			return
		}
		utils.WriteJson(w, http.StatusInternalServerError, struct {
			Error string `json:"error"`
		}{Error: err.Error()})
		return
	}
	utils.WriteJson(w, http.StatusOK, Posts)
}

func (H *Handler) GetContactHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		utils.WriteJson(w, http.StatusNonAuthoritativeInfo, err)
		return
	}

	contact, err := H.Service.Database.GetContact(cookie.Value)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, "err")
		return
	}
	utils.WriteJson(w, http.StatusOK, contact)
}
