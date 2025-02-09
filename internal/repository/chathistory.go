package repository

import (
	"errors"
	"math"

	"real-time-forum/internal/models"
)

const messagesperpage float64 = 10

func (Database *Database) HistoryMessages(from, to int, pagenm int) ([]models.Conversations, error) {
	conversations_id := 0
	var result []models.Conversations
	Database.Db.QueryRow("SELECT id FROM conversations WHERE (user_one = ? AND user_two = ?) OR (user_two = ? AND user_one = ?)", from, to, from, to).Scan(&conversations_id)
	count, err := Database.Getchatmessagescount(conversations_id)
	if err != nil {
		return []models.Conversations{}, err
	}

	floatpages := math.Ceil(float64(count) / messagesperpage)
	start := (pagenm - int(floatpages)) * int(messagesperpage)
	if start >= count || start < 0 {
		return []models.Conversations{}, errors.New(models.CommentErrors.InvalidPage)
	}
	rows, err := Database.Db.Query("SELECT u.Nickname,m.content,m.created_at FROM messages m JOIN user u ON m.sender_id = u.id WHERE conversation_id = ? ORDER BY m.created_at ASC LIMIT ? OFFSET ?", conversations_id, messagesperpage, start)
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

func (Database *Database) Getchatmessagescount(conversation_id int) (int, error) {
	var count int

	err := Database.Db.QueryRow("SELECT COUNT(*) FROM messages WHERE conversation_id = ?", conversation_id).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
