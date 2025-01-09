package service

import (
	"errors"

	"real-time-forum/internal/models"
)

func (data *Service) CheckReactInput(react models.React) error {
	if (react.React != 2 && react.React != 1) || (react.Thread_type != "post" && react.Thread_type != "comment") || react.Thread_id < 0 || !data.Database.CheckIfThreadExists(react) {
		return errors.New("bad request")
	}
	return nil
}
