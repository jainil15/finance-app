package budget

import (
	"errors"

	"github.com/google/uuid"
)

var ErrorBudgetNotFound = errors.New("Budget not found")

type Repo interface {
	Add(userID uuid.UUID, budget *Budget) (*Budget, error)
	GetByUserID(userID uuid.UUID) (*Budget, error)
	Update(userID uuid.UUID, budget *Budget) (*Budget, error)
}
