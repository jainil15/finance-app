package routes

import (
	"financeapp/repository/postgres"
	budgetService "financeapp/service/budget"
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
	return nil
}
