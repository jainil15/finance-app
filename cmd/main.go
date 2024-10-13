package main

import (
	"financeapp/pkg/config"
	"financeapp/repository/postgres"
	"financeapp/routes"
	"fmt"

	"github.com/labstack/echo/v4"
)

func Run() {
	env := config.Env
	db := postgres.New(env.DBUser, env.DBPassword, env.DBHost, env.DBPort, env.DBName)
	e := echo.New()
	routes.AddViews(e, db)
	routes.AddRoutes(e, db)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", env.Port)))
}

func main() {
	Run()
}
