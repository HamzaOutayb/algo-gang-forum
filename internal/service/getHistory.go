package service

import (
	"errors"

	"real-time-forum/internal/models"
)

const messagesperpage = 10

func (s *Service) GetHistory(user string, to string, pagenm int) ([]models.Conversations, error) {
	user_id, to_id, err := s.Database.GetId2(user, to)
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
