package models

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Uuid     string `json:"uid"`
}


var UserErrors struct {
	InvalidEmail     string
	InvalidUsername  string
	InvalidPassword  string
	UserAlreadyExist string
	UserNotExist     string
} = struct {
	InvalidEmail     string
	InvalidUsername  string
	InvalidPassword  string
	UserAlreadyExist string
	UserNotExist     string
}{
	InvalidEmail:     "invalid email",
	InvalidUsername:  "invalid username",
	InvalidPassword:  "invalid password",
	UserAlreadyExist: "user already exist",
	UserNotExist:     "user doesn't exist",
}