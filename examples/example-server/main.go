package main

import (
	"fmt"
	"log"
	"net/http"

	gosocketio "github.com/TymekV/gosf-socketio"
	"github.com/TymekV/gosf-socketio/transport"
)

func main() {
	// Create a new WebSocket server with the default WebSocket transport
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	// Define event handlers
	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		fmt.Println("New client connected:", c.Id())
		c.Join("chat")
	})

	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		fmt.Println("Client disconnected:", c.Id())
	})

	server.On("chat message", func(c *gosocketio.Channel, message string) {
		fmt.Printf("Received message from %s: %s\n", c.Id(), message)
		// Broadcast the message to all clients in the "chat" room
		server.BroadcastTo("chat", "chat message", message)
	})

	// Set up the HTTP server to handle WebSocket connections
	http.Handle("/socket.io/", server)

	// Start the HTTP server
	addr := "localhost:8080"
	fmt.Printf("Serving at %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
