package handler

import (
	"encoding/json"
	"net/http"
	"real-time-forum/internal/api/bcryptp"
	"real-time-forum/internal/api/database"
)

type Login_session struct {
	Email    string
	Password string
}

func Signin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		//error message
		return
	}

	var Log Login_session
	if erro := json.NewDecoder(r.Body).Decode(&Log); erro != nil {
		//error message
		return
	}
	r.ParseForm()
	_, err := r.Cookie("session")
	if err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	doz, err := database.GetLogin(Log.Email, Log.Password)
	if err != nil {
		//error message
		return
	}
	if doz {
		session, err := bcryptp.CreateSession()
		if err != nil {
			//error message
			return
		}
		err = database.AddSession(session.String(), Log.Email)
		if err != nil {
			//error message
			return
		}
		cookie := http.Cookie{
			Name:  "session",
			Value: session.String(),
		}
		http.SetCookie(w, &cookie)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("login ok")
	} else {
		errorMessage := ""
		if Log.Email != "" || Log.Password != "" {
			errorMessage = "Password or email not working"
		}
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorMessage)
	}
}
