package category

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrorCategoryNotFound      = errors.New("Category not found")
	ErrorDuplicateCategoryName = errors.New("Categoy with name %s already exists")
)

type Repo interface {
	Add(c *Category) (*Category, error)
	GetByUserID(userID uuid.UUID) ([]Category, error)
	GetByID(userID, categoryID uuid.UUID) (*Category, error)
}
