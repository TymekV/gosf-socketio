package gosocketio

import (
	"net"
	"net/http"
	"strconv"

	transport "github.com/TymekV/gosf-socketio/transport"
)

type Client struct {
	methods
	Channel
	opts *transport.Options
}

const (
	webSocketProtocol       = "ws://"
	webSocketSecureProtocol = "wss://"
	socketioUrl             = "/socket.io/?EIO=4&transport=websocket"
)

// GetUrl returns the ws/wss url by host and port
func GetUrl(host string, port int, secure bool) string {
	var prefix string
	if secure {
		prefix = webSocketSecureProtocol
	} else {
		prefix = webSocketProtocol
	}
	return prefix + net.JoinHostPort(host, strconv.Itoa(port)) + socketioUrl
}

// Dial connects to the server with the given options
func Dial(url string, tr transport.Transport, opts *transport.Options) (*Client, error) {
	c := &Client{opts: opts}
	c.initChannel()
	c.initMethods()

	requestHeaders := http.Header{}
	for _, extraHeader := range c.opts.ExtraHeaders {
		requestHeaders.Set(extraHeader.Key, extraHeader.Value)
	}

	conn, err := tr.Connect(url, opts)
	if err != nil {
		return nil, err
	}

	c.conn = conn

	go inLoop(&c.Channel, &c.methods)
	go outLoop(&c.Channel, &c.methods)
	go pinger(&c.Channel)

	return c, nil
}

// Close client connection
func (c *Client) Close() {
	closeChannel(&c.Channel, &c.methods)
}
