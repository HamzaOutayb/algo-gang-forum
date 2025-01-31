package repository

import (
	"errors"
	"fmt"
)

func (database *Database) GetId(from string) (string,int, error) {
	From:= 0
	name := ""
	_ = database.Db.QueryRow("SELECT id,Nickname FROM user WHERE uid = ?", from).Scan(&From,&name)
	//_ = database.Db.QueryRow("SELECT  id,uid FROM user WHERE Nickname = ?", to).Scan(&To,&uid)
	if From == 0 {
		fmt.Println(From)
		return name,From, errors.New("not exist")
	}
	return name,From, nil
}


func (database *Database) GetId2(from,to string) (int,int, error) {
	From:= 0
	To:= 0
	
	_ = database.Db.QueryRow("SELECT id FROM user WHERE uid = ?", from).Scan(&From)
	_ = database.Db.QueryRow("SELECT id FROM user WHERE Nickname = ?", to).Scan(&To)
	if From == 0 || To == 0 {
		fmt.Println(From)
		return From,To, errors.New("not exist")
	}
	return From,To, nil
}
