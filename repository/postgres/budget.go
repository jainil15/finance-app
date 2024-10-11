package postgres

import (
	"database/sql"
	"errors"
	"financeapp/domain/budget"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type BudgetRepo struct {
	db *sqlx.DB
}

func NewBudgetRepo(db *sqlx.DB) *BudgetRepo {
	return &BudgetRepo{
		db: db,
	}
}

type pgBudget struct {
	UserID   uuid.UUID       `db:"user_id"`
	Currency budget.Currency `db:"currency"`
	Value    float64         `db:"value"`
}

func ToPgBudget(b *budget.Budget) *pgBudget {
	return &pgBudget{
		UserID:   b.UserID,
		Currency: b.Currency,
		Value:    b.Value,
	}
}

func ToBudget(b *pgBudget) *budget.Budget {
	return &budget.Budget{
		UserID:   b.UserID,
		Currency: b.Currency,
		Value:    b.Value,
	}
}

func (br BudgetRepo) Add(userID uuid.UUID, b *budget.Budget) (*budget.Budget, error) {
	tx, err := br.db.Beginx()
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(
		"INSERT INTO budgets (user_id, currency, value) VALUES ($1, $2, $3)",
		userID,
		b.Currency,
		b.Value,
	)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return b, nil
}

func (br BudgetRepo) GetByUserID(userID uuid.UUID) (*budget.Budget, error) {
	pb := pgBudget{}
	err := br.db.Get(
		&pb,
		"SELECT user_id, currency, value FROM budgets WHERE user_id=$1",
		userID,
	)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, budget.ErrorBudgetNotFound
		}
		return nil, err
	}
	return ToBudget(&pb), nil
}

func (br BudgetRepo) Update(userID uuid.UUID, b *budget.Budget) (*budget.Budget, error) {
	tx, err := br.db.Beginx()
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(
		"UPDATE budgets SET currency=$1, value=$2 WHERE user_id=$3",
		b.Currency,
		b.Value,
		userID,
	)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return b, nil
}
