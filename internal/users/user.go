package users

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id               uuid.UUID `db:"id" json:"id"`
	RegistrationTime time.Time `db:"registration_time" json:"registration_time"`
	Password         string    `db:"password" json:"password"`
	Username         string    `db:"username" json:"username"`
	Avatar           string    `db:"avatar" json:"avatar"`
	Sex              string    `db:"sex" json:"sex"`
	Email            string    `db:"email" json:"email"`
}
