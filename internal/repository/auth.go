package repository

import (
	"database/sql"
	"real-time-forum/internal/api/bcryptp"
)

type Database struct {
	Db *sql.DB
}

func GetLogin(email, password string) (bool, error) {
	var passwords string
	err := db.QueryRow("SELECT password FROM users WHERE email =  $1", email).Scan(&passwords)
	if err != nil {
		return false, err
	}

	if bcryptp.CheckPasswordHash(password, passwords) {
		return true, nil
	}

	return false, nil
}

func AddSession(session string, email string) error {
	user_id := 0
	err := db.QueryRow("SELECT user_id FROM users WHERE email = ?", email).Scan(&user_id)
	if err != nil {
		return err
	}
	_, _ = db.Exec("DELETE from session WHERE user_id = ?", user_id)

	_, err = db.Exec("INSERT INTO  session (session,user_id) VALUES (?,?)", session, user_id)
	return err
}
