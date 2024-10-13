package transactionService

import (
	"financeapp/domain/category"
	"financeapp/domain/transaction"
	errx "financeapp/pkg/errors"
	"financeapp/pkg/middleware"
	"financeapp/pkg/utils"
	"financeapp/repository/repository"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TransactionService struct {
	transactionRepo repository.TransactionRepo
	categoryRepo    category.Repo
}

func NewTransactionService(tr repository.TransactionRepo, cr category.Repo) *TransactionService {
	return &TransactionService{
		transactionRepo: tr,
		categoryRepo:    cr,
	}
}

func NewTransactionRoutes(g *echo.Group, ts *TransactionService) {
	g.POST("/user/:user_id/transaction", ts.Add, middleware.AuthMiddleware, middleware.CheckUser)
	g.GET(
		"/user/:user_id/transactions",
		ts.GetByUser,
		middleware.AuthMiddleware,
		middleware.CheckUser,
	)
}

type transactionRequest struct {
	UserID          uuid.UUID `param:"user_id"`
	CategoryID      uuid.UUID `                json:"category_id"      form:"category_id"`
	TransactionType string    `                json:"transaction_type" form:"transaction_type"`
	Currency        string    `                json:"currency"         form:"currency"`
	Value           float64   `                json:"value"            form:"value"`
}
type transactionResponse struct {
	ID              uuid.UUID                   `json:"id"`
	UserID          uuid.UUID                   `json:"user_id"`
	CategoryID      uuid.UUID                   `json:"category_id"`
	TransactionType transaction.TransactionType `json:"transaction_type"`
	Currency        transaction.Currency        `json:"currency"`
	Value           transaction.Value           `json:"value"`
}

func ToTransactinResponse(t *transaction.Transaction) *transactionResponse {
	return &transactionResponse{
		ID:              t.ID,
		UserID:          t.UserID,
		CategoryID:      t.CategoryID,
		TransactionType: t.TransactionType,
		Currency:        t.Currency,
		Value:           t.Value,
	}
}

func (ts TransactionService) GetByUser(c echo.Context) error {
	uID := c.Param("user_id")
	userID, err := uuid.Parse(uID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Error{
			Message: "Invalid user id",
		})
	}
	trans, err := ts.transactionRepo.GetByUserId(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.Error{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, trans)
}

func (ts TransactionService) Add(c echo.Context) error {
	var t transactionRequest
	err := c.Bind(&t)
	if err != nil {
		return err
	}
	errs := errx.New()
	newID := uuid.New()
	tt, err := transaction.NewTransactionType(t.TransactionType)
	if err != nil {
		errs.Add("transaction_type", err.Error())
	}
	fmt.Printf("%v\n", t)
	curr, err := transaction.NewCurrency(t.Currency)
	if err != nil {
		errs.Add("currency", err.Error())
	}
	val, err := transaction.NewValue(t.Value)
	if err != nil {
		errs.Add("value", err.Error())
	}
	if len(errs) > 0 {
		return c.JSON(http.StatusBadRequest, utils.Error{
			Message: "Bad Request",
			Error:   errs,
		})
	}
	cat, err := ts.categoryRepo.GetByID(t.UserID, t.CategoryID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Error{
			Message: "Category not found",
		})
	}
	if cat == nil {
		return c.JSON(http.StatusBadRequest, utils.Error{
			Message: "Category not found",
		})
	}
	fmt.Println("Error:", cat)
	tran := transaction.New(newID, t.UserID, t.CategoryID, curr, val, tt)
	tran, err = ts.transactionRepo.Add(tran)
	log.Println(tran)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.Error{
			Message: "Internal server error",
			Error:   err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, utils.Response{
		Message: "Transaction created",
		Result:  ToTransactinResponse(tran),
	})
}
