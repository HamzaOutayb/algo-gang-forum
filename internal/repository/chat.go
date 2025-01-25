package repository

import (
	"real-time-forum/internal/models"
)

func (Database *Database) InsertChat(From, To int, Message []byte) error {
	Insertchat, err := Database.Db.Exec("INSERT INTO conversation (user_one, user_two) VALUES (?, ?)", From, To)
	if err != nil {
		return err
	}
	lastid, err := Insertchat.LastInsertId()
	if err != nil {
		return err
	}
	_, err = Database.Db.Exec("INSERT INTO messages (sender_id, content, conversation_id) VALUES (?, ?, ?)", From, string(Message), lastid)
	if err != nil {
		return err
	}
	return nil
}

func (Database *Database) GetChatAndUserID(pagenm int, usrid int) ([]models.Chat, error) {
	rows, err := Database.Db.Query(`
	SELECT DISTINCT
  chat_id,
  CASE 
    WHEN user_one = ${1} THEN user_two
    ELSE user_one
  END AS chatted_user
FROM chats
WHERE ${1} IN (user_one, user_two);
`, usrid)
	if err != nil {
		return []models.Chat{}, err
	}

	var Chats []models.Chat
	for rows.Next() {
		var chat models.Chat
		err := rows.Scan(&chat.Id, &chat.FriendsId)
		if err != nil {
			return []models.Chat{}, err
		}
		chat.Nickname, err = Database.GetuserNickname(chat.FriendsId)
		if err != nil {
			return []models.Chat{}, err
		}

		Chats = append(Chats, chat)
	}

	return Chats, nil
}

func (Database *Database) GetuserNickname(Friendid string) (string, error) {
	var Nickname string
	err := Database.Db.QueryRow(`SELECT Nickname FROM user  WHERE id = ?`, Friendid).Scan(&Nickname)
	if err != nil {
		return Nickname, err
	}

	return Nickname, nil
}