package server

import (
	"fmt"
	"github.com/google/uuid"
	"soa_hw_4/internal/users"
)

type UserBody struct {
	Password string `json:"password"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Sex      string `json:"sex"`
	Email    string `json:"email"`
}

type AuthUserBody struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type RegistrationResponce struct {
	Id    uuid.UUID `json:"id"`
	Token string    `json:"token"`
}

type CreateStatsTaskResponce struct {
	Id uuid.UUID `json:"id"`
}

type SeveralUserBody struct {
	Users []*users.User `json:"users"`
}

func Validate(userBody *UserBody) error {
	if userBody.Username == "" {
		return fmt.Errorf("username is empty")
	}
	if userBody.Sex != "male" && userBody.Sex != "female" {
		return fmt.Errorf("only male and female sex accepted")
	}
	return nil
}
