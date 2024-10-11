package account

import "github.com/google/uuid"

type Account struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func New(userID uuid.UUID) *Account {
	return &Account{
		ID:     uuid.New(),
		UserID: userID,
	}
}
