package transport

import (
	"time"
)

type ExtraHeaders struct {
	Key   string
	Value string
}

type Options struct {
	PingTimeout  time.Duration
	PingInterval time.Duration

	Transports   []Transport
	ExtraHeaders []ExtraHeaders
}
