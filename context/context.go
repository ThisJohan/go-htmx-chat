package context

import (
	"context"

	"github.com/ThisJohan/go-htmx-chat/models"
	"github.com/labstack/echo/v4"
)

type key = string

var (
	userKey key = "app_user"
)

type AppContext struct {
	echo.Context
}

func RegisterAppContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := &AppContext{c}
		return next(cc)
	}
}

func (ctx *AppContext) View() context.Context {
	newCtx := context.WithValue(ctx.Request().Context(), userKey, ctx.Get(userKey))

	return newCtx
}

func (ctx *AppContext) WithUser(data *models.UserCache) {
	ctx.Set(userKey, data)
}

func (ctx *AppContext) User() *models.UserCache {
	return ctx.Get(userKey).(*models.UserCache)
}

func User(ctx context.Context) *models.UserCache {
	return ctx.Value(userKey).(*models.UserCache)
}
