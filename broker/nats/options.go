package nats

import (
	"time"

	"github.com/micro/go-micro/broker"
	"github.com/nats-io/nats"
)

// default NATS client values
var (
	DefaultNatsMaxReconnect     = 60
	DefaultNatsReconnectWait    = 2 * time.Second
	DefaultNatsTimeout          = 2 * time.Second
	DefaultNatsPingInterval     = 2 * time.Minute
	DefaultNatsMaxPingOut       = 2
	DefaultNatsMaxChanLen       = 8192            // 8k
	DefaultNatsReconnectBufSize = 8 * 1024 * 1024 // 8MB
	DefaultNatsAllowReconnect   = true

	optionsKey = optionsKeyType{}
)

// natsOptions contains NATS specific options
type natsOptions struct {
	maxReconnect             int
	name                     string
	reconnectWait            time.Duration
	timeout                  time.Duration
	pingInterval             time.Duration
	maxPingOut               int
	maxChanLen               int
	reconnectBufSize         int
	allowReconnect           bool
	username                 string
	password                 string
	token                    string
	closedHandler            func(*nats.Conn)
	disconnectHandler        func(*nats.Conn)
	discoveredServersHandler func(*nats.Conn)
	reconnectHandler         func(*nats.Conn)
	errorHandler             func(*nats.Conn, *nats.Subscription, error)
}

type optionsKeyType struct{}

func MaxReconnect(n int) broker.Option {
	return func(o *broker.Options) {
		no := o.Context.Value(optionsKey).(*natsOptions)
		no.maxReconnect = n
	}
}

func ReconnectWait(d time.Duration) broker.Option {
	return func(o *broker.Options) {
		no := o.Context.Value(optionsKey).(*natsOptions)
		no.reconnectWait = d
	}
}

func Timeout(d time.Duration) broker.Option {
	return func(o *broker.Options) {
		no := o.Context.Value(optionsKey).(*natsOptions)
		no.timeout = d
	}
}

func PingInterval(d time.Duration) broker.Option {
	return func(o *broker.Options) {
		no := o.Context.Value(optionsKey).(*natsOptions)
		no.pingInterval = d
	}
}

func MaxPingOut(n int) broker.Option {
	return func(o *broker.Options) {
		no := o.Context.Value(optionsKey).(*natsOptions)
		no.maxPingOut = n
	}
}

func MaxChanLen(n int) broker.Option {
	return func(o *broker.Options) {
		no := o.Context.Value(optionsKey).(*natsOptions)
		no.maxChanLen = n
	}
}

func ReconnectBufSize(n int) broker.Option {
	return func(o *broker.Options) {
		no := o.Context.Value(optionsKey).(*natsOptions)
		no.reconnectBufSize = n
	}
}

func AllowReconnect(b bool) broker.Option {
	return func(o *broker.Options) {
		no := o.Context.Value(optionsKey).(*natsOptions)
		no.allowReconnect = b
	}
}

func Username(s string) broker.Option {
	return func(o *broker.Options) {
		no := o.Context.Value(optionsKey).(*natsOptions)
		no.username = s
	}
}

func Password(s string) broker.Option {
	return func(o *broker.Options) {
		no := o.Context.Value(optionsKey).(*natsOptions)
		no.password = s
	}
}

func Token(s string) broker.Option {
	return func(o *broker.Options) {
		no := o.Context.Value(optionsKey).(*natsOptions)
		no.token = s
	}
}

func ClosedHandler(cb nats.ConnHandler) broker.Option {
	return func(o *broker.Options) {
		no := o.Context.Value(optionsKey).(*natsOptions)
		no.closedHandler = cb
	}
}

func DisconnectHandler(cb nats.ConnHandler) broker.Option {
	return func(o *broker.Options) {
		no := o.Context.Value(optionsKey).(*natsOptions)
		no.disconnectHandler = cb
	}
}

func DiscoveredServersHandler(cb nats.ConnHandler) broker.Option {
	return func(o *broker.Options) {
		no := o.Context.Value(optionsKey).(*natsOptions)
		no.discoveredServersHandler = cb
	}
}

func ReconnectHandler(cb nats.ConnHandler) broker.Option {
	return func(o *broker.Options) {
		no := o.Context.Value(optionsKey).(*natsOptions)
		no.reconnectHandler = cb
	}
}

func ErrorHandler(cb nats.ErrHandler) broker.Option {
	return func(o *broker.Options) {
		no := o.Context.Value(optionsKey).(*natsOptions)
		no.errorHandler = cb
	}
}

func Name(s string) broker.Option {
	return func(o *broker.Options) {
		no := o.Context.Value(optionsKey).(*natsOptions)
		no.name = s
	}
}
