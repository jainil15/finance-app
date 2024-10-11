package postgres

import (
	"financeapp/domain/account"

	"github.com/jmoiron/sqlx"
)

type AccountRepo struct {
	db *sqlx.DB
}

func NewAccountRepo(db *sqlx.DB) *AccountRepo {
	return &AccountRepo{
		db: db,
	}
}

func (a *AccountRepo) Add(newAccount *account.Account) (*account.Account, error) {
	tx, err := a.db.Beginx()
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(
		"INSERT INTO accounts (id, user_id) VALUES ($1, $2)",
		newAccount.ID,
		newAccount.UserID,
	)
	if err != nil {
		return nil, err
	}
	return newAccount, nil
}
