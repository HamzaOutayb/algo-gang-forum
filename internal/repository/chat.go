package repository

func (Database *Database) InsertChat(From,To int,Message []byte) error {
	Insertchat,err := Database.Db.Exec("INSERT INTO conversation (user_one, user_two) VALUES (?, ?)",From,To)
	if err != nil {
		return err
	}
	lastid, err := Insertchat.LastInsertId()
	if err != nil {
		return err
	}
	_,err = Database.Db.Exec("INSERT INTO messages (sender_id, content, conversation_id) VALUES (?, ?, ?)",From,string(Message),lastid)
	if err != nil {
		return err
	}
	return nil
}