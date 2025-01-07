package utils

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

func WriteJson(w http.ResponseWriter, statuscode int, Data any) error {
	w.WriteHeader(statuscode)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(Data)
	if err != nil {
		return err
	}
	return nil
}

func CheckExpiredCookie(uid string, date time.Time, db *sql.DB) bool {
	var expired time.Time
	db.QueryRow("SELECT expired_at FROM user_profile WHERE uid = ?", uid).Scan(&expired)

	return date.Compare(expired) <= -1
}

func Contains(slice []string, str string) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == str {
			return true
		}
	}
	return false
}
