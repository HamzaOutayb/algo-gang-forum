package database

import (
	"database/sql"
	"real-time-forum/internal/api/bcryptp"
)

var db *sql.DB

func OpenDb() (*sql.DB, error) {
	var err error
	db, err = sql.Open("sqlite3", "./forum.db")
	return db, err
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
