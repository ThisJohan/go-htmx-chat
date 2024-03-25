package handler

import (
	"github.com/ThisJohan/go-htmx-chat/models"
	views "github.com/ThisJohan/go-htmx-chat/views/user"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserService *models.UserService
}

func (h *UserHandler) ShowUser(c echo.Context) error {
	return render(c, views.Show(), 200)
}
