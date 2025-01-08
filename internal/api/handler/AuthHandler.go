package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

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

	err := H.Service.LoginUser(&user); if err != nil {
			if err == sqlite3.ErrLocked {
				http.Error(w, "Database Is Busy!", http.StatusLocked)
				return
			}
			// Email
			if err.Error() == models.Errors.InvalidEmail {
				http.Error(w, models.Errors.InvalidEmail, http.StatusBadRequest)
				return
			}
			if err.Error() == models.Errors.LongEmail {
				http.Error(w, models.Errors.LongEmail, http.StatusBadRequest)
				return
			}

			// Password
			if err.Error() == models.Errors.InvalidPassword {
				http.Error(w, models.Errors.InvalidPassword, http.StatusBadRequest)
				return
			}
			// General: User Doesn't Exist
			if err.Error() == models.Errors.InvalidCredentials {
				http.Error(w, models.Errors.InvalidCredentials, http.StatusUnauthorized)
				return
			}

			if err == sql.ErrNoRows {
				http.Error(w, models.Errors.InvalidCredentials, http.StatusUnauthorized)
				return
			}

			log.Println("Unexpected error:", err)
			http.Error(w, "Error While logging To An  Account.", http.StatusInternalServerError)
			return
		}

	utils.SetSessionCookie(w, user.Uuid)
	utils.WriteJson(w, http.StatusOK, "You Logged In Successfuly!")
}

func (H *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodPost {
		utils.WriteJson(w,http.StatusMethodNotAllowed, "Method not allowed")
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
			http.Error(w, "Database Is Busy!", http.StatusLocked)
			return
		}
		// Username
		if err.Error() == models.Errors.InvalidUsername {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Password
		if err.Error() == models.Errors.InvalidPassword {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Email
		if err.Error() == models.Errors.InvalidEmail {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err.Error() == models.Errors.LongEmail {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// General
		if err.Error() == models.Errors.UserAlreadyExist {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Error While Registering The User.", http.StatusInternalServerError)
		return
	}
	err = utils.WriteJson(w,http.StatusOK,"You'v loged succesfuly"); if err != nil {
		http.Error(w,"internal server error",http.StatusInternalServerError)
		return
	}
}