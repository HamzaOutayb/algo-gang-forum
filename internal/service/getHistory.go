package service

import "real-time-forum/internal/models"

func (S *Service) GetHistory(from, to int) ([]models.Conversations, error) {
	messages, err := S.Database.HistoryMessages(from, to)
	if err != nil {
		return []models.Conversations{}, err
	}
	return messages, nil
}
