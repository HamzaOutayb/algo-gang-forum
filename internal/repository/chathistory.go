package repository

func (Database *Database) HistoryMessages(from,to int) ([]string,error) {
	row,err := Database.Db.Query("s")
	if err != nil {
		return []string{}, err
	}
}