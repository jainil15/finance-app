package postgres

import (
	"errors"
	"financeapp/domain/user"

	"github.com/jmoiron/sqlx"
)

var ErrorUserNotFound = errors.New("User not found")

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (u UserRepo) Add(*user.User) (*user.User, error) {
	return nil, nil
}

func (u UserRepo) GetAll() ([]*user.User, error) {
	return nil, nil
}
