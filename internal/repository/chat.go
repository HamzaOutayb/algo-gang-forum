package repository

import "real-time-forum/internal/models"

func (Database *Database) InsertChat(From, To int, Message []byte) error {
	var Conversations_ID int64
	Database.Db.QueryRow("SELECT id FROM conversations WHERE (user_one = ? AND user_two = ?) OR (user_one = ? AND user_two = ?)", From, To, To, From).Scan(&Conversations_ID)
	if Conversations_ID == 0 {
		Insertchat, err := Database.Db.Exec("INSERT INTO conversations (user_one, user_two) VALUES (?, ?)", From, To)
		if err != nil {
			return err
		}
		Conversations_ID, err = Insertchat.LastInsertId()
		if err != nil {
			return err
		}
	}
	_, err := Database.Db.Exec("INSERT INTO messages (sender_id, content, conversation_id) VALUES (?, ?, ?)", From, string(Message), Conversations_ID)
	if err != nil {
		return err
	}
	return nil
}

const pagesize = 10

func (Database *Database) GetChatWith(pagenm int, usrid int) ([]models.Chat, error) {
	start := pagenm * pagesize
	rows, err := Database.Db.Query(`
	SELECT DISTINCT
  id,
  CASE 
    WHEN user_one = ${1} THEN user_two
    ELSE user_one
  END AS chatted_user
FROM chats
WHERE ${1} IN (user_one, user_two) OFFSET ${2} LIMIT 15 ORDER BY id;
`, usrid, start)
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

func (Database *Database) GetChat(pagenm int, usrid int) ([]models.Chat, error) {
	start := pagenm * pagesize
	rows, err := Database.Db.Query(`SELECT Nickname
	FROM user
	WHERE id NOT IN (SELECT id FROM conversation WHERE user_one = ${1} OR user_two  = ${2}) OFFSET ? ORDERBY id;
`, usrid, start)
	if err != nil {
		return []models.Chat{}, err
	}

	var Chats []models.Chat
	for rows.Next() {
		var chat models.Chat
		err := rows.Scan(&chat.FriendsId)
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
