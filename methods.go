package gosocketio

import (
	"encoding/json"
	"reflect"
	"sync"

	"github.com/TymekV/gosf-socketio/protocol"
)

const (
	OnConnection    = "connection"
	OnDisconnection = "disconnection"
	OnError         = "error"
)

type systemHandler func(c *Channel)

type methods struct {
	messageHandlers sync.Map

	onConnection    systemHandler
	onDisconnection systemHandler
}

func (m *methods) initMethods() {
}

func (m *methods) On(method string, f interface{}) error {
	c, err := newCaller(f)
	if err != nil {
		return err
	}

	m.messageHandlers.Store(method, c)
	return nil
}

func (m *methods) findMethod(method string) (*caller, bool) {
	if f, ok := m.messageHandlers.Load(method); ok {
		return f.(*caller), true
	}

	return nil, false
}

func (m *methods) callLoopEvent(c *Channel, event string) {
	if m.onConnection != nil && event == OnConnection {
		m.onConnection(c)
	}
	if m.onDisconnection != nil && event == OnDisconnection {
		m.onDisconnection(c)
	}

	f, ok := m.findMethod(event)
	if !ok {
		return
	}

	f.callFunc(c, &struct{}{})
}

func (m *methods) processIncomingMessage(c *Channel, msg *protocol.Message) {
	switch msg.Type {
	case protocol.MessageTypeEmit:
		f, ok := m.findMethod(msg.Method)
		if !ok {
			return
		}

		if !f.ArgsPresent {
			f.callFunc(c, &struct{}{})
			return
		}

		data := f.getArgs()
		err := json.Unmarshal([]byte(msg.Args), &data)
		if err != nil {
			return
		}

		f.callFunc(c, data)

	case protocol.MessageTypeAckRequest:
		f, ok := m.findMethod(msg.Method)
		if !ok || !f.Out {
			return
		}

		var result []reflect.Value
		if f.ArgsPresent {
			data := f.getArgs()
			err := json.Unmarshal([]byte(msg.Args), &data)
			if err != nil {
				return
			}

			result = f.callFunc(c, data)
		} else {
			result = f.callFunc(c, &struct{}{})
		}

		ack := &protocol.Message{
			Type:  protocol.MessageTypeAckResponse,
			AckId: msg.AckId,
		}
		send(ack, c, result[0].Interface())

	case protocol.MessageTypeAckResponse:
		waiter, err := c.ack.getWaiter(msg.AckId)
		if err == nil {
			waiter <- msg.Args
		}
	}
}
