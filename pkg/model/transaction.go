package model

import (
	"financeapp/domain/transaction"

	"github.com/google/uuid"
)

type Transaction struct {
	CategoryID      uuid.UUID
	Currency        transaction.Currency
	Value           transaction.Value
	TransactionType transaction.TransactionType
}
