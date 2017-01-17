// Package nats provides a NATS registry using broadcast queries
package nats

import (
	"crypto/tls"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/registry"
	"github.com/nats-io/nats"
)

type natsRegistry struct {
	addrs []string
	opts  registry.Options

	sync.RWMutex
	conn      *nats.Conn
	services  map[string][]*registry.Service
	listeners map[string]chan bool
}

func init() {
	cmd.DefaultRegistries["nats"] = NewRegistry
}

var (
	QueryTopic = "micro.registry.nats.query"
	WatchTopic = "micro.registry.nats.watch"

	DefaultTimeout = time.Millisecond * 100
)

func newConn(addrs []string, secure bool, config *tls.Config) (*nats.Conn, error) {
	var cAddrs []string
	for _, addr := range addrs {
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

	opts := nats.DefaultOptions
	opts.Servers = cAddrs
	opts.Secure = secure
	opts.TLSConfig = config

	// secure might not be set
	if config != nil {
		opts.Secure = true
	}

	c, err := opts.Connect()
	if err != nil {
		return nil, err
	}
	return c, err
}

func (n *natsRegistry) getConn() (*nats.Conn, error) {
	n.Lock()
	defer n.Unlock()
	if n.conn == nil {
		c, err := newConn(n.addrs, n.opts.Secure, n.opts.TLSConfig)
		if err != nil {
			return nil, err
		}
		n.conn = c
	}
	return n.conn, nil
}

func (n *natsRegistry) register(s *registry.Service) error {
	conn, err := n.getConn()
	if err != nil {
		return err
	}

	n.Lock()
	defer n.Unlock()

	// cache service
	n.services[s.Name] = addServices(n.services[s.Name], []*registry.Service{s})

	// create query listener
	if n.listeners[s.Name] == nil {
		listener := make(chan bool)

		// create a subscriber that responds to queries
		sub, err := conn.Subscribe(QueryTopic, func(m *nats.Msg) {
			var result *registry.Result

			if err := json.Unmarshal(m.Data, &result); err != nil {
				return
			}

			var services []*registry.Service

			switch result.Action {
			// is this a get query and we own the service?
			case "get":
				if result.Service.Name != s.Name {
					return
				}
				n.RLock()
				services = n.services[s.Name]
				n.RUnlock()
			// it's a list request, but we're still only a
			// subscriber for this service... so just get this service
			// totally suboptimal
			case "list":
				n.RLock()
				services = n.services[s.Name]
				n.RUnlock()
			default:
				// does not match
				return
			}

			// respond to query
			for _, service := range services {
				b, err := json.Marshal(service)
				if err != nil {
					continue
				}
				conn.Publish(m.Reply, b)
			}
		})
		if err != nil {
			return err
		}

		// Unsubscribe if we're told to do so
		go func() {
			<-listener
			sub.Unsubscribe()
		}()

		n.listeners[s.Name] = listener
	}

	return nil
}

func (n *natsRegistry) deregister(s *registry.Service) error {
	n.Lock()
	defer n.Unlock()

	// cache leftover service
	services := addServices(n.services[s.Name], []*registry.Service{s})
	if len(services) > 0 {
		n.services[s.Name] = services
		return nil
	}

	// delete cached service
	delete(n.services, s.Name)

	// delete query listener
	if listener, lexists := n.listeners[s.Name]; lexists {
		close(listener)
		delete(n.listeners, s.Name)
	}

	return nil
}

func (n *natsRegistry) query(s string, quorum int) ([]*registry.Service, error) {
	conn, err := n.getConn()
	if err != nil {
		return nil, err
	}

	var action string
	var service *registry.Service

	if len(s) > 0 {
		action = "get"
		service = &registry.Service{Name: s}
	} else {
		action = "list"
	}

	inbox := nats.NewInbox()

	response := make(chan *registry.Service, 10)

	sub, err := conn.Subscribe(inbox, func(m *nats.Msg) {
		var service *registry.Service
		if err := json.Unmarshal(m.Data, &service); err != nil {
			return
		}
		select {
		case response <- service:
		case <-time.After(DefaultTimeout):
		}
	})
	if err != nil {
		return nil, err
	}
	defer sub.Unsubscribe()

	b, err := json.Marshal(&registry.Result{Action: action, Service: service})
	if err != nil {
		return nil, err
	}

	if err := conn.PublishMsg(&nats.Msg{
		Subject: QueryTopic,
		Reply:   inbox,
		Data:    b,
	}); err != nil {
		return nil, err
	}

	timeoutChan := time.After(n.opts.Timeout)

	serviceMap := make(map[string]*registry.Service)

loop:
	for {
		select {
		case service := <-response:
			key := service.Name + "-" + service.Version
			srv, ok := serviceMap[key]
			if ok {
				srv.Nodes = append(srv.Nodes, service.Nodes...)
				serviceMap[key] = srv
			} else {
				serviceMap[key] = service
			}

			if quorum > 0 && len(serviceMap[key].Nodes) >= quorum {
				break loop
			}
		case <-timeoutChan:
			break loop
		}
	}

	var services []*registry.Service
	for _, service := range serviceMap {
		services = append(services, service)
	}
	return services, nil
}

func (n *natsRegistry) Register(s *registry.Service, opts ...registry.RegisterOption) error {
	if err := n.register(s); err != nil {
		return err
	}

	conn, err := n.getConn()
	if err != nil {
		return err
	}

	b, err := json.Marshal(&registry.Result{Action: "create", Service: s})
	if err != nil {
		return err
	}

	return conn.Publish(WatchTopic, b)
}

func (n *natsRegistry) Deregister(s *registry.Service) error {
	if err := n.deregister(s); err != nil {
		return err
	}

	conn, err := n.getConn()
	if err != nil {
		return err
	}

	b, err := json.Marshal(&registry.Result{Action: "delete", Service: s})
	if err != nil {
		return err
	}
	return conn.Publish(WatchTopic, b)
}

func (n *natsRegistry) GetService(s string) ([]*registry.Service, error) {
	services, err := n.query(s, getQuorum(n.opts))
	if err != nil {
		return nil, err
	}
	return services, nil
}

func (n *natsRegistry) ListServices() ([]*registry.Service, error) {
	s, err := n.query("", 0)
	if err != nil {
		return nil, err
	}

	var services []*registry.Service
	serviceMap := make(map[string]*registry.Service)

	for _, v := range s {
		serviceMap[v.Name] = &registry.Service{Name: v.Name}
	}

	for _, v := range serviceMap {
		services = append(services, v)
	}

	return services, nil
}

func (n *natsRegistry) Watch() (registry.Watcher, error) {
	conn, err := n.getConn()
	if err != nil {
		return nil, err
	}

	sub, err := conn.SubscribeSync(WatchTopic)
	if err != nil {
		return nil, err
	}

	return &natsWatcher{sub}, nil
}

func (n *natsRegistry) String() string {
	return "nats"
}

func NewRegistry(opts ...registry.Option) registry.Registry {
	options := registry.Options{
		Timeout: DefaultTimeout,
	}
	for _, o := range opts {
		o(&options)
	}
	return &natsRegistry{
		addrs:     options.Addrs,
		opts:      options,
		services:  make(map[string][]*registry.Service),
		listeners: make(map[string]chan bool),
	}
}
