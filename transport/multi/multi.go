package multi

import (
	"errors"

	"github.com/micro/go-micro/transport"
)

// Option to be passed to NewTransport
type Option func(*multi)

type multi struct {
	transports []transport.Transport
	listen     transport.Transport
}

func (m *multi) Dial(addr string, opts ...transport.DialOption) (client transport.Client, err error) {
	err = errors.New("No transports provided")

	for _, transport := range m.transports {
		client, err = transport.Dial(addr, opts...)
		if client != nil && err == nil {
			return client, nil
		}
	}

	return client, err
}

func (m *multi) Listen(addr string, opts ...transport.ListenOption) (transport.Listener, error) {
	if m.listen != nil {
		return m.listen.Listen(addr, opts...)
	}
	return nil, errors.New("Not supported")
}

func (m *multi) String() string {
	return "multi"
}

// NewTransport will return a new multi Transport with the given options.
func NewTransport(opts ...Option) transport.Transport {
	m := &multi{}
	for _, o := range opts {
		o(m)
	}
	return m
}

// WithTransports allows you to specify which transports in the order
// of use that you wish to have them attempted.
func WithTransports(transport ...transport.Transport) Option {
	return func(m *multi) {
		m.transports = append(m.transports, transport...)
	}
}

// WithListenTransport allows you to specify an optional listening transport.
func WithListenTransport(transport transport.Transport) Option {
	return func(m *multi) {
		m.listen = transport
	}
}
