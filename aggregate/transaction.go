package aggregate

import (
	"financeapp/domain/category"
	"financeapp/domain/transaction"
)

type Transaction struct {
	Transaction transaction.Transaction
	Category    category.Category
}

func NewTransaction(t transaction.Transaction, c category.Category) Transaction {
	return Transaction{
		t, c,
	}
}
