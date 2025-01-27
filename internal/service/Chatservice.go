package service

import "real-time-forum/internal/models"

func (s *Service) GetLastconversations(page int, usrid int) ([]models.Chat, error) {
	chat, err := s.Database.GetChatWith(page, usrid)
	if err != nil {
		return []models.Chat{}, err
	}
	return chat, nil
}

func (s *Service) Getconversations(page int, usrid int) ([]models.Chat, error) {
	chat, err := s.Database.GetChat(page, usrid)
	if err != nil {
		return []models.Chat{}, err
	}
	return chat, nil
}
