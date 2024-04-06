package handler

import (
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
	user, err := h.UserService.CreateUser(data)
	if err != nil {
		return err
	}
	redisData := map[string]interface{}{
		"id":         user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	}

	sessionToken, err := h.SessionService.Create(c.Request().Context(), redisData)
	if err != nil {
		return err
	}

	fmt.Println(sessionToken)
	return render(c, views.SignupForm(), 200)
}
