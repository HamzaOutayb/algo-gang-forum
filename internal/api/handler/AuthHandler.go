package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"real-time-forum/internal/models"
	utils "real-time-forum/pkg"

	"github.com/mattn/go-sqlite3"
)

func (H *Handler) Signin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJson(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var user models.User
	if erro := json.NewDecoder(r.Body).Decode(&user); erro != nil {
		utils.WriteJson(w, http.StatusBadRequest, "Bad request")
		return
	}
	err := H.Service.LoginUser(&user)
	if err != nil {
		if err == sqlite3.ErrLocked {
			utils.WriteJson(w, http.StatusLocked, "Database Is Busy!")
			return
		}
		// Email
		if err.Error() == models.Errors.InvalidEmail {
			utils.WriteJson(w, http.StatusBadRequest, models.Errors.InvalidEmail)
			return
		}
		if err.Error() == models.Errors.LongEmail {
			utils.WriteJson(w, http.StatusBadRequest, models.Errors.LongEmail)
			return
		}

		// Password
		if err.Error() == models.Errors.InvalidPassword {
			utils.WriteJson(w, http.StatusBadRequest, models.Errors.InvalidPassword)
			return
		}
		// General: User Doesn't Exist
		if err.Error() == models.Errors.InvalidCredentials {
			utils.WriteJson(w, http.StatusBadRequest, models.Errors.InvalidCredentials)
			return
		}

		if err == sql.ErrNoRows {
			utils.WriteJson(w, http.StatusBadRequest, models.Errors.InvalidCredentials)
			return
		}

		utils.WriteJson(w, http.StatusInternalServerError,"Error While logging To An  Account.")
		return
	}

	utils.SetSessionCookie(w, user.Uuid)
	utils.WriteJson(w, http.StatusOK, "You Logged In Successfuly!")
}

func (H *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJson(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var user models.User
	if erro := json.NewDecoder(r.Body).Decode(&user); erro != nil {
		utils.WriteJson(w, http.StatusBadRequest, "Bad request")
		return
	}
	// Proccess Data and Insert it
	err := H.Service.RegisterUser(&user)
	if err != nil {
		if err == sqlite3.ErrLocked {
			utils.WriteJson(w, http.StatusLocked, "Database Is Busy!")
			return
		}
		// Username
		if err.Error() == models.Errors.InvalidUsername {
			utils.WriteJson(w, http.StatusBadRequest, err.Error())

			return
		}

		// Age
		if err.Error() == models.UserErrors.InvalideAge {
			utils.WriteJson(w, http.StatusBadRequest, err.Error())

			return
		}

		// Password
		if err.Error() == models.Errors.InvalidPassword {
			utils.WriteJson(w, http.StatusBadRequest, err.Error())

			return
		}
		// Email
		if err.Error() == models.Errors.InvalidEmail {
			utils.WriteJson(w, http.StatusBadRequest, err.Error())

			return
		}
		if err.Error() == models.Errors.LongEmail {
			utils.WriteJson(w, http.StatusBadRequest, err.Error())
			return
		}
		// General
		if err.Error() == models.Errors.UserAlreadyExist {
			utils.WriteJson(w, http.StatusConflict, models.Errors.UserAlreadyExist)
			return
		}
		utils.WriteJson(w, http.StatusInternalServerError, "Error While Registering The User.")
		return
	}
	utils.WriteJson(w, http.StatusOK, "You'v loged succesfuly")
}

func (H *Handler) LougoutHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, "bad request")
	}
	_, user_id, err := H.Service.Database.GetId(user.Uuid)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, "bad request")
		return
	}
	utils.DeleteSessionCookie(w, user.Uuid)
	mu.Lock()
	statusmap[user_id] = false
	mu.Unlock()
	broadcast(conns, statusmap)
	fmt.Println("test")
	utils.WriteJson(w, http.StatusOK, "You Logged Out Successfuly!")
}

func (H *Handler) InfoHandler(w http.ResponseWriter, r *http.Request) {
	var Authorized bool
	// parse user uid
	uid := ""
	userUID, errCookie := r.Cookie("session_token")
	if errCookie == nil {
		uid = userUID.Value
	}

	// Get Info Data
	Authorized, err := H.Service.GetInfoData(uid)
	if err != nil {
		if err == sqlite3.ErrLocked {
			utils.WriteJson(w, http.StatusLocked, struct {
				Message string `json:"message"`
			}{Message: "Database Locked"})
			return
		}

		utils.WriteJson(w, http.StatusInternalServerError, struct {
			Message string `json:"message"`
		}{Message: "Internal Server Error"})
		return
	}
	if errCookie == nil && !H.Service.Database.CheckExpiredCookie(userUID.Value, time.Now()) {
		Authorized = false
		utils.DeleteSessionCookie(w, userUID.Value)
	}

	utils.WriteJson(w, http.StatusOK, Authorized)
}
