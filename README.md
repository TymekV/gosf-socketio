# GoSocketIO

[![GoDoc](https://godoc.org/github.com/TymekV/gosf-socketio?status.svg)](https://godoc.org/github.com/TymekV/gosf-socketio)
[![Go Report Card](https://goreportcard.com/badge/github.com/TymekV/gosf-socketio)](https://goreportcard.com/report/github.com/TymekV/gosf-socketio)

GoSocketIO is a [Socket.IO](http://socket.io) client and server library for Go, compatible with Socket.IO v3 and v4. This package allows you to build real-time, bidirectional communication applications with ease.

This library is based off the [GoLang SocketIO Framework](https://github.com/ambelovsky/gosf)

This library was built with previous contributions by:
- [ambelovsky](https://github.com/ambelovsky)
- [joaopandolfi](https://github.com/joaopandolfi)


## Features

- **Socket.IO v3 and v4 Compatible**: Fully compatible with the latest versions of Socket.IO.
- **Extra Headers Support**: Customize HTTP headers for WebSocket connections.
- **Event-based Architecture**: Simplified event handling for connection, disconnection, and custom events.
- **Room Support**: Easily broadcast messages to specific groups of clients.
- **Ping/Pong Mechanism**: Built-in heartbeat mechanism to keep connections alive.

## Installation

To install GoSocketIO, use `go get`:

```bash
go get github.com/TymekV/gosf-socketio
```

## Usage

### Server Example

Here’s a simple WebSocket server using GoSocketIO:

```go
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/TymekV/gosf-socketio"
    "github.com/TymekV/gosf-socketio/transport"
)

func main() {
    server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

    server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
        fmt.Println("New client connected:", c.Id())
        c.Join("chat")
    })

    server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
        fmt.Println("Client disconnected:", c.Id())
    })

    server.On("chat message", func(c *gosocketio.Channel, message string) {
        fmt.Printf("Received message from %s: %s\n", c.Id(), message)
        server.BroadcastTo("chat", "chat message", message)
    })

    http.Handle("/socket.io/", server)

    addr := "localhost:8080"
    fmt.Printf("Serving at %s...\n", addr)
    log.Fatal(http.ListenAndServe(addr, nil))
}
```
For more details, see the [Server Example](https://github.com/TymekV/gosf-socketio/examples/example-server)

### Go Client Example

Here’s how you can create a Go client that connects to the server:

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/TymekV/gosf-socketio"
    "github.com/TymekV/gosf-socketio/transport"
)

func main() {
    wsTransport := transport.GetDefaultWebsocketTransport()
    options := &gosocketio.Options{
        ExtraHeaders: []gosocketio.ExtraHeaders{
            {Key: "Content-Type", Value: "application/json"},
            {Key: "Authorization", Value: "Bearer YOUR_ACCESS_TOKEN"},
        },
        PingTimeout:  60 * time.Second,
        PingInterval: 30 * time.Second,
    }

    client, err := gosocketio.Dial("ws://localhost:8080/socket.io/?EIO=3&transport=websocket", wsTransport, options)
    if err != nil {
        log.Fatalf("Failed to connect to server: %v", err)
    }
    defer client.Close()

    client.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
        fmt.Println("Connected to server")
    })

    client.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
        fmt.Println("Disconnected from server")
    })

    client.On("chat message", func(c *gosocketio.Channel, message string) {
        fmt.Printf("Received message: %s\n", message)
    })

    err = client.Emit("chat message", "Hello, Server!")
    if err != nil {
        log.Printf("Failed to send message: %v", err)
    }

    select {}
}

```

For more details, see the [Go Client Example](https://github.com/TymekV/gosf-socketio/examples/example-go-client)

### JavaScript Client Example

Here’s a simple JavaScript client using Socket.IO:

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Socket.IO Client</title>
</head>
<body>
    <h1>WebSocket Chat</h1>
    <div id="chat">
        <ul id="messages"></ul>
    </div>
    <form id="form" action="">
        <input id="input" autocomplete="off" /><button>Send</button>
    </form>
    <script src="https://cdn.socket.io/4.0.1/socket.io.min.js"></script>
    <script>
        var socket = io('http://localhost:8080');

        socket.on('connect', function() {
            console.log('Connected to server');
        });

        socket.on('disconnect', function() {
            console.log('Disconnected from server');
        });

        socket.on('chat message', function(msg) {
            var item = document.createElement('li');
            item.textContent = msg;
            document.getElementById('messages').appendChild(item);
        });

        document.getElementById('form').addEventListener('submit', function(e) {
            e.preventDefault();
            var input = document.getElementById('input');
            socket.emit('chat message', input.value);
            input.value = '';
        });
    </script>
</body>
</html>
```
For more details, see the [JavaScript Client Example](https://github.com/TymekV/gosf-socketio/examples/example-js-client)


## Contributing

We welcome contributions! Please check out the [Contributing Guide](https://github.com/TymekV/gosf-socketio/CONTRIBUTING.md) for more details.

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/TymekV/gosf-socketio/LICENSE) file for details.

