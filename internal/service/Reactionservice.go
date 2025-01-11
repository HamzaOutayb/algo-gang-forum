package service

import (
	"errors"

	"real-time-forum/internal/models"
)

func (s *Service) CheckReactInput(react models.React) error {
	if (react.React != 2 && react.React != 1) || (react.Thread_type != "post" && react.Thread_type != "comment") || react.Thread_id < 0 || !s.Database.CheckIfThreadExists(react) {
		return errors.New("bad request")
	}
	return nil
}

func (s *Service) LikesTotal(thread_type string, thread_id int) (models.ReactResponse, error) {
	response := models.ReactResponse{}
	var err error
	if thread_type == "post" {
		response.Like, response.Dislike, err = s.Database.CountPostLikes(thread_id)
	} else {
		response.Like, response.Dislike, err = s.Database.CountCommentLikes(thread_id)
	}
	if err != nil {
		return response, err
	}
	return response, nil
}

func (s *Service) GetLikedThread(thread_type string, thread_id, user_id int) (bool, bool) {
	isLiked, isDisliked := false, false
	if thread_type == "post" {
		isLiked, isDisliked = s.Database.CheckIfLikedPost(thread_id, user_id)
	} else {
		isLiked, isDisliked = s.Database.CheckIfLikedComment(thread_id, user_id)
	}
	return isLiked, isDisliked
}