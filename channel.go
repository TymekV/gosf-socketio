package gosocketio

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/erock530/gosf-socketio/protocol"
	"github.com/erock530/gosf-socketio/transport"
)

const (
	queueBufferSize = 10000
)

type Channel struct {
	conn transport.Connection

	out    chan string
	header Header

	alive     bool
	aliveLock sync.Mutex

	ack ackProcessor

	server  *Server
	ip      string
	request *http.Request
}

// Header represents engine.io header to send or receive
type Header struct {
	Sid          string   `json:"sid"`
	Upgrades     []string `json:"upgrades"`
	PingInterval int      `json:"pingInterval"`
	PingTimeout  int      `json:"pingTimeout"`
}

var ErrorWrongHeader = errors.New("Wrong header")

// Initialize channel, map, and set active
func (c *Channel) initChannel() {
	c.out = make(chan string, queueBufferSize)
	c.setAliveValue(true)
}

// Get the ID of the current socket connection
func (c *Channel) Id() string {
	return c.header.Sid
}

// Check if the Channel is still alive
func (c *Channel) IsAlive() bool {
	c.aliveLock.Lock()
	isAlive := c.alive
	c.aliveLock.Unlock()

	return isAlive
}

func (c *Channel) setAliveValue(value bool) {
	c.aliveLock.Lock()
	c.alive = value
	c.aliveLock.Unlock()
}

// Close the channel
func closeChannel(c *Channel, m *methods, args ...interface{}) error {
	if !c.IsAlive() {
		return nil
	}

	c.conn.Close()
	c.setAliveValue(false)

	for len(c.out) > 0 {
		<-c.out
	}

	c.out <- protocol.CloseMessage
	m.callLoopEvent(c, OnDisconnection)

	deleteOverflooded(c)

	return nil
}

// Incoming messages loop
func inLoop(c *Channel, m *methods) error {
	for {
		pkg, err := c.conn.GetMessage()
		if err != nil {
			return closeChannel(c, m, err)
		}
		msg, err := protocol.Decode(pkg)
		if err != nil {
			closeChannel(c, m, protocol.ErrorWrongPacket)
			return err
		}

		switch msg.Type {
		case protocol.MessageTypeOpen:
			if err := json.Unmarshal([]byte(msg.Source[1:]), &c.header); err != nil {
				closeChannel(c, m, ErrorWrongHeader)
			}
			m.callLoopEvent(c, OnConnection)
		case protocol.MessageTypePing:
			c.out <- protocol.PongMessage
		case protocol.MessageTypePong:
		default:
			go m.processIncomingMessage(c, msg)
		}
	}
}

var overflooded sync.Map

func deleteOverflooded(c *Channel) {
	overflooded.Delete(c)
}

func storeOverflow(c *Channel) {
	overflooded.Store(c, struct{}{})
}

// Outgoing messages loop
func outLoop(c *Channel, m *methods) error {
	for {
		outBufferLen := len(c.out)
		if outBufferLen >= queueBufferSize-1 {
			return closeChannel(c, m, ErrorSocketOverflood)
		} else if outBufferLen > int(queueBufferSize/2) {
			storeOverflow(c)
		} else {
			deleteOverflooded(c)
		}

		msg := <-c.out
		if msg == protocol.CloseMessage {
			return nil
		}

		err := c.conn.WriteMessage(msg)
		if err != nil {
			return closeChannel(c, m, err)
		}
	}
}

// Pinger sends ping messages to keep connection alive
func pinger(c *Channel) {
	interval, _ := c.conn.PingParams()
	ticker := time.NewTicker(interval)
	for {
		<-ticker.C
		if !c.IsAlive() {
			return
		}
		c.out <- protocol.PingMessage
	}
}
