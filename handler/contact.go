package handler

import (
	"net/http"
	"strconv"

	"github.com/ThisJohan/go-htmx-chat/context"
	"github.com/ThisJohan/go-htmx-chat/models"
	view "github.com/ThisJohan/go-htmx-chat/views/chat"
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

func (h *ContactHandler) SelectContact(c echo.Context) error {
	ac := c.(*context.AppContext)
	user := ac.User()
	contactId, _ := strconv.Atoi(ac.Param("id"))
	contact, err := h.ContactService.GetContactByIdAndValidate(contactId, user.ID)
	if err != nil {
		return err
	}

	render(ac, view.Chat(contact), 200)
	return render(ac, view.Contact(contact, false), 200)
}

func (h *ContactHandler) CreateContact(c echo.Context) error {
	ac := c.(*context.AppContext)
	user := ac.User()
	email := ac.FormValue("email")
	contact, err := h.ContactService.CreateContact(email, user.ID)
	if err != nil {
		return err
	}
	return render(ac, view.OOBContact(contact), http.StatusCreated)
}
