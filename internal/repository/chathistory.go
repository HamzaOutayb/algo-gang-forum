package repository

import "fmt"

type Conversations struct {
	Sender     string
	Content    string
	Created_at string
}

func (Database *Database) HistoryMessages(from, to int) ([]Conversations, error) {
	conversations_id := 0
	var result []Conversations
	Database.Db.QueryRow("SELECT id FROM conversations WHERE (user_one = ? AND user_two = ?) OR (user_two = ? AND user_one = ?)", from, to, from, to).Scan(&conversations_id)
	rows, err := Database.Db.Query("SELECT u.Nickname,m.content,m.created_at FROM messages m JOIN user u ON m.sender_id = u.id WHERE conversation_id = ? ", conversations_id)
	fmt.Println(from,to,conversations_id)
	if err != nil {
		return []Conversations{}, err
	}
	for rows.Next() {
		var sender, message, date string
		err := rows.Scan(&sender, &message, &date)
		if err != nil {
			return nil, err
		}
		var messages = Conversations{
			Sender:     sender,
			Content:    message,
			Created_at: date,
		}
		result = append(result, messages)
	}
	fmt.Println(result)
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil

}
