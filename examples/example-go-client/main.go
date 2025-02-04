package main

import (
	"fmt"
	"log"
	"time"

	gosocketio "github.com/TymekV/gosf-socketio"
	"github.com/TymekV/gosf-socketio/transport"
)

func main() {
	// Define the WebSocket transport with default settings
	wsTransport := transport.GetDefaultWebsocketTransport()

	// Set the client options, including any custom headers if needed
	options := &transport.Options{
		ExtraHeaders: []transport.ExtraHeaders{
			{Key: "Content-Type", Value: "application/json"},
			{Key: "Authorization", Value: "Bearer YOUR_ACCESS_TOKEN"},
		},
		PingTimeout:  60 * time.Second,
		PingInterval: 30 * time.Second,
	}

	// Connect to the WebSocket server
	client, err := gosocketio.Dial("ws://localhost:8080/socket.io/?EIO=3&transport=websocket", wsTransport, options)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer client.Close()

	// Define event handlers
	client.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		fmt.Println("Connected to server")
	})

	client.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		fmt.Println("Disconnected from server")
	})

	client.On("chat message", func(c *gosocketio.Channel, message string) {
		fmt.Printf("Received message: %s\n", message)
	})

	// Emit a message to the server
	err = client.Emit("chat message", "Hello, Server!")
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}

	// Keep the application running to listen for messages
	select {}
}
