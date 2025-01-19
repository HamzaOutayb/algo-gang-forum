package service

import "real-time-forum/internal/models"

func (S *Service) GetHistory(from, to int) []models.Messagesbody {
	messages, err := S.Database.HistoryMessages(from, to)
	if err != nil {
		// err message
	}
	return messages
}
