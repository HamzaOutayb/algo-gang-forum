package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	err := H.Service.LoginUser(&user)
	if err != nil {
		if err == sqlite3.ErrLocked {
			http.Error(w, "Database Is Busy!", http.StatusLocked)
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
			fmt.Println("1", err.Error())
			utils.WriteJson(w, http.StatusBadRequest, models.Errors.InvalidCredentials)
			return
		}

		if err == sql.ErrNoRows {
			fmt.Println("2", err.Error())
			utils.WriteJson(w, http.StatusBadRequest, models.Errors.InvalidCredentials)
			return
		}

		http.Error(w, "Error While logging To An  Account.", http.StatusInternalServerError)
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
