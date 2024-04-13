package context

import (
	"fmt"

	"github.com/ThisJohan/go-htmx-chat/models"
	"github.com/labstack/echo/v4"
)

const (
	userKey = "app_user"
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

func (ctx *AppContext) WithUser(data *models.UserCache) {
	ctx.Set(userKey, data)
}

func (ctx *AppContext) User() (*models.UserCache, error) {
	user, ok := ctx.Get(userKey).(*models.UserCache)
	if !ok {
		return nil, fmt.Errorf("user context not found")
	}
	return user, nil
}
