package models

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type ChatService struct {
	Hub   *UsersHub
	DB    *sqlx.DB
	Redis *redis.Client
}

type Message struct {
	ToUser   int
	FromUser int
	Content  string
}

type UsersHub struct {
	// map of userIds point to clients
	clients    map[int]*HubClient
	register   chan *HubClient
	unregister chan *HubClient
	deliver    chan *Message
}

type HubClient struct {
	hub  *UsersHub
	conn *websocket.Conn
	// For sending messages from server to ws client
	send            chan []byte
	userId          int
	activeContactId int
}

// Send new messages to client received from chan hub.deliver
func (client *HubClient) WritePump() {
	for {
		select {
		case message, ok := <-client.send:
			if ok {
				client.conn.WriteMessage(websocket.TextMessage, message)
			}
		}
	}
}

// Read incoming message from user and handle it to hub
func (c *HubClient) ReadPump() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		if string(message) == "" {
			continue
		}
		var data map[string]interface{}
		err = json.Unmarshal(message, &data)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if data["type"] == "select_contact" {
			// TODO: fix this with real data
			// think about how to have userId and not contactId
			c.activeContactId = 2
		}
	}
}

func (hub *UsersHub) Deliver(message *Message) {
	hub.deliver <- message
}

func (hub *UsersHub) Register(conn *websocket.Conn, userId int) *HubClient {
	client := &HubClient{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte),
		userId: userId,
	}
	hub.register <- client
	return client
}

func (hub *UsersHub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.clients[client.userId] = client
		case client := <-hub.unregister:
			if _, ok := hub.clients[client.userId]; ok {
				delete(hub.clients, client.userId)
				close(client.send)
			}
		case message := <-hub.deliver:
			if client, ok := hub.clients[message.ToUser]; ok {
				// TODO: client must know message is sent from who ;)
				// TODO: save message to redis and make update to database
				if client.activeContactId == 2 {
					client.send <- []byte(message.Content)
				}
			}
		}
	}
}

func NewChatService(DB *sqlx.DB, Redis *redis.Client) *ChatService {
	return &ChatService{
		Hub: &UsersHub{
			clients:    make(map[int]*HubClient),
			register:   make(chan *HubClient),
			unregister: make(chan *HubClient),
			deliver:    make(chan *Message),
		},
		DB:    DB,
		Redis: Redis,
	}
}

func (s *ChatService) HandleNewMessage(message *Message) error {
	return nil
}
