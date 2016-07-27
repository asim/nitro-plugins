package nats

import (
	"strings"

	"github.com/micro/go-micro/cmd"
	microtransport "github.com/micro/go-micro/transport"
	"github.com/nats-io/nats"
)

type transport struct {
	addrs []string
	opts  microtransport.Options
}

func init() {
	cmd.DefaultTransports["nats"] = NewTransport
}

func (n *transport) Dial(addr string, dialOpts ...microtransport.DialOption) (microtransport.Client, error) {
	dopts := microtransport.DialOptions{
		Timeout: microtransport.DefaultDialTimeout,
	}

	for _, o := range dialOpts {
		o(&dopts)
	}

	opts := nats.DefaultOptions
	opts.Servers = n.addrs
	opts.Secure = n.opts.Secure
	opts.TLSConfig = n.opts.TLSConfig
	opts.Timeout = dopts.Timeout

	// secure might not be set
	if n.opts.TLSConfig != nil {
		opts.Secure = true
	}

	c, err := opts.Connect()
	if err != nil {
		return nil, err
	}

	id := nats.NewInbox()
	sub, err := c.SubscribeSync(id)
	if err != nil {
		return nil, err
	}

	return &client{
		conn: c,
		addr: addr,
		id:   id,
		sub:  sub,
	}, nil
}

func (n *transport) Listen(addr string, listenOpts ...microtransport.ListenOption) (microtransport.Listener, error) {
	opts := nats.DefaultOptions
	opts.Servers = n.addrs
	opts.Secure = n.opts.Secure
	opts.TLSConfig = n.opts.TLSConfig

	// secure might not be set
	if n.opts.TLSConfig != nil {
		opts.Secure = true
	}

	c, err := opts.Connect()
	if err != nil {
		return nil, err
	}

	return &listener{
		addr: nats.NewInbox(),
		conn: c,
		exit: make(chan bool, 1),
		so:   make(map[string]*socket),
	}, nil
}

func (n *transport) String() string {
	return "nats"
}

func NewTransport(opts ...microtransport.Option) microtransport.Transport {
	var options microtransport.Options
	for _, o := range opts {
		o(&options)
	}

	var cAddrs []string

	for _, addr := range options.Addrs {
		if len(addr) == 0 {
			continue
		}
		if !strings.HasPrefix(addr, "nats://") {
			addr = "nats://" + addr
		}
		cAddrs = append(cAddrs, addr)
	}

	if len(cAddrs) == 0 {
		cAddrs = []string{nats.DefaultURL}
	}

	return &transport{
		addrs: cAddrs,
		opts:  options,
	}
}
