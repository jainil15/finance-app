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

func WriteHTMLWithStatus(c echo.Context, s int, t ...templ.Component) error {
	c.Response().Header().Add("content-type", "text/html")
	c.Response().WriteHeader(s)
	for _, t := range t {
		err := t.Render(context.Background(), c.Response().Writer)
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteHTML(c echo.Context, t ...templ.Component) error {
	c.Response().Header().Add("content-type", "text/html")
	c.Response().WriteHeader(200)
	for _, t := range t {
		err := t.Render(context.Background(), c.Response().Writer)
		if err != nil {
			return err
		}
	}
	return nil
}
