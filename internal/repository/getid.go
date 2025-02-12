package repository

import (
	"errors"
)

func (database *Database) GetId(from string) (string, int, error) {
	From := 0
	name := ""
	_ = database.Db.QueryRow("SELECT id,Nickname FROM user WHERE uid = ?", from).Scan(&From, &name)
	//_ = database.Db.QueryRow("SELECT  id,uid FROM user WHERE Nickname = ?", to).Scan(&To,&uid)
	if From == 0 {
		return name, From, errors.New("not exist")
	}
	return name, From, nil
}

func (database *Database) GetId2(from string) (int, error) {
	From := 0
	_ = database.Db.QueryRow("SELECT id FROM user WHERE uid = ?", from).Scan(&From)
	if From == 0 {
		return From, errors.New("not exist")
	}
	return From, nil
}
