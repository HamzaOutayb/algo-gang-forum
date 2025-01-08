package service

import (
	"errors"
	"net/mail"
	"strings"

	"real-time-forum/internal/models"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) LoginUser(user *models.User) error {
	// email
	(*user).Email = strings.ToLower((*user).Email)
	if EmailChecker((*user).Email) {
		return errors.New(models.Errors.InvalidEmail)
	}
	if len((*user).Email) > 50 {
		return errors.New(models.Errors.LongEmail)
	}
	// Password
	if len((*user).Password) < 6 || len((*user).Password) > 30 {
		return errors.New(models.Errors.InvalidPassword)
	}
	// check existance
	if !s.Database.CheckIfUserExists((*user).Username, (*user).Email) {
		return errors.New(models.Errors.InvalidCredentials)
	}
	// get user password
	UserPassword, err := s.Database.GetUserPassword((*user).Email)
	if err != nil {
		return err
	}

	// Check Password Validity
	if !CheckPasswordValidity(UserPassword, (*user).Password) {
		return errors.New(models.Errors.InvalidCredentials)
	}

	// generate new uuid
	(*user).Uuid = GenerateUuid()

	// Update uuid
	s.Database.UpdateUuid((*user).Uuid, (*user).Email)

	return nil
}

func EmailChecker(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func CheckPasswordValidity(hashedPass, entredPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(entredPass))
	return err == nil
}

func GenerateUuid() string {
	return uuid.Must(uuid.NewV4()).String()
}
