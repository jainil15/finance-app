package category

import (
	"errors"

	"github.com/google/uuid"
)

var ErrorEmptyName = errors.New("Name is empty")

type Category struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Name   Name
}

type Name string

func NewName(n string) (Name, error) {
	if len(n) == 0 {
		return "", ErrorEmptyName
	}
	return Name(n), nil
}

func New(id uuid.UUID, userID uuid.UUID, name Name) *Category {
	return &Category{
		ID:     id,
		UserID: userID,
		Name:   name,
	}
}

func (n Name) String() string {
	return string(n)
}
