package repository

func (database *Database) GetId(from, to string) (int, int, error) {
	From, To := 0, 0
	_ = database.Db.QueryRow("SELECTE user_id FROM user WHERE Nickname = ?", from).Scan(&From)
	_ = database.Db.QueryRow("SELECTE user_id FROM user WHERE Nickname = ?", to).Scan(&To)
	return From, To, nil
}
