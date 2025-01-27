package repository

import (
	"fmt"

	"real-time-forum/internal/models"
)

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
	rows, err := Database.Db.Query(`
	SELECT
    id,
    CASE 
        WHEN user_one = ? THEN user_two
        ELSE user_one
    END
FROM conversations
WHERE ? IN (user_one, user_two) 
ORDER BY created_at;
`, usrid, usrid)
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
	// start := pagenm * pagesize
	rows, err := Database.Db.Query(`SELECT id Nickname
FROM user
WHERE id != ?
AND id NOT IN (
    SELECT 
        CASE
            WHEN c.user_one = ? THEN c.user_two
            WHEN c.user_two = ? THEN c.user_one
        END
    FROM conversations c
    WHERE ? IN (c.user_one, c.user_two)
) ORDER BY Nickname;
`, usrid, usrid, usrid, usrid)
	if err != nil {
		fmt.Println("er", err)
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
