package repository

func (Database *Database) InsertChat(From,To int,Message []byte) error {
	var Conversations_ID int64
	Database.Db.QueryRow("SELECT id FROM conversations WHERE (user_one = ? AND user_two = ?) OR (user_one = ? AND user_two = ?)",From,To,To,From).Scan(&Conversations_ID)
	if Conversations_ID == 0 {
		Insertchat,err := Database.Db.Exec("INSERT INTO conversations (user_one, user_two) VALUES (?, ?)",From,To)
		if err != nil {
			return err
		}
		Conversations_ID, err = Insertchat.LastInsertId()
		if err != nil {
			return err
		}
	}
	_,err := Database.Db.Exec("INSERT INTO messages (sender_id, content, conversation_id) VALUES (?, ?, ?)",From,string(Message),Conversations_ID)
	if err != nil {
		return err
	}
	return nil
}