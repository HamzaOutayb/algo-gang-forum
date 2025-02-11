package service

import (
	"errors"

	"real-time-forum/internal/models"
)

func (s *Service) GetHistory(user_uid string, to_id string, pagenm int) ([]models.Conversations, error) {
	user_id, err := s.Database.GetId2(user_uid)
	if err != nil {
		return []models.Conversations{}, err
	}
	if pagenm < 1 {
		return []models.Conversations{}, errors.New(models.CommentErrors.InvalidPage)
	}

	HistoryMessages, err := s.Database.HistoryMessages(user_id, to_id, pagenm)
	if err != nil {
		return []models.Conversations{}, err
	}
	return HistoryMessages, nil
}
