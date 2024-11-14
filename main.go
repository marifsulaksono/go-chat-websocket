package main

import (
	"log"
	"marfs-websocket/chat"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	manager := chat.NewWebSocketManager()
	go manager.Start() // Start WebSocketManager in a separate goroutine

	e := echo.New()
	e.Use(middleware.Logger())

	// WebSocket endpoint
	e.GET("/ws/chat", func(c echo.Context) error {
		return chat.WebSocketHandler(manager, c)
	})

	log.Fatal(e.Start(":8080"))
}
