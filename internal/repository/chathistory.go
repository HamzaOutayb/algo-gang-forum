package repository

import (
	"real-time-forum/internal/models"
)

func (Database *Database) HistoryMessages(from, to int) ([]models.Messagesbody, error) {

	conversations_id := 0

	Database.Db.QueryRow("SELECTE id FROM conversations WHERE user_one = ? AND user_two = ?", from, to).Scan(conversations_id)
	rows, err := Database.Db.Query("SELECTE sender_id,content,created_at FROM messages WHERE conversation_id = ?", from, to, conversations_id)
	if err != nil {
		return []models.Messagesbody{}, err
	}

	conversations := []models.Messagesbody{}
	conversation := models.Messagesbody{}
	for rows.Next() {
		var sender, message, date string
		err := rows.Scan(&sender, &message, &date)
		if err != nil {
			return nil, err
		}

		conversation.Sender = sender
		conversation.Message = message
		conversation.Date = date
		conversations = append(conversations, conversation)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return conversations, nil
}