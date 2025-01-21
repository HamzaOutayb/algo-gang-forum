package repository

import (
	"errors"
	"fmt"
)

func (database *Database) GetId(from, to string) (int, int, error) {
	From, To := 0, 0
	_ = database.Db.QueryRow("SELECT id FROM user WHERE uid = ?", from).Scan(&From)
	_ = database.Db.QueryRow("SELECT id FROM user WHERE Nickname = ?", to).Scan(&To)
	if From == 0 || To == 0 {
		fmt.Println(From, To)
		return From, To, errors.New("not exist")
	}
	return From, To, nil
}
