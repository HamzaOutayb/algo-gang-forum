package service

import (
	"errors"
	"html"
	"net/mail"
	"strconv"
	"strings"

	"real-time-forum/internal/models"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) LoginUser(user *models.User) error {
	// email
	(*user).Email = strings.ToLower((*user).Email)
	if !EmailChecker((*user).Email) {
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
	if !s.Database.CheckIfUserExists((*user).Nickname, (*user).Email) {
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

func (s *Service) RegisterUser(user *models.User) error {
	// Username
	if len((*user).Nickname) < 3 || len((*user).Nickname) > 15 {
		return errors.New(models.Errors.InvalidUsername)
	}

	// Age
	if !CheckAgeValidation((*user).Age) {
		return errors.New(models.UserErrors.InvalideAge)
	}

	//First_Name
	 if len((*user).First_Name) < 3 || len((*user).First_Name) > 15 {
		return errors.New(models.UserErrors.InvalideFirst_Name)
	}

	// Last_Name
	if len((*user).Last_Name) < 3 || len((*user).Last_Name) > 15 {
		return errors.New(models.UserErrors.InvalideLast_Name)
	}

	// Password
	if len((*user).Password) < 6 || len((*user).Password) > 30 {
		return errors.New(models.Errors.InvalidPassword)
	}

	// email
	(*user).Email = strings.ToLower((*user).Email)
	if EmailChecker((*user).Email) {
		return errors.New(models.Errors.InvalidEmail)
	}
	if len((*user).Email) > 50 {
		return errors.New(models.Errors.LongEmail)
	}

	// username or email existance
	if s.Database.CheckIfUserExists((*user).Nickname, (*user).Email) {
		return errors.New(models.Errors.UserAlreadyExist)
	}

	// Generate Uuid
	(*user).Uuid = GenerateUuid()
	
	// Encrypt Pass
	var err error
	(*user).Password, err = EncyptPassword((*user).Password)
	if err != nil {
		return err
	}

	// Fix username html
	(*user).Nickname = html.EscapeString((*user).Nickname)

	// Insert the user
	return s.Database.InsertUser(*user)
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

func EncyptPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPass), nil
}

func CheckAgeValidation(age string) bool {
	Age, err := strconv.Atoi(age)
	if err != nil {
		return false
	}
	if Age > 1000 || Age < 0 {
		return false
	}

	return true
}
