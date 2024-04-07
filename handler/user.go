package handler

import (
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
	var data models.User
	if err := c.Bind(&data); err != nil {
		return err
	}
	user, err := h.UserService.CreateUser(data)
	if err != nil {
		return err
	}

	sessionToken, err := h.SessionService.Create(c.Request().Context(), models.UserCache{
		ID: user.ID, Email: user.Email, FirstName: user.FirstName, LastName: user.LastName,
	})
	if err != nil {
		return err
	}
	writeCookie(c, sessionTokenCookie, sessionToken)

	return render(c, views.SignupForm(), 200)
}

func (h *UserHandler) Me(c echo.Context) error {
	cookie, err := readCookie(c, sessionTokenCookie)
	if err != nil {
		return err
	}
	userCache, err := h.SessionService.Get(c.Request().Context(), cookie.Value)
	if err != nil {
		deleteCookie(c, sessionTokenCookie)
		return err
	}

	return c.JSON(200, userCache)
}
