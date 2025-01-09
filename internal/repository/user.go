package repository

func (Database *Database) GetUser(uid string) (int, error) {
	var id int
	err := Database.Db.QueryRow("SELECT id FROM user WHERE uid = ?", uid).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}