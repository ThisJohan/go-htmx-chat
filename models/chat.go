package models

import "github.com/gorilla/websocket"

type ChatService struct {
	Hub *UsersHub
}

type UsersHub struct {
	clients    map[*HubClient]bool
	register   chan *HubClient
	unregister chan *HubClient
}

type HubClient struct {
	hub  *UsersHub
	conn *websocket.Conn
	// For sending messages from server to ws client
	send chan []byte
}

func (hub *UsersHub) Register(conn *websocket.Conn) *HubClient {
	client := &HubClient{
		hub:  hub,
		conn: conn,
		send: make(chan []byte),
	}
	hub.register <- client
	return client
}

func NewChatService() *ChatService {
	return &ChatService{
		Hub: &UsersHub{
			clients:    make(map[*HubClient]bool),
			register:   make(chan *HubClient),
			unregister: make(chan *HubClient),
		},
	}
}
