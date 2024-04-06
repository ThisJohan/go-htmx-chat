package handler

import (
	"context"
	"fmt"

	"github.com/ThisJohan/go-htmx-chat/models"
	views "github.com/ThisJohan/go-htmx-chat/views/user"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserService    *models.UserService
	SessionService *models.SessionService
}

func (h *UserHandler) ShowUser(c echo.Context) error {
	return render(c, views.Show(), 200)
}

func (h *UserHandler) Signup(c echo.Context) error {
	return render(c, views.Signup(), 200)
}

func (h *UserHandler) ProcessSignup(c echo.Context) error {
	var data models.CreateUserDTO
	if err := c.Bind(&data); err != nil {
		return err
	}
	_, err := h.UserService.CreateUser(data)
	if err != nil {
		return err
	}
	// TODO
	sessionToken, _ := h.SessionService.Create(context.Background(), "Test")

	fmt.Println(sessionToken)
	return render(c, views.SignupForm(), 200)
}
