package service

import "real-time-forum/internal/models"

func (s *Service) GetLastconversations(page int, usrid int) ([]models.Chat, error) {
	chat, err := s.Database.GetChatAndUserID(page, usrid); if err != nil {
		return []models.Chat{}, err
	}
	return chat, nil
}
