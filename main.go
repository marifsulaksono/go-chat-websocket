package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

/*
WebSocketManager handles WebSocket connections and message broadcasting

Clients: a map of user ID to WebSocket connections
Broadcast: a channel to broadcast messages
Register: a channel to register new clients
Unregister: a channel to unregister clients
*/
type WebSocketManager struct {
	Clients    map[string]*websocket.Conn
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
}

/*
Client represents a single WebSocket connection
*/
type Client struct {
	ID     string
	Socket *websocket.Conn
	Send   chan Message
}

type Message struct {
	Sender   string    `json:"sender"`
	Receiver string    `json:"receiver"`
	Content  string    `json:"content"`
	Time     time.Time `json:"time"`
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		Clients:    make(map[string]*websocket.Conn),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (manager *WebSocketManager) Start() {
	for {
		select {
		case client := <-manager.Register:
			manager.Clients[client.ID] = client.Socket
			log.Printf("Client %s connected", client.ID)

		case client := <-manager.Unregister:
			if _, ok := manager.Clients[client.ID]; ok {
				close(client.Send)
				delete(manager.Clients, client.ID)
				log.Printf("Client %s successfully disconnected", client.ID)
			}

		case message := <-manager.Broadcast:
			// send messages to the recipient (personal chat for now)
			for clientID, conn := range manager.Clients {
				if strings.EqualFold(clientID, message.Receiver) {
					err := conn.WriteJSON(message)
					if err != nil {
						log.Printf("Error sending message to client %s: %v", clientID, err)
						conn.Close()
						delete(manager.Clients, clientID)
					} else {
						log.Printf("Message sent to %s", clientID)
					}
				}
			}
		}
	}
}

/*
WebSocketHandler handles WebSocket connections and message broadcasts

This function sets up a WebSocket connection, registers the client with the manager, and starts a goroutine to handle incoming messages from the client.
*/
func WebSocketHandler(manager *WebSocketManager, c echo.Context) error {
	// Set up WebSocket connection
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true }, // allow all origins for simplicity
	}).Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return fmt.Errorf("could not upgrade to WebSocket: %v", err)
	}

	// Get user ID form query parameter for the client identifier
	userID := c.QueryParam("user_id")
	if userID == "" {
		return fmt.Errorf("user_id is required")
	}

	client := &Client{
		ID:     userID,
		Socket: conn,
		Send:   make(chan Message),
	}

	// Register the client with the manager
	manager.Register <- client

	// Handle incoming messages from the client
	go client.ListenForMessages(manager)

	// Keep connection open
	select {}
}

/*
ListenForMessages listens for incoming messages from the client and broadcasts them

This function reads messages from the client's WebSocket connection and sends them to the manager for broadcasting
*/
func (client *Client) ListenForMessages(manager *WebSocketManager) {
	defer func() {
		manager.Unregister <- client
		client.Socket.Close()
	}()

	// Keep receiving messages from the WebSocket
	for {
		var msg Message
		err := client.Socket.ReadJSON(&msg)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Printf("Client %s disconnected gracefully", client.ID)
			} else {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		msg.Time = time.Now().Truncate(time.Second) // Truncate(time.Second) removes any sub-second precision

		log.Printf("Received message from %s to %s at %s: %s", msg.Sender, msg.Receiver, msg.Time, msg.Content)
		// Send message to broadcast channel
		manager.Broadcast <- msg
	}
}

func main() {
	manager := NewWebSocketManager()

	go manager.Start() // Start WebSocketManager in a separate goroutine

	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/ws/chat", func(c echo.Context) error {
		return WebSocketHandler(manager, c)
	})

	log.Fatal(e.Start(":8080"))
}