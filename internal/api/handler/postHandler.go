package handler

import (
	"fmt"
	"net/http"
	utils "real-time-forum/pkg"
)

/*
	func (H *Handler) GetPostByIdHandler(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := H.Service.GetPostbyid(idString)
	}
*/
func (H *Handler) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := H.Service.Database.Posts(1)
	if err != nil {
		fmt.Println(err)
		utils.WriteJson(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.WriteJson(w, http.StatusOK, posts)
}
