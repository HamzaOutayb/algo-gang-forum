package service

import "real-time-forum/internal/models"

func (S *Service) GetHistory(from, to int) ([]models.Messagesbody, error) {
	messages, err := S.Database.HistoryMessages(from, to)
	if err != nil {
	return []models.Messagesbody{},err
	}
	return messages, nil
}