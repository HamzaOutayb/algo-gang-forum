package repository

import "real-time-forum/internal/models"

func (d *Database) CheckIfThreadExists(react models.React) bool {
	exist := false
	if react.Thread_type == "post" {
		d.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM post WHERE id = ?)", react.Thread_id).Scan(&exist)
		return exist
	} else {
		d.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM comment WHERE id = ?)", react.Thread_id).Scan(&exist)
		return exist
	}
}
