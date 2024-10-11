package account

import "github.com/google/uuid"

type Cost struct {
	Currency string
	Value    float64
}
type Account struct {
	ID     uuid.UUID
	UserID uuid.UUID
	// Budget Cost
}

func NewBudget(c string, v float64) (*Cost, error) {
	return &Cost{
		Currency: c,
		Value:    v,
	}, nil
}

func New(userID uuid.UUID) *Account {
	return &Account{
		ID:     uuid.New(),
		UserID: userID,
		// Budget: budget,
	}
}
