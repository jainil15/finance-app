package routes

import (
	"context"
	"financeapp/pkg/middleware"
	"financeapp/pkg/model"
	"financeapp/repository/postgres"
	userService "financeapp/service/user"
	"financeapp/web/components/forms"
	"financeapp/web/components/home"
	"financeapp/web/layout"

	"github.com/a-h/templ"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func RenderHtml(t templ.Component) echo.HandlerFunc {
	return func(c echo.Context) error {
		// c.Response().Header().Add("HX-Redirect", "/register")
		return t.Render(context.Background(), c.Response().Writer)
	}
}

func AddViews(e *echo.Echo, db *sqlx.DB) {
	e.GET("/", RenderHtml(layout.Layout(layout.Navbar(), home.Home())))
	e.GET("/login", RenderHtml(layout.Layout(layout.Navbar(), home.Home())))
	e.GET(
		"/register",
		RenderHtml(layout.Layout(layout.Navbar(), forms.Register(model.RegisterUser{}, nil))),
	)
	e.GET("/fragment/login", RenderHtml(home.Home()))
	e.GET(
		"/fragment/register",
		RenderHtml(forms.Register(model.RegisterUser{}, nil)),
	)

	us := userService.NewUserService(
		postgres.NewUserRepo(db),
		postgres.NewAccountRepo(db),
		postgres.NewTransactionRepo(db),
		postgres.NewCategoryRepo(db),
		postgres.NewBudgetRepo(db),
	)
	e.GET(
		"/fragment/user/:user_id/transaction-form",
		us.TransactionForm,
		middleware.AuthMiddleware,
		middleware.CheckUser,
	)
	e.GET("/home", us.GetUserInfoView, middleware.AuthMiddleware)
}
