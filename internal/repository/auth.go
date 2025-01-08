package repository

import (
	"database/sql"
	"time"

	"real-time-forum/internal/api/bcryptp"
)

type Database struct {
	Db *sql.DB
}

func GetLogin(email, Inpassword string) (bool, error) {
	var password string
	err := db.QueryRow("SELECT password FROM user WHERE email =  $1", email).Scan(&password)
	if err != nil {
		return false, err
	}

	if bcryptp.CheckPasswordHash(Inpassword, password) {
		return true, nil
	}

	return false, nil
}

func AddSession(session string, email string) error {
	user_id := 0
	err := db.QueryRow("SELECT user_id FROM user WHERE email = ?", email).Scan(&user_id)
	if err != nil {
		return err
	}
	_, _ = db.Exec("DELETE from session WHERE user_id = ?", user_id)

	_, err = db.Exec("INSERT INTO  session (session,user_id) VALUES (?,?)", session, user_id)
	return err
}

func (database *Database) CheckIfUserExists(username, email string) bool {
	var uname string
	var uemail string
	database.Db.QueryRow("SELECT username, email FROM user WHERE username = ? AND email = ?",
		username, email).Scan(&uname, &uemail)
	return uname == username && uemail == email
}

func (database *Database) GetUserPassword(email string) (string, error) {
	var password string
	err := database.Db.QueryRow("SELECT password FROM user WHERE email = ?",
		email).Scan(&password)
	return password, err
}

func (database *Database) UpdateUuid(uuid, email string) error {
	expire := time.Now().Add(time.Hour)
	_, err := database.Db.Exec("UPDATE user SET uid = ?, expired_at = ? WHERE email = ?", uuid, expire, email)
	return err
}
