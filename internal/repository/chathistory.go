package repository

func (Database *Database) HistoryMessages(from,to int) (map[string]map[string]string,error) {
	conversations_id := 0
	conversations := make(map[string]map[string]string)
	Database.Db.QueryRow("SELECTE id FROM conversations WHERE user_one = ? AND user_two = ?",from,to).Scan(conversations_id)
	rows,err := Database.Db.Query("SELECTE sender_id,content,created_at FROM messages WHERE conversation_id = ? ",from,to,conversations_id)
	if (err != nil){
		return conversations,err
	}
	for rows.Next() {
		var sender, message, date string
		err := rows.Scan(&sender, &message, &date)
		if err != nil {
			return nil, err
		}

		if conversations[sender] == nil {
			conversations[sender] = make(map[string]string)
		}
		conversations[sender][message] = date
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
return conversations,nil
	
}