package mangos

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/rep"
	"github.com/go-mangos/mangos/protocol/req"
	"github.com/go-mangos/mangos/transport/tcp"
	"github.com/go-mangos/mangos/transport/tlstcp"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/transport"

	mls "github.com/micro/misc/lib/tls"
)

type ntport struct {
	opts transport.Options
}

type transportListener struct {
	listener mangos.Listener
	socket   mangos.Socket
}

type transportSocket struct {
	socket mangos.Socket
	msg    *mangos.Message
}

type transportClient struct {
	socket mangos.Socket
}

func listen(addr string, fn func(string) (mangos.Listener, error)) (mangos.Listener, error) {
	// host:port || host:min-max
	parts := strings.Split(addr, ":")

	//
	if len(parts) < 2 {
		return fn(addr)
	}

	// try to extract port range
	ports := strings.Split(parts[len(parts)-1], "-")

	// single port
	if len(ports) < 2 {
		return fn(addr)
	}

	// we have a port range

	// extract min port
	min, err := strconv.Atoi(ports[0])
	if err != nil {
		return nil, errors.New("unable to extract port range")
	}

	// extract max port
	max, err := strconv.Atoi(ports[1])
	if err != nil {
		return nil, errors.New("unable to extract port range")
	}

	// set host
	host := parts[:len(parts)-1]

	// range the ports
	for port := min; port <= max; port++ {
		// try bind to host:port
		ln, err := fn(fmt.Sprintf("%s:%d", host, port))
		if err == nil {
			return ln, nil
		}

		// hit max port
		if port == max {
			return nil, err
		}
	}

	// why are we here?
	return nil, fmt.Errorf("unable to bind to %s", addr)
}

func init() {
	cmd.DefaultTransports["mangos"] = NewTransport
}

func (t *transportSocket) Recv(m *transport.Message) error {
	if m == nil {
		return errors.New("message passed in is nil")
	}

	if err := json.Unmarshal(t.msg.Body, &m); err != nil {
		return err
	}

	return nil
}

func (t *transportSocket) Send(m *transport.Message) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	t.msg.Body = b

	return t.socket.SendMsg(t.msg)
}

func (t *transportSocket) Close() error {
	return nil
}

func (t *transportClient) Recv(m *transport.Message) error {
	if m == nil {
		return errors.New("message passed in is nil")
	}

	dat, err := t.socket.Recv()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(dat, &m); err != nil {
		return err
	}
	return nil
}

func (t *transportClient) Send(m *transport.Message) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return t.socket.Send(b)
}

func (t *transportClient) Close() error {
	return t.socket.Close()
}

func (t *transportListener) Addr() string {
	addr := strings.TrimPrefix(t.listener.Address(), "tls+tcp://")
	addr = strings.TrimPrefix(addr, "tcp://")
	return addr
}

func (t *transportListener) Close() error {
	return t.listener.Close()
}

func (t *transportListener) Accept(fn func(transport.Socket)) error {
	for {
		msg, err := t.socket.RecvMsg()
		if err != nil {
			return err
		}

		go fn(&transportSocket{
			socket: t.socket,
			msg:    msg,
		})
	}
}

func (n *ntport) Dial(addr string, dialOpts ...transport.DialOption) (transport.Client, error) {
	var sock mangos.Socket
	var d mangos.Dialer
	var err error
	options := transport.DialOptions{
		Timeout: transport.DefaultDialTimeout,
	}
	for _, o := range dialOpts {
		o(&options)
	}

	if sock, err = req.NewSocket(); err != nil {
		return nil, err
	}

	// TODO: support dial option here rather than using internal config
	if n.opts.Secure || n.opts.TLSConfig != nil {
		sock.AddTransport(tlstcp.NewTransport())
		config := n.opts.TLSConfig
		if config == nil {
			config = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		d, err = sock.NewDialer("tls+tcp://"+addr, map[string]interface{}{
			mangos.OptionTLSConfig: config,
		})
		if err != nil {
			return nil, err
		}
	} else {
		sock.AddTransport(tcp.NewTransport())
		d, err = sock.NewDialer("tcp://"+addr, nil)
	}

	if err != nil {
		return nil, err
	}

	sock.SetOption(mangos.OptionRecvDeadline, options.Timeout)
	sock.SetOption(mangos.OptionSendDeadline, options.Timeout)
	d.Dial()

	return &transportClient{
		socket: sock,
	}, err
}

func (n *ntport) Listen(addr string, opts ...transport.ListenOption) (transport.Listener, error) {
	var options transport.ListenOptions
	for _, o := range opts {
		o(&options)
	}

	var sock mangos.Socket
	var l mangos.Listener
	var err error

	if sock, err = rep.NewSocket(); err != nil {
		return nil, err
	}
	if err = sock.SetOption(mangos.OptionRaw, true); err != nil {
		return nil, err
	}

	// TODO: support use of listen options
	if n.opts.Secure || n.opts.TLSConfig != nil {
		sock.AddTransport(tlstcp.NewTransport())
		config := n.opts.TLSConfig

		fn := func(addr string) (mangos.Listener, error) {
			if config == nil {
				hosts := []string{addr}

				// check if its a valid host:port
				if host, _, err := net.SplitHostPort(addr); err == nil {
					if len(host) == 0 {
						hosts = getIPAddrs()
					} else {
						hosts = []string{host}
					}
				}

				// generate a certificate
				cert, err := mls.Certificate(hosts...)
				if err != nil {
					return nil, err
				}
				config = &tls.Config{Certificates: []tls.Certificate{cert}}
			}
			return sock.NewListener("tls+tcp://"+addr, map[string]interface{}{
				mangos.OptionTLSConfig: config,
			})
		}

		l, err = listen(addr, fn)
	} else {
		sock.AddTransport(tcp.NewTransport())
		fn := func(addr string) (mangos.Listener, error) {
			return sock.NewListener("tcp://"+addr, nil)
		}

		l, err = listen(addr, fn)
	}

	if err != nil {
		return nil, err
	}

	l.Listen()
	return &transportListener{
		listener: l,
		socket:   sock,
	}, err
}

func (n *ntport) String() string {
	return "mangos"
}

func getIPAddrs() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	var ipAddrs []string

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue
			}

			ipAddrs = append(ipAddrs, ip.String())
		}
	}

	return ipAddrs
}

func NewTransport(opts ...transport.Option) transport.Transport {
	var options transport.Options
	for _, o := range opts {
		o(&options)
	}

	return &ntport{opts: options}
}
