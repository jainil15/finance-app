package transaction

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrorInvalidCurrency        = errors.New("Invalid currency")
	ErrorEmptyValue             = errors.New("Value required")
	ErrorInvalidTransactionType = errors.New("Invalid transaction type")
	ErrorEmptyTransactionType   = errors.New("Transaction type required")
)

type (
	TransactionType string
	Currency        string
	Value           float64
)

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Transaction struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	CategoryID      uuid.UUID
	Currency        Currency
	Value           Value
	TransactionType TransactionType
}

func New(
	id, userID, categoryID uuid.UUID,
	currency Currency,
	value Value,
	transactionType TransactionType,
) *Transaction {
	return &Transaction{
		ID:              id,
		UserID:          userID,
		CategoryID:      categoryID,
		Currency:        currency,
		Value:           value,
		TransactionType: transactionType,
	}
}

func NewCurrency(c string) (Currency, error) {
	var err error
	if c != "INR" && c != "USD" {
		err = errors.Join(err, ErrorInvalidCurrency)
	}
	return Currency(c), err
}

func NewValue(v float64) (Value, error) {
	var err error
	if v == 0.0 {
		err = errors.Join(err, ErrorEmptyValue)
	}
	return Value(v), err
}

func NewTransactionType(t string) (TransactionType, error) {
	var err error
	if t != string(Income) && t != string(Expense) {
		fmt.Printf("Debug: %s\n", t)
		err = errors.Join(err, ErrorInvalidTransactionType)
	}
	return TransactionType(t), err
}
