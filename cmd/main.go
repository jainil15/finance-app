package main

import (
	"financeapp/pkg/config"
	"financeapp/repository/postgres"
	"financeapp/routes"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Run() {
	env := config.Env
	db := postgres.New(env.DBUser, env.DBPassword, env.DBHost, env.DBPort, env.DBName)
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusTeapot, map[string]string{
			"Hello": "World",
		})
	})
	routes.AddRoutes(e, db)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", env.Port)))
}

func main() {
	Run()
}
