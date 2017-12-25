package nats

import (
	"time"

	"github.com/micro/go-micro/broker"
	"github.com/nats-io/nats"
)

// default NATS client values
var (
	DefaultMaxReconnect     = 60
	DefaultReconnectWait    = 2 * time.Second
	DefaultTimeout          = 2 * time.Second
	DefaultPingInterval     = 2 * time.Minute
	DefaultMaxPingOut       = 2
	DefaultMaxChanLen       = 8192            // 8k
	DefaultReconnectBufSize = 8 * 1024 * 1024 // 8MB

	optionsKey = optionsKeyType{}
)

// brokerOptions contains NATS specific options
type brokerOptions struct {
	maxReconnect             int
	name                     string
	reconnectWait            time.Duration
	timeout                  time.Duration
	pingInterval             time.Duration
	maxPingOut               int
	maxChanLen               int
	reconnectBufSize         int
	closedHandler            func(*nats.Conn)
	disconnectHandler        func(*nats.Conn)
	discoveredServersHandler func(*nats.Conn)
	reconnectHandler         func(*nats.Conn)
	errorHandler             func(*nats.Conn, *nats.Subscription, error)
}

type optionsKeyType struct{}

func MaxReconnect(n int) broker.Option {
	return func(o *broker.Options) {
		bo := o.Context.Value(optionsKey).(*brokerOptions)
		bo.maxReconnect = n
	}
}

func ReconnectWait(d time.Duration) broker.Option {
	return func(o *broker.Options) {
		bo := o.Context.Value(optionsKey).(*brokerOptions)
		bo.reconnectWait = d
	}
}

func Timeout(d time.Duration) broker.Option {
	return func(o *broker.Options) {
		bo := o.Context.Value(optionsKey).(*brokerOptions)
		bo.timeout = d
	}
}

func PingInterval(d time.Duration) broker.Option {
	return func(o *broker.Options) {
		bo := o.Context.Value(optionsKey).(*brokerOptions)
		bo.pingInterval = d
	}
}

func MaxPingOut(n int) broker.Option {
	return func(o *broker.Options) {
		bo := o.Context.Value(optionsKey).(*brokerOptions)
		bo.maxPingOut = n
	}
}

func MaxChanLen(n int) broker.Option {
	return func(o *broker.Options) {
		bo := o.Context.Value(optionsKey).(*brokerOptions)
		bo.maxChanLen = n
	}
}

func ReconnectBufSize(n int) broker.Option {
	return func(o *broker.Options) {
		bo := o.Context.Value(optionsKey).(*brokerOptions)
		bo.reconnectBufSize = n
	}
}

func ClosedHandler(cb nats.ConnHandler) broker.Option {
	return func(o *broker.Options) {
		bo := o.Context.Value(optionsKey).(*brokerOptions)
		bo.closedHandler = cb
	}
}

func DisconnectHandler(cb nats.ConnHandler) broker.Option {
	return func(o *broker.Options) {
		bo := o.Context.Value(optionsKey).(*brokerOptions)
		bo.disconnectHandler = cb
	}
}

func DiscoveredServersHandler(cb nats.ConnHandler) broker.Option {
	return func(o *broker.Options) {
		bo := o.Context.Value(optionsKey).(*brokerOptions)
		bo.discoveredServersHandler = cb
	}
}

func ReconnectHandler(cb nats.ConnHandler) broker.Option {
	return func(o *broker.Options) {
		bo := o.Context.Value(optionsKey).(*brokerOptions)
		bo.reconnectHandler = cb
	}
}

func ErrorHandler(cb nats.ErrHandler) broker.Option {
	return func(o *broker.Options) {
		bo := o.Context.Value(optionsKey).(*brokerOptions)
		bo.errorHandler = cb
	}
}

func Name(s string) broker.Option {
	return func(o *broker.Options) {
		bo := o.Context.Value(optionsKey).(*brokerOptions)
		bo.name = s
	}
}
