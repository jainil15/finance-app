package budgetService

import (
	"errors"
	"financeapp/domain/budget"
	"financeapp/pkg/middleware"
	"financeapp/pkg/utils"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
)

type addBudgetRequest struct {
	UserID   uuid.UUID `param:"user_id"`
	Currency string    `json:"currency"`
	Value    float64   `json:"value"`
}

type BudgetService struct {
	budgetRepo budget.Repo
}

func NewBudgetService(budgetRepo budget.Repo) *BudgetService {
	return &BudgetService{
		budgetRepo: budgetRepo,
	}
}

type budgetRespose struct {
	UserID   uuid.UUID       `param:"user_id"`
	Currency budget.Currency `json:"currency"`
	Value    float64         `json:"value"`
}

func ToBudgetResponse(b *budget.Budget) *budgetRespose {
	return &budgetRespose{
		UserID:   b.UserID,
		Currency: b.Currency,
		Value:    b.Value,
	}
}

func NewBudgetRoutes(e *echo.Echo, bs *BudgetService) {
	e.GET(
		"user/:user_id/budget",
		bs.GetByUserID,
		middleware.AuthMiddleware,
		middleware.CheckUser,
	)
	e.POST("user/:user_id/budget", bs.Add, middleware.AuthMiddleware, middleware.CheckUser)
}

func (bs BudgetService) GetByUserID(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Error{
			Message: fmt.Sprintf("Invalid user id: %s", userID),
		})
	}
	b, err := bs.budgetRepo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, budget.ErrorBudgetNotFound) {
			return c.JSON(http.StatusNotFound, utils.Error{
				Message: "Budget not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, utils.Error{
			Error:   err.Error(),
			Message: "Internal server error",
		})
	}
	return c.JSON(http.StatusOK, ToBudgetResponse(b))
}

func (bs BudgetService) Add(c echo.Context) error {
	bRequest := addBudgetRequest{}
	err := c.Bind(&bRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Error{
			Message: "Bad request",
		})
	}
	b, errs := budget.NewBudget(
		bRequest.UserID,
		bRequest.Currency,
		bRequest.Value,
	)
	if errs != nil {
		return c.JSON(http.StatusBadRequest, utils.Error{
			Message: "Bad request",
			Error:   errs,
		})
	}
	b, err = bs.budgetRepo.Add(b.UserID, b)
	if err != nil {
		pgxError, ok := err.(*pgconn.PgError)
		if ok {
			switch pgxError.Code {
			case "23505":
				return c.JSON(http.StatusConflict, utils.Error{
					Message: "Budget already exists",
				})
			}
		}
		return c.JSON(http.StatusInternalServerError, utils.Error{
			Message: "Error adding budget",
		})
	}
	return c.JSON(http.StatusCreated, utils.Response{
		Message: "Budget Created",
		Result:  ToBudgetResponse(b),
	})
}
