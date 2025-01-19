package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"real-time-forum/internal/models"
)

func (H *Handler) AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	react := models.React{}
	json.NewDecoder(r.Body).Decode(&react)
	cookie, err := r.Cookie("session_token")
	if err != nil || !H.Service.Database.CheckExpiredCookie(cookie.Value, time.Now()) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, _ := H.Service.Database.GetUser(cookie.Value)

	err = H.Service.CheckReactInput(react)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = H.Service.Database.ReactionService(react, id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response, err := H.Service.LikesTotal(react.Thread_type, react.Thread_id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response.IsLiked, response.IsDisliked = H.Service.GetLikedThread(react.Thread_type, react.Thread_id, id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response)
}
