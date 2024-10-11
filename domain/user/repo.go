package user

import (
	"errors"

	"github.com/google/uuid"
)

var ErrorUserNotFound = errors.New("User not found")

type Repo interface {
	Add(user *User) (*User, error)
	GetAll() ([]User, error)
	GetById(userID uuid.UUID) (*User, error)
	GetByEmail(email Email) (*User, error)
}
