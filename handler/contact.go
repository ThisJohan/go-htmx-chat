package handler

import (
	"github.com/ThisJohan/go-htmx-chat/context"
	"github.com/ThisJohan/go-htmx-chat/models"
	"github.com/labstack/echo/v4"
)

type ContactHandler struct {
	ContactService *models.ContactService
}

func (h *ContactHandler) GetContacts(c echo.Context) error {
	ac := c.(*context.AppContext)
	user := ac.User()

	contacts, err := h.ContactService.GetContacts(user.ID)
	if err != nil {
		return err
	}

	return ac.JSON(200, contacts)
}
