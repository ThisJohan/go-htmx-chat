package handler

import (
	"fmt"
	"strconv"

	"github.com/ThisJohan/go-htmx-chat/context"
	"github.com/ThisJohan/go-htmx-chat/models"
	view "github.com/ThisJohan/go-htmx-chat/views/chat"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{}
)

type SocketHandler struct {
	ChatService    *models.ChatService
	ContactService *models.ContactService
}

func (h *SocketHandler) Demo(c echo.Context) error {
	ac := c.(*context.AppContext)
	user := ac.User()
	contacts, err := h.ContactService.GetContacts(user.ID)
	if err != nil {
		return err
	}
	return render(c, view.Demo(contacts), 200)
}

func (h *SocketHandler) SelectContact(c echo.Context) error {
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

func (h *SocketHandler) Hello(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		err = ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			return err
		}
		_, message, err := ws.ReadMessage()
		if err != nil {
			return err
		}

		fmt.Println(string(message))
	}
}

func (h *SocketHandler) Chat(c echo.Context) error {
	ac := c.(*context.AppContext)
	user := ac.User()

	ws, err := upgrader.Upgrade(ac.Response(), ac.Request(), nil)
	if err != nil {
		return err
	}
	// defer ws.Close()

	client := h.ChatService.Hub.Register(ws, user.ID)

	go client.ReadPump()
	go client.WritePump()

	return nil
}
