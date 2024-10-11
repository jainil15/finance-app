package transaction

import (
	"github.com/google/uuid"
)

type (
	Cost float64
)

type Transaction struct {
	ID         uuid.UUID
	AccountID  uuid.UUID
	CategoryID uuid.UUID
	Cost       Cost
}

func New() {
}
