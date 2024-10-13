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
	e.Static("/static", "web/static")
	g := e.Group("/api")
	us := userService.NewUserService(
		postgres.NewUserRepo(db),
		postgres.NewAccountRepo(db),
		postgres.NewTransactionRepo(db),
		postgres.NewCategoryRepo(db),
		postgres.NewBudgetRepo(db),
	)
	userService.NewUserRoutes(g, us)

	bs := budgetService.NewBudgetService(postgres.NewBudgetRepo(db))
	budgetService.NewBudgetRoutes(g, bs)

	cs := categoryService.NewCategoryService(postgres.NewCategoryRepo(db))
	categoryService.NewCategoryRoutes(g, cs)

	ts := transactionService.NewTransactionService(
		postgres.NewTransactionRepo(db),
		postgres.NewCategoryRepo(db),
	)
	transactionService.NewTransactionRoutes(g, ts)
	return nil
}
