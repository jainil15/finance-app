package repository

import (
	"financeapp/aggregate"
	"financeapp/domain/transaction"

	"github.com/google/uuid"
)

type TransactionRepo interface {
	Add(*transaction.Transaction) (*transaction.Transaction, error)
	GetByUserId(userID *uuid.UUID) (*[]aggregate.Transaction, error)
}
