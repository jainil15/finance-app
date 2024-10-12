package routes

import (
	"financeapp/repository/postgres"
	budgetService "financeapp/service/budget"
	categoryService "financeapp/service/category"
	transactionService "financeapp/service/transaction"
	userService "financeapp/service/user"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func AddRoutes(e *echo.Echo, db *sqlx.DB) error {
	e.Use(middleware.Logger())

	us := userService.NewUserService(postgres.NewUserRepo(db), postgres.NewAccountRepo(db))
	userService.NewUserRoutes(e, us)

	bs := budgetService.NewBudgetService(postgres.NewBudgetRepo(db))
	budgetService.NewBudgetRoutes(e, bs)

	cs := categoryService.NewCategoryService(postgres.NewCategoryRepo(db))
	categoryService.NewCategoryRoutes(e, cs)

	ts := transactionService.NewTransactionService(postgres.NewTransactionRepo(db))
	transactionService.NewTransactionRoutes(e, ts)
	return nil
}
