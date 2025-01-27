package repository

import "real-time-forum/internal/models"

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
