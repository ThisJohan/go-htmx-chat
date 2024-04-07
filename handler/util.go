package handler

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

const (
	sessionTokenCookie = "session"
)

// This custom Render replaces Echo's echo.Context.Render() with templ's templ.Component.Render().
func render(ctx echo.Context, t templ.Component, status int) error {
	ctx.Response().Writer.WriteHeader(status)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}

func writeCookie(c echo.Context, key, value string) {
	cookie := new(http.Cookie)
	cookie.Name = key
	cookie.Value = value
	c.SetCookie(cookie)
}

func readCookie(c echo.Context, key string) (*http.Cookie, error) {
	cookie, err := c.Cookie(key)
	if err != nil {
		return nil, err
	}
	return cookie, nil
}

func deleteCookie(c echo.Context, name string) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = ""
	cookie.MaxAge = -1
	c.SetCookie(cookie)
}
