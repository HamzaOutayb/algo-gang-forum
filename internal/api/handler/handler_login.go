package handler

import (
	"encoding/json"
	"net/http"

	"real-time-forum/internal/api/bcryptp"
	"real-time-forum/internal/repository"
	utils "real-time-forum/pkg"
)

type Login_session struct {
	Email    string
	Password string
}

func Signin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// error message
		return
	}

	var Log Login_session
	if erro := json.NewDecoder(r.Body).Decode(&Log); erro != nil {
		// error message
		return
	}
	r.ParseForm()
	_, err := r.Cookie("session")
	if err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	doz, err := repository.GetLogin(Log.Email, Log.Password)
	if err != nil {
		// error message
		return
	}
	if doz {
		session, err := bcryptp.CreateSession()
		if err != nil {
			// error message
			return
		}
		err = repository.AddSession(session.String(), Log.Email)
		if err != nil {
			// error message
			return
		}
		cookie := http.Cookie{
			Name:  "session",
			Value: session.String(),
		}
		http.SetCookie(w, &cookie)
		utils.WriteJson(w, 200, "log succesfuly")
	} else {
		errorMessage := ""
		if Log.Email != "" || Log.Password != "" {
			errorMessage = "Password or email not working"
		}
		utils.WriteJson(w,400,errorMessage)
	}
}
