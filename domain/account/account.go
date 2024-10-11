package account

import (
	"github.com/google/uuid"
)

type Budget struct {
	UserID   uuid.UUID
	Currency string
	Value    float64
}

type Account struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func NewBudget(userID uuid.UUID, c string, v float64) (*Budget, error) {
	return &Budget{
		userID, c, v,
	}, nil
}

func New(userID uuid.UUID) *Account {
	return &Account{
		ID:     uuid.New(),
		UserID: userID,
	}
}
