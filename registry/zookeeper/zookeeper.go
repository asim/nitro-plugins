// Package zookeeper provides a zookeeper registry
package zookeeper

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"sync"
	"time"

	"github.com/go-zookeeper/zk"
	log "github.com/micro/go-micro/v3/logger"
	"github.com/micro/go-micro/v3/registry"
	hash "github.com/mitchellh/hashstructure"
)

var (
	prefix        = "/micro/registry"
	defaultDomain = "micro"
)

type zookeeperRegistry struct {
	client       *zk.Conn
	options      registry.Options
	register     map[string]register
	leases       map[string]leases
	ttlSupported bool
	sync.RWMutex
}

type serviceInfo struct {
	service     *registry.Service
	leaseID     int64
	leaseOK     bool
	pathKey     string
	opts        registry.RegisterOptions
	currentNode *registry.Node
}

type register map[string]uint64
type leases map[string]int64

func configure(z *zookeeperRegistry, opts ...registry.Option) error {
	cAddrs := z.options.Addrs

	for _, o := range opts {
		o(&z.options)
	}

	if z.options.Timeout == 0 {
		z.options.Timeout = 5
	}

	// already set
	if z.client != nil && len(z.options.Addrs) == len(cAddrs) {
		return nil
	}

	// reset
	cAddrs = nil

	for _, addr := range z.options.Addrs {
		if len(addr) == 0 {
			continue
		}
		cAddrs = append(cAddrs, addr)
	}

	if len(cAddrs) == 0 {
		cAddrs = []string{"127.0.0.1:2181"}
	}

	// connect to zookeeper
	c, _, err := zk.Connect(cAddrs, time.Second*z.options.Timeout)
	if err != nil {
		log.Errorf("connect to zk err: %s", err)
		return err
	}

	// create our prefix path
	if err := createPath(prefix, []byte{}, c, -1); err != nil {
		log.Errorf("create zk path err: %s", err)
		return err
	}

	z.client = c
	z.ttlSupported = z.checkSupportTTL()

	return nil
}

func (z *zookeeperRegistry) Init(opts ...registry.Option) error {
	return configure(z, opts...)
}

func (z *zookeeperRegistry) Options() registry.Options {
	return z.options
}

func (z *zookeeperRegistry) Deregister(s *registry.Service, opts ...registry.DeregisterOption) error {
	if len(s.Nodes) == 0 {
		return errors.New("Require at least one currentNode")
	}

	// parse the options
	var options registry.DeregisterOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.Domain) == 0 {
		options.Domain = defaultDomain
	}

	for _, node := range s.Nodes {
		z.Lock()
		// delete our hash of the service
		nodes, ok := z.register[options.Domain]
		if ok {
			delete(nodes, s.Name+node.Id)
			z.register[options.Domain] = nodes
		}

		// delete our lease of the service
		leases, ok := z.leases[options.Domain]
		if ok {
			delete(leases, s.Name+node.Id)
			z.leases[options.Domain] = leases
		}
		z.Unlock()

		if log.V(log.TraceLevel, log.DefaultLogger) {
			log.Tracef("Deregister %s id %s", s.Name, node.Id)
		}

		if err := z.client.Delete(nodePath(options.Domain, s.Name, node.Id), -1); err != nil {
			return err
		}
	}

	return nil
}

func (z *zookeeperRegistry) Register(s *registry.Service, opts ...registry.RegisterOption) error {
	if len(s.Nodes) == 0 {
		return errors.New("Require at least one currentNode ")
	}

	var gerr error

	// register each currentNode individually
	for _, node := range s.Nodes {
		if err := z.registerNode(s, node, opts...); err != nil {
			gerr = err
		}
	}

	return gerr
}

func (z *zookeeperRegistry) GetService(name string, opts ...registry.GetOption) ([]*registry.Service, error) {
	var options registry.GetOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.Domain) == 0 {
		options.Domain = defaultDomain
	}

	path := servicePath(options.Domain, name)
	l, _, err := z.client.Children(path)
	if err != nil {
		return nil, fmt.Errorf("name: [%s] path: [%s] err: [%s] ", name, path, err)
	}

	serviceMap := make(map[string]*registry.Service)

	for _, n := range l {
		_, stat, err := z.client.Children(nodePath(options.Domain, name, n))
		if err != nil {
			log.Errorf("get children err: %s", err)
			return nil, err
		}

		if stat.NumChildren > 0 {
			continue
		}

		b, _, err := z.client.Get(nodePath(options.Domain, name, n))
		if err != nil {
			log.Errorf("get currentNode path err: %s", err)
			return nil, err
		}

		sn, err := decode(b)
		if err != nil {
			log.Errorf("decode currentNode data err: %s", err)
			return nil, err
		}

		s, ok := serviceMap[sn.Version]
		if !ok {
			s = &registry.Service{
				Name:      sn.Name,
				Version:   sn.Version,
				Metadata:  sn.Metadata,
				Endpoints: sn.Endpoints,
			}
			serviceMap[s.Version] = s
		}

		for _, node := range sn.Nodes {
			s.Nodes = append(s.Nodes, node)
		}
	}

	services := make([]*registry.Service, 0, len(serviceMap))

	for _, service := range serviceMap {
		services = append(services, service)
	}

	return services, nil
}

func (z *zookeeperRegistry) ListServices(opts ...registry.ListOption) ([]*registry.Service, error) {
	var options registry.ListOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.Domain) == 0 {
		options.Domain = defaultDomain
	}

	var p = prefixWithDomain(options.Domain)
	srv, _, err := z.client.Children(p)
	if err != nil {
		log.Errorf("list get children err: %s", err)
		return nil, err
	}

	serviceMap := make(map[string]*registry.Service)

	for _, key := range srv {
		s := servicePath(options.Domain, key)
		nodes, _, err := z.client.Children(s)
		if err != nil {
			return nil, err
		}

		for _, node := range nodes {
			_, stat, err := z.client.Children(nodePath(options.Domain, key, node))
			if err != nil {
				return nil, err
			}

			if stat.NumChildren == 0 {
				b, _, err := z.client.Get(nodePath(options.Domain, key, node))
				if err != nil {
					return nil, err
				}
				i, err := decode(b)
				if err != nil {
					return nil, err
				}
				serviceMap[s] = &registry.Service{Name: i.Name}
			}
		}
	}

	var services []*registry.Service

	for _, service := range serviceMap {
		services = append(services, service)
	}

	return services, nil
}

func (z *zookeeperRegistry) String() string {
	return "zookeeper"
}

func (z *zookeeperRegistry) Watch(opts ...registry.WatchOption) (registry.Watcher, error) {
	return newZookeeperWatcher(z, opts...)
}

func NewRegistry(opts ...registry.Option) registry.Registry {
	z := &zookeeperRegistry{
		options:  registry.Options{},
		register: make(map[string]register),
		leases:   make(map[string]leases),
	}

	if err := configure(z, opts...); err != nil {
		log.Errorf("configure zk err: %s", err)
	}

	return z
}

func (z *zookeeperRegistry) registerNode(s *registry.Service, node *registry.Node, opts ...registry.RegisterOption) error {
	if len(s.Nodes) == 0 {
		return errors.New("Require at least one currentNode ")
	}

	si := z.prepareService(s, node, opts...)
	srv, _ := encode(si.service)

	if z.ttlSupported {
		z.registerWithTTL(si, srv)
	} else {
		z.registerWithoutTTL(si, srv)
	}

	return nil
}

func (z *zookeeperRegistry) registerWithoutTTL(si *serviceInfo, srv []byte) (err error) {
	// create hash of service; uint64
	h, err := hash.Hash(si.service, nil)
	if err != nil {
		return
	}

	// get existing hash
	z.Lock()
	v, ok := z.register[si.opts.Domain][si.service.Name+si.currentNode.Id]
	z.Unlock()

	// the service is unchanged, skip registering
	if ok && v == h {
		return nil
	}

	exists, _, err := z.client.Exists(si.pathKey)
	if err != nil {
		return err
	}

	if exists {
		_, err := z.client.Set(si.pathKey, srv, -1)
		if err != nil {
			return err
		}
	} else {
		err := createPath(si.pathKey, srv, z.client, -1)
		if err != nil {
			return err
		}
	}

	// save our hash of the service
	z.Lock()
	z.register[si.opts.Domain][si.service.Name+si.currentNode.Id] = h
	z.Unlock()

	return
}

func (z *zookeeperRegistry) registerWithTTL(si *serviceInfo, srv []byte) {
	if !si.leaseOK {
		// look for the existing key
		bytes, stat, err := z.client.Get(si.pathKey)
		if err != nil {
			log.Infof("[register currentNode] get key [%s] err: %s", si.pathKey, err)
		} else if bytes != nil { // currentNode exits
			si.leaseID = stat.Mzxid

			// decode the existing currentNode
			srv, err := decode(bytes)
			if err != nil {
				// dont return
				log.Error(err)
			} else {
				// create hash of service; uint64
				h, err := hash.Hash(srv.Nodes[0], nil)
				if err != nil {
					log.Error(err)
				} else {
					// save the info
					z.Lock()
					z.leases[si.opts.Domain][si.service.Name+si.currentNode.Id] = si.leaseID
					z.register[si.opts.Domain][si.service.Name+si.currentNode.Id] = h
					z.Unlock()
				}
			}
		}
	}

	if si.leaseID > 0 {
		// renew the lease if it exists
		if log.V(log.TraceLevel, log.DefaultLogger) {
			log.Tracef("Renewing existing lease for %s %d", si.service.Name, si.leaseID)
		}

		// delete first, sdk doest have setTTL method.
		err := z.client.Delete(si.pathKey, -1)
		if err != nil {
			log.Errorf("delete before renew key: [%s] err: %", si.pathKey, err)
			return
		}

		err = createPath(si.pathKey, srv, z.client, si.opts.TTL)
		if err != nil {
			log.Errorf("renew after delete key: [%s] err: %", si.pathKey, err)
			return
		}
	}

	z.Lock()
	z.leases[si.opts.Domain][si.service.Name+si.currentNode.Id] = si.leaseID
	z.Unlock()
}

func (z *zookeeperRegistry) prepareService(s *registry.Service, node *registry.Node, opts ...registry.RegisterOption) *serviceInfo {
	// parse the options
	var options registry.RegisterOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.Domain) == 0 {
		options.Domain = defaultDomain
	}

	if s.Metadata == nil {
		s.Metadata = map[string]string{}
	}
	s.Metadata["domain"] = options.Domain

	if node.Metadata == nil {
		node.Metadata = map[string]string{}
	}
	node.Metadata["domain"] = options.Domain

	z.Lock()
	if _, ok := z.register[options.Domain]; !ok {
		z.register[options.Domain] = make(register)
	}

	// ensure the leases and registers are setup for this domain
	if _, ok := z.leases[options.Domain]; !ok && z.ttlSupported {
		z.leases[options.Domain] = make(leases)
	}

	// check to see if we already have a lease cached
	leaseID, ok := z.leases[options.Domain][s.Name+node.Id]
	z.Unlock()

	si := &serviceInfo{
		service: &registry.Service{
			Name:      s.Name,
			Version:   s.Version,
			Metadata:  s.Metadata,
			Endpoints: s.Endpoints,
			Nodes:     []*registry.Node{node},
		},
		currentNode: node,
		leaseID:     leaseID,
		leaseOK:     ok,
		pathKey:     nodePath(options.Domain, s.Name, node.Id),
		opts:        options,
	}

	return si
}

// checkSupportTTL zk doest support TTL currentNode by default. we need to check it.
func (z *zookeeperRegistry) checkSupportTTL() bool {
	tempID := uuid.New().String()
	key := nodePath(defaultDomain, tempID, tempID)
	_, err := z.client.CreateTTL(key, []byte{}, zk.FlagTTL, nil, time.Second)
	if err != nil && strings.Contains(err.Error(), "-6") {
		log.Infof("[check ttl] this zk doesn't support TTL")
		return false
	}

	// delete it as soon as possible and ignore the error
	// because if success above, it will be deleted in one second.
	_ = z.client.Delete(key, -1)

	return true
}
