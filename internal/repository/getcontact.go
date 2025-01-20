package repository

import (
	"fmt"

)

type Contact struct {
	Contact_list []string
}


func (database *Database) GetContact(uid string) (Contact, error) {
    result := []string{}
    rows, err := database.Db.Query("SELECT Nickname FROM user WHERE uid != ?",uid)
    if err != nil {
        fmt.Println(err)
        return Contact{}, err
    }
    for rows.Next() {
        name := ""
        err := rows.Scan(&name)
		if err != nil {
			fmt.Println(err)
		}
        result = append(result, name)
    }
	contact := Contact{
		Contact_list: result,
	}
    return contact, nil
}
