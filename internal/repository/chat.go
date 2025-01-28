package repository

import (
	"fmt"

	"real-time-forum/internal/models"
)

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

func (Database *Database) HistoryMessages(from, to int) ([]models.Conversations, error) {
	conversations_id := 0
	var result []models.Conversations
	Database.Db.QueryRow("SELECT id FROM conversations WHERE (user_one = ? AND user_two = ?) OR (user_two = ? AND user_one = ?)", from, to, from, to).Scan(&conversations_id)
	rows, err := Database.Db.Query("SELECT u.Nickname,m.content,m.created_at FROM messages m JOIN user u ON m.sender_id = u.id WHERE conversation_id = ? ", conversations_id)
	if err != nil {
		return []models.Conversations{}, err
	}
	for rows.Next() {
		var sender, message, date string
		err := rows.Scan(&sender, &message, &date)
		if err != nil {
			return nil, err
		}
		messages := models.Conversations{
			Sender:     sender,
			Content:    message,
			Created_at: date,
		}
		result = append(result, messages)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
