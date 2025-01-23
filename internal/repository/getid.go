package repository

import (
	"errors"
	"fmt"
)

func (database *Database) GetId(from, to string) (string,string,int, int, error) {
	From, To := 0, 0
	uid,name := "",""
	_ = database.Db.QueryRow("SELECT id,Nickname FROM user WHERE uid = ?", from).Scan(&From,&name)
	_ = database.Db.QueryRow("SELECT  id,uid FROM user WHERE Nickname = ?", to).Scan(&To,&uid)
	if From == 0 || To == 0 || uid == ""{
		fmt.Println(From, To,uid)
		return name,uid,From, To, errors.New("not exist")
	}
	return name,uid,From, To, nil
}
