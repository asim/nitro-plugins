// Package Gossip provides a gossip registry based on hashicorp/memberlist
package gossip

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/memberlist"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/registry"
	"github.com/mitchellh/hashstructure"
	"github.com/pborman/uuid"
)

type action int

const (
	addAction action = iota
	delAction
	syncAction
)

type broadcast struct {
	msg    []byte
	notify chan<- struct{}
}

type delegate struct {
	broadcasts *memberlist.TransmitLimitedQueue
	updates    chan *update
}

type gossipRegistry struct {
	broadcasts *memberlist.TransmitLimitedQueue
	updates    chan *update

	sync.RWMutex
	services map[string][]*registry.Service

	s    sync.RWMutex
	subs map[string]chan *registry.Result
}

type update struct {
	Action    action
	Service   *registry.Service
	Timestamp int64
	Expires   int64
	sync      chan *registry.Service
}

var (
	// You should change this if using secure
	DefaultKey = []byte("gossipKey")
	ExpiryTick = time.Second * 10
)

func init() {
	cmd.DefaultRegistries["gossip"] = NewRegistry
}

func addNodes(old, neu []*registry.Node) []*registry.Node {
	for _, n := range neu {
		var seen bool
		for i, o := range old {
			if o.Id == n.Id {
				seen = true
				old[i] = n
				break
			}
		}
		if !seen {
			old = append(old, n)
		}
	}
	return old
}

func addServices(old, neu []*registry.Service) []*registry.Service {
	for _, s := range neu {
		var seen bool
		for i, o := range old {
			if o.Version == s.Version {
				s.Nodes = addNodes(o.Nodes, s.Nodes)
				seen = true
				old[i] = s
				break
			}
		}
		if !seen {
			old = append(old, s)
		}
	}
	return old
}

func delNodes(old, del []*registry.Node) []*registry.Node {
	var nodes []*registry.Node
	for _, o := range old {
		var rem bool
		for _, n := range del {
			if o.Id == n.Id {
				rem = true
				break
			}
		}
		if !rem {
			nodes = append(nodes, o)
		}
	}
	return nodes
}

func delServices(old, del []*registry.Service) []*registry.Service {
	var services []*registry.Service
	for i, o := range old {
		var rem bool
		for _, s := range del {
			if o.Version == s.Version {
				old[i].Nodes = delNodes(o.Nodes, s.Nodes)
				if len(old[i].Nodes) == 0 {
					rem = true
				}
			}
		}
		if !rem {
			services = append(services, o)
		}
	}
	return services
}

func (b *broadcast) Invalidates(other memberlist.Broadcast) bool {
	return false
}

func (b *broadcast) Message() []byte {
	return b.msg
}

func (b *broadcast) Finished() {
	if b.notify != nil {
		close(b.notify)
	}
}

func (d *delegate) NodeMeta(limit int) []byte {
	return []byte{}
}

func (d *delegate) NotifyMsg(b []byte) {
	if len(b) == 0 {
		return
	}

	buf := make([]byte, len(b))
	copy(buf, b)

	go func() {
		switch buf[0] {
		case 'd': // data
			var updates []*update
			if err := json.Unmarshal(buf[1:], &updates); err != nil {
				return
			}
			for _, u := range updates {
				d.updates <- u
			}
		}
	}()
}

func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return d.broadcasts.GetBroadcasts(overhead, limit)
}

func (d *delegate) LocalState(join bool) []byte {
	if !join {
		return []byte{}
	}

	syncCh := make(chan *registry.Service, 1)
	m := map[string][]*registry.Service{}

	d.updates <- &update{
		Action: syncAction,
		sync:   syncCh,
	}

	for s := range syncCh {
		m[s.Name] = append(m[s.Name], s)
	}

	b, _ := json.Marshal(m)
	return b
}

func (d *delegate) MergeRemoteState(buf []byte, join bool) {
	if len(buf) == 0 {
		return
	}
	if !join {
		return
	}

	var m map[string][]*registry.Service
	if err := json.Unmarshal(buf, &m); err != nil {
		return
	}

	for _, services := range m {
		for _, service := range services {
			d.updates <- &update{
				Action:  addAction,
				Service: service,
				sync:    nil,
			}
		}
	}
}

func (m *gossipRegistry) publish(action string, services []*registry.Service) {
	m.s.RLock()
	for _, sub := range m.subs {
		go func() {
			for _, service := range services {
				sub <- &registry.Result{Action: action, Service: service}
			}
		}()
	}
	m.s.RUnlock()
}

func (m *gossipRegistry) subscribe() (chan *registry.Result, chan bool) {
	next := make(chan *registry.Result, 10)
	exit := make(chan bool)

	id := uuid.NewUUID().String()

	m.s.Lock()
	m.subs[id] = next
	m.s.Unlock()

	go func() {
		<-exit
		m.s.Lock()
		delete(m.subs, id)
		close(next)
		m.s.Unlock()
	}()

	return next, exit
}

func (m *gossipRegistry) run() {
	var mtx sync.Mutex
	updates := map[uint64]*update{}

	// expiry loop
	go func() {
		t := time.NewTicker(ExpiryTick)
		defer t.Stop()

		for _ = range t.C {
			now := time.Now().Unix()

			mtx.Lock()
			for k, v := range updates {
				// check if expiry time has passed
				if d := (v.Timestamp + v.Expires) - now; d < 0 {
					// delete from records
					delete(updates, k)
					// set to delete
					v.Action = delAction
					// fire a new update
					m.updates <- v
				}
			}
			mtx.Unlock()
		}
	}()

	for u := range m.updates {
		switch u.Action {
		case addAction:
			m.Lock()
			if service, ok := m.services[u.Service.Name]; !ok {
				m.services[u.Service.Name] = []*registry.Service{u.Service}

			} else {
				m.services[u.Service.Name] = addServices(service, []*registry.Service{u.Service})
			}
			m.Unlock()
			go m.publish("add", []*registry.Service{u.Service})

			// we need to expire the node at some point in the future
			if u.Expires > 0 {
				// create a hash of this service
				if hash, err := hashstructure.Hash(u.Service, nil); err == nil {
					mtx.Lock()
					updates[hash] = u
					mtx.Unlock()
				}
			}
		case delAction:
			m.Lock()
			if service, ok := m.services[u.Service.Name]; ok {
				if services := delServices(service, []*registry.Service{u.Service}); len(services) == 0 {
					delete(m.services, u.Service.Name)
				} else {
					m.services[u.Service.Name] = services
				}
			}
			m.Unlock()
			go m.publish("delete", []*registry.Service{u.Service})

			// delete from expiry checks
			if hash, err := hashstructure.Hash(u.Service, nil); err == nil {
				mtx.Lock()
				delete(updates, hash)
				mtx.Unlock()
			}
		case syncAction:
			if u.sync == nil {
				continue
			}
			m.RLock()
			for _, services := range m.services {
				for _, service := range services {
					u.sync <- service
				}
				go m.publish("add", services)
			}
			m.RUnlock()
			close(u.sync)
		}
	}
}

func (m *gossipRegistry) Register(s *registry.Service, opts ...registry.RegisterOption) error {
	m.Lock()
	if service, ok := m.services[s.Name]; !ok {
		m.services[s.Name] = []*registry.Service{s}
	} else {
		m.services[s.Name] = addServices(service, []*registry.Service{s})
	}
	m.Unlock()

	var options registry.RegisterOptions
	for _, o := range opts {
		o(&options)
	}

	b, _ := json.Marshal([]*update{
		&update{
			Action:    addAction,
			Service:   s,
			Timestamp: time.Now().Unix(),
			Expires:   int64(options.TTL.Seconds()),
		},
	})

	m.broadcasts.QueueBroadcast(&broadcast{
		msg:    append([]byte("d"), b...),
		notify: nil,
	})

	return nil
}

func (m *gossipRegistry) Deregister(s *registry.Service) error {
	m.Lock()
	if service, ok := m.services[s.Name]; ok {
		if services := delServices(service, []*registry.Service{s}); len(services) == 0 {
			delete(m.services, s.Name)
		} else {
			m.services[s.Name] = services
		}
	}
	m.Unlock()

	b, _ := json.Marshal([]*update{
		&update{
			Action:  delAction,
			Service: s,
		},
	})

	m.broadcasts.QueueBroadcast(&broadcast{
		msg:    append([]byte("d"), b...),
		notify: nil,
	})

	return nil
}

func (m *gossipRegistry) GetService(name string) ([]*registry.Service, error) {
	m.RLock()
	service, ok := m.services[name]
	m.RUnlock()
	if !ok {
		return nil, fmt.Errorf("Service %s not found", name)
	}
	return service, nil
}

func (m *gossipRegistry) ListServices() ([]*registry.Service, error) {
	var services []*registry.Service
	m.RLock()
	for _, service := range m.services {
		services = append(services, service...)
	}
	m.RUnlock()
	return services, nil
}

func (m *gossipRegistry) Watch() (registry.Watcher, error) {
	n, e := m.subscribe()
	return newGossipWatcher(n, e)
}

func (m *gossipRegistry) String() string {
	return "gossip"
}

func NewRegistry(opts ...registry.Option) registry.Registry {
	var options registry.Options
	for _, o := range opts {
		o(&options)
	}

	cAddrs := []string{}
	hostname, _ := os.Hostname()
	updates := make(chan *update, 100)

	for _, addr := range options.Addrs {
		if len(addr) > 0 {
			cAddrs = append(cAddrs, addr)
		}
	}

	broadcasts := &memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			return len(cAddrs)
		},
		RetransmitMult: 3,
	}

	mr := &gossipRegistry{
		broadcasts: broadcasts,
		services:   make(map[string][]*registry.Service),
		updates:    updates,
		subs:       make(map[string]chan *registry.Result),
	}

	go mr.run()

	c := memberlist.DefaultLocalConfig()
	c.BindPort = 0
	c.Name = hostname + "-" + uuid.NewUUID().String()
	c.Delegate = &delegate{
		updates:    updates,
		broadcasts: broadcasts,
	}

	if options.Secure {
		k, ok := options.Context.Value(contextSecretKey{}).([]byte)
		if !ok {
			k = DefaultKey
		}
		c.SecretKey = k
	}

	m, err := memberlist.Create(c)
	if err != nil {
		log.Fatalf("Error creating memberlist: %v", err)
	}

	if len(cAddrs) > 0 {
		_, err := m.Join(cAddrs)
		if err != nil {
			log.Fatalf("Error joining members: %v", err)
		}
	}

	log.Printf("Local memberlist node %s:%d\n", m.LocalNode().Addr, m.LocalNode().Port)
	return mr
}
