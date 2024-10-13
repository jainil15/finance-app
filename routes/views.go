package routes

import (
	"context"
	"financeapp/pkg/model"
	"financeapp/web/components/forms"
	"financeapp/web/components/home"
	"financeapp/web/layout"
	"financeapp/web/views"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func RenderHtml(t templ.Component) echo.HandlerFunc {
	return func(c echo.Context) error {
		// c.Response().Header().Add("HX-Redirect", "/register")
		return t.Render(context.Background(), c.Response().Writer)
	}
}

func AddViews(e *echo.Echo) {
	e.GET("/", RenderHtml(views.Home()))
	e.GET("/login", RenderHtml(layout.Layout(home.Home())))
	e.GET("/register", RenderHtml(layout.Layout(forms.Register(model.RegisterUser{}, nil))))
	e.GET("/fragment/login", RenderHtml(home.Home()))
	e.GET(
		"/fragment/register",
		RenderHtml(forms.Register(model.RegisterUser{}, nil)),
	)
}
