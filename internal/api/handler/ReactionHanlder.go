package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"real-time-forum/internal/models"
	utils "real-time-forum/pkg"
)

func (H *Handler) ReactionHandler(w http.ResponseWriter, r *http.Request) {
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

	_,err  = H.Service.Database.GetUser(cookie.Value); if err != nil {
		utils.WriteJson(w, http.StatusBadRequest,"bad request")
		return
	}

	err = H.Service.CheckReactInput(react)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	/*err = service.ReactService.ReactionService(react, id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response, err := service.ReactService.LikesTotal(react.Thread_type, react.Thread_id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response.IsLiked, response.IsDisliked = service.ReactService.GetLikedThread(react.Thread_type, react.Thread_id, id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response)*/
}
