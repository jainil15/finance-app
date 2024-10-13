package routes

import (
	"context"
	"financeapp/web/views"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func RenderHtml(t templ.Component) echo.HandlerFunc {
	return func(c echo.Context) error {
		return t.Render(context.Background(), c.Response().Writer)
	}
}

func AddViews(e *echo.Echo) {
	e.GET("/", RenderHtml(views.Home()))
}
