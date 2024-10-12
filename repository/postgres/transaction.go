package postgres

import (
	"financeapp/aggregate"
	"financeapp/domain/transaction"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TransactionRepo struct {
	db *sqlx.DB
}

func NewTransactionRepo(db *sqlx.DB) *TransactionRepo {
	return &TransactionRepo{
		db: db,
	}
}

// TODO: Add(*transaction.Transaction) (*transaction.Transaction, error)
func (tr TransactionRepo) Add(t *transaction.Transaction) (*transaction.Transaction, error) {
	tx, err := tr.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	_, err = tx.Exec(
		"INSERT INTO transactions (id, user_id, category_id, currency, value, type) VALUES($1, $2, $3, $4, $5, $6)",
		t.ID,
		t.UserID,
		t.CategoryID,
		t.Currency,
		t.Value,
		t.TransactionType,
	)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return t, nil
}

// TODO: GetByUserId(userID *uuid.UUID) (*[]aggregate.Transaction, error)
func (tr TransactionRepo) GetByUserId(userID uuid.UUID) (*[]aggregate.Transaction, error) {
	rows, err := tr.db.Queryx(
		"SELECT t.id, t.user_id, t.category_id, t.currency, t.value, t.type, c.id, c.user_id, c.Name FROM transactions AS t JOIN categories c on t.category_id=c.id WHERE t.user_id=$1",
		userID,
	)
	if err != nil {
		return nil, err
	}
	transactions := []aggregate.Transaction{}
	for rows.Next() {
		t := aggregate.Transaction{}
		rows.Scan(
			&t.Transaction.ID,
			&t.Transaction.UserID,
			&t.Transaction.CategoryID,
			&t.Transaction.Currency,
			&t.Transaction.Value,
			&t.Transaction.TransactionType,
			&t.Category.ID,
			&t.Category.UserID,
			&t.Category.Name,
		)
		transactions = append(transactions, t)
	}
	return &transactions, nil
}
