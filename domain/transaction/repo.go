package transaction

import (
	"github.com/google/uuid"
)

type Repo interface {
	Add(*Transaction) (*Transaction, error)
	GetByUserId(userID *uuid.UUID) (*Transaction, error)
}
