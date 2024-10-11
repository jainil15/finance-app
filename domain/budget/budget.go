package budget

import (
	"errors"
	errx "financeapp/pkg/errors"

	"github.com/google/uuid"
)

var ErrorInvalidCurrency = errors.New("Currency is not valid")

type Budget struct {
	UserID   uuid.UUID
	Currency Currency
	Value    float64
}

type Currency string

func NewCurrency(c string) (Currency, error) {
	if c != "INR" && c != "USD" {
		return "", ErrorInvalidCurrency
	}
	return Currency(c), nil
}

func NewBudget(userID uuid.UUID, c string, v float64) (*Budget, errx.Error) {
	errs := errx.New()
	curr, err := NewCurrency(c)
	if c == "" {
		errs.Add("currency", "Currency Required")
	}
	if v == 0.0 {
		errs.Add("value", "Value Required")
	}
	if err != nil {
		errs.Add("currency", err.Error())
	}
	if len(errs) > 0 {
		return nil, errs
	}
	return &Budget{
		userID, curr, v,
	}, nil
}
