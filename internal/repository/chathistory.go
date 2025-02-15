package repository

import (
	"errors"
	"fmt"

	"real-time-forum/internal/models"
)

const messagesperpage = 10

func (Database *Database) HistoryMessages(from int, to string, pagenm int) ([]models.Conversations, error) {
	conversations_id := 0
	var result []models.Conversations
	err := Database.Db.QueryRow("SELECT id FROM conversations WHERE (user_one = ? AND user_two = ?) OR (user_two = ? AND user_one = ?)", from, to, from, to).Scan(&conversations_id)
	if err != nil {
		return []models.Conversations{}, err
	}

	count, err := Database.Getchatmessagescount(conversations_id)
	if err != nil {
		return []models.Conversations{}, err
	}
	start := (pagenm - 1) * messagesperpage
	if start >= count || start < 0 {
		return []models.Conversations{}, errors.New(models.CommentErrors.InvalidPage)
	}
	fmt.Println("start", start)
	rows, err := Database.Db.Query(`
    SELECT u.Nickname, m.content, m.created_at 
    FROM messages m 
    JOIN user u ON m.sender_id = u.id 
    WHERE m.conversation_id = ? 
    ORDER BY m.id DESC 
    LIMIT ? OFFSET ?`, conversations_id, messagesperpage, start)
	if err != nil {
		return []models.Conversations{}, err
	}
	defer rows.Close()

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

func (Database *Database) Getchatmessagescount(conversation_id int) (int, error) {
	var count int

	err := Database.Db.QueryRow("SELECT COUNT(*) FROM messages WHERE conversation_id = ?", conversation_id).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
