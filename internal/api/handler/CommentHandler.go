package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"real-time-forum/internal/models"
	utils "real-time-forum/pkg"

	"github.com/mattn/go-sqlite3"
)

func (H *Handler) AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	// parse data
	comment := models.Comment{}
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	cookie, err := r.Cookie("session_token")
	if err != nil || !H.Service.Database.CheckExpiredCookie(cookie.Value, time.Now()) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Validate Inputs
	err = H.Service.ValidateInput(comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the user exist
	usrid, _ := H.Service.Database.GetUser(cookie.Value)
	if usrid == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// add userUid to the comment
	comment.UserUID = cookie.Value

	// add the comment using the AddComment from the service layer
	err = H.Service.AddComment(comment)
	if err != nil {
		switch err.Error() {
		case sqlite3.ErrLocked.Error():
			http.Error(w, "Database locked", http.StatusLocked)
			return
		case models.PostErrors.PostNotExist:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		default:
			log.Printf("Unexpected Error when we add comment %s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	// return a success response
	utils.WriteJson(w, http.StatusCreated, struct{ Message string }{
		Message: "Your comment added successfuly",
	})
}