package utils

import (
	"context"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func RenderHtml(t templ.Component) echo.HandlerFunc {
	return func(c echo.Context) error {
		return t.Render(context.Background(), c.Response().Writer)
	}
}

func WriteHTML(c echo.Context, statusCode int, t templ.Component) error {
	c.Response().Header().Add("content-type", "text/html")
	c.Response().WriteHeader(statusCode)
	return t.Render(context.Background(), c.Response().Writer)
}
