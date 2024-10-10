package postgres

import (
	"financeapp/domain/user"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	*sqlx.DB
}

func (u UserRepo) Add(*user.User) (*user.User, error) {
	return nil, nil
}

func (u UserRepo) GetAll() ([]*user.User, error) {
	return nil, nil
}
