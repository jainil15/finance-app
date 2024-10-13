package categoryService

import (
	"errors"
	"financeapp/domain/category"
	errx "financeapp/pkg/errors"
	"financeapp/pkg/middleware"
	"financeapp/pkg/utils"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CategoryService struct {
	catergoryRepo category.Repo
}

func NewCategoryService(cr category.Repo) *CategoryService {
	return &CategoryService{
		catergoryRepo: cr,
	}
}

func NewCategoryRoutes(g *echo.Group, cs *CategoryService) {
	g.GET(
		"user/:user_id/categories",
		cs.GetByUser,
		middleware.AuthMiddleware,
		middleware.CheckUser,
	)
	g.POST("user/:user_id/categories", cs.Add, middleware.AuthMiddleware, middleware.CheckUser)
}

type catergoryRequest struct {
	UserID uuid.UUID `param:"user_id"`
	Name   string    `                json:"name" form:"name"`
}
type catergoryRespone struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
}

func ToCategoryRespose(c *category.Category) *catergoryRespone {
	return &catergoryRespone{
		ID:     c.ID,
		UserID: c.UserID,
		Name:   string(c.Name),
	}
}

func (cs CategoryService) Add(c echo.Context) error {
	var categoryReq catergoryRequest
	if err := c.Bind(&categoryReq); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Error{
			Message: "Bad request",
		})
	}
	errs := errx.New()
	newID := uuid.New()
	name, err := category.NewName(categoryReq.Name)
	if err != nil {
		errs.Add("name", err.Error())
	}
	if len(errs) > 0 {
		return c.JSON(http.StatusBadRequest, utils.Error{
			Message: "Bad request",
			Error:   errs,
		})
	}
	cat := category.New(newID, categoryReq.UserID, name)
	cat, err = cs.catergoryRepo.Add(cat)
	if err != nil {
		if errors.Is(err, category.ErrorDuplicateCategoryName) {
			return c.JSON(http.StatusConflict, utils.Error{
				Error: map[string]string{
					"name": fmt.Sprintf("Category with name %s already exists", name),
				},
				Message: fmt.Sprintf("Category with name %s already exists", name),
			})
		}
		return c.JSON(http.StatusInternalServerError, utils.Error{
			Message: "Internal server error",
		})
	}
	return c.JSON(http.StatusCreated, utils.Response{
		Result:  ToCategoryRespose(cat),
		Message: "Category created",
	})
}

func (cs CategoryService) GetByUser(c echo.Context) error {
	uID := c.Param("user_id")
	userID, err := uuid.Parse(uID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Error{
			Message: "Invalid user id",
		})
	}
	categories, err := cs.catergoryRepo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, category.ErrorCategoryNotFound) {
			return c.JSON(http.StatusNotFound, utils.Error{
				Message: "Catergories not found for user",
			})
		}
		return c.JSON(http.StatusInternalServerError, utils.Error{
			Error:   err.Error(),
			Message: "Internal server error",
		})
	}
	catRes := make([]catergoryRespone, len(categories))
	for i, cat := range categories {
		catRes[i] = *ToCategoryRespose(&cat)
	}
	return c.JSON(http.StatusOK, utils.Response{
		Result: catRes,
	})
}
