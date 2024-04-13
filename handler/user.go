package handler

import (
	"github.com/ThisJohan/go-htmx-chat/context"
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

func (h *UserHandler) Login(c echo.Context) error {
	return render(c, views.Login(), 200)
}

func (h *UserHandler) ProcessLogin(c echo.Context) error {
	var data struct {
		Email    string `form:"email"`
		Password string `form:"password"`
	}
	if err := c.Bind(&data); err != nil {
		return err
	}
	user, err := h.UserService.Authenticate(data.Email, data.Password)
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
	return render(c, views.LoginForm(), 200)
}

func (h *UserHandler) Me(c echo.Context) error {
	ac := c.(*context.AppContext)
	user := ac.User()

	return c.JSON(200, user)
}

func (h *UserHandler) AuthRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ac := c.(*context.AppContext)
		cookie, err := readCookie(ac, sessionTokenCookie)
		if err != nil {
			ac.Redirect(302, "/login")
			return err
		}
		userCache, err := h.SessionService.Get(ac.Request().Context(), cookie.Value)
		if err != nil {
			deleteCookie(ac, sessionTokenCookie)
			ac.Redirect(302, "/login")
			return err
		}
		ac.WithUser(userCache)
		return next(ac)
	}
}
