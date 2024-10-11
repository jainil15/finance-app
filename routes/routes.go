package routes

import (
	"financeapp/repository/postgres"
	service "financeapp/service/user"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func AddRoutes(e *echo.Echo, db *sqlx.DB) error {
	e.Use(middleware.Logger())

	userService := service.NewUserService(postgres.NewUserRepo(db), postgres.NewAccountRepo(db))
	service.NewUserRoutes(e, userService)
	return nil
}
