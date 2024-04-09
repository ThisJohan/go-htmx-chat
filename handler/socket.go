package handler

import (
	"fmt"

	view "github.com/ThisJohan/go-htmx-chat/views/chat"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{}
)

type SocketHandler struct{}

func (h *SocketHandler) Demo(c echo.Context) error {
	return render(c, view.Demo(), 200)
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
