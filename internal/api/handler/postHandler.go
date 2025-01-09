package handler

import (
	"net/http"
	utils "real-time-forum/pkg"
)

func (H *Handler) GetPostByIdHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := H.Service.GetPostbyid(idString); if err != nil {
		if err.Error() == "bad request" {
			utils.WriteJson(w, http.StatusBadRequest, "bad request")
			return
		}
	}
}
