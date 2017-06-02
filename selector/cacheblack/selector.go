package cacheblack

import (
	"log"
	"sync"
	"time"

	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/selector"
)

type cacheBlackSelector struct {
	so       selector.Options
	cacheTTL time.Duration

	// registry cache
	mu    sync.Mutex
	cache map[string][]*registry.Service
	ttls  map[string]time.Time

	// used to close or reload watcher
	reload chan bool
	exit   chan bool

	bl *blacklist
}

func (self *cacheBlackSelector) Init(opts ...selector.Option) error {
	for _, o := range opts {
		o(&self.so)
	}

	self.bl.Close()
	self.bl = newOptionBlacklist(&self.so)

	// reload the watcher
	go func() {
		select {
		case <-self.exit:
			return
		default:
			self.reload <- true
		}
	}()

	return nil
}

func (self *cacheBlackSelector) Options() selector.Options {
	return self.so
}

func (self *cacheBlackSelector) quit() bool {
	select {
	case <-self.exit:
		return true
	default:
		return false
	}
}

// cp copies a service. Because we're caching handing back pointers would
// create a race condition, so we do this instead
// its fast enough
func (self *cacheBlackSelector) cp(current []*registry.Service) []*registry.Service {
	var services []*registry.Service

	for _, service := range current {
		// copy service
		s := new(registry.Service)
		*s = *service

		// copy nodes
		var nodes []*registry.Node
		for _, node := range service.Nodes {
			n := new(registry.Node)
			*n = *node
			nodes = append(nodes, n)
		}
		s.Nodes = nodes

		// copy endpoints
		var eps []*registry.Endpoint
		for _, ep := range service.Endpoints {
			e := new(registry.Endpoint)
			*e = *ep
			eps = append(eps, e)
		}
		s.Endpoints = eps

		// append service
		services = append(services, s)
	}

	return services
}

func (self *cacheBlackSelector) del(service string) {
	delete(self.cache, service)
	delete(self.ttls, service)
}

func (self *cacheBlackSelector) get(service string) ([]*registry.Service, error) {
	self.mu.Lock()
	defer self.mu.Unlock()

	// check the cache first
	services, ok := self.cache[service]
	ttl, kk := self.ttls[service]

	// got results, copy and return
	if ok && len(services) > 0 {
		// only return if its less than the ttl
		if kk && time.Since(ttl) < self.cacheTTL {
			return self.cp(services), nil
		}
	}

	// cache miss or ttl expired

	// now ask the registry
	services, err := self.so.Registry.GetService(service)
	if err != nil {
		return nil, err
	}

	// we didn't have any results so cache
	self.cache[service] = self.cp(services)
	self.ttls[service] = time.Now().Add(self.cacheTTL)
	return services, nil
}

func (self *cacheBlackSelector) set(service string, services []*registry.Service) {
	self.cache[service] = services
	self.ttls[service] = time.Now().Add(self.cacheTTL)
}

func (self *cacheBlackSelector) update(res *registry.Result) {
	if res == nil || res.Service == nil {
		return
	}

	self.mu.Lock()
	defer self.mu.Unlock()

	services, ok := self.cache[res.Service.Name]
	if !ok {
		// we're not going to cache anything
		// unless there was already a lookup
		return
	}

	if len(res.Service.Nodes) == 0 {
		switch res.Action {
		case "delete":
			self.del(res.Service.Name)
		}
		return
	}

	// existing service found
	var service *registry.Service
	var index int
	for i, s := range services {
		if s.Version == res.Service.Version {
			service = s
			index = i
		}
	}

	switch res.Action {
	case "create", "update":
		if service == nil {
			self.set(res.Service.Name, append(services, res.Service))
			return
		}

		// append old nodes to new service
		for _, cur := range service.Nodes {
			var seen bool
			for _, node := range res.Service.Nodes {
				if cur.Id == node.Id {
					seen = true
					break
				}
			}
			if !seen {
				res.Service.Nodes = append(res.Service.Nodes, cur)
			}
		}

		services[index] = res.Service
		self.set(res.Service.Name, services)
	case "delete":
		if service == nil {
			return
		}

		var nodes []*registry.Node

		// filter cur nodes to remove the dead one
		for _, cur := range service.Nodes {
			var seen bool
			for _, del := range res.Service.Nodes {
				if del.Id == cur.Id {
					seen = true
					break
				}
			}
			if !seen {
				nodes = append(nodes, cur)
			}
		}

		// still got nodes, save and return
		if len(nodes) > 0 {
			service.Nodes = nodes
			services[index] = service
			self.set(service.Name, services)
			return
		}

		// zero nodes left

		// only have one thing to delete
		// nuke the thing
		if len(services) == 1 {
			self.del(service.Name)
			return
		}

		// still have more than 1 service
		// check the version and keep what we know
		var srvs []*registry.Service
		for _, s := range services {
			if s.Version != service.Version {
				srvs = append(srvs, s)
			}
		}

		// save
		self.set(service.Name, srvs)
	}
}

// run starts the cache watcher loop
// it creates a new watcher if there's a problem
// reloads the watcher if Init is called
// and returns when Close is called
func (self *cacheBlackSelector) run() {
	go self.tick()

	for {
		// exit early if already dead
		if self.quit() {
			return
		}

		// create new watcher
		w, err := self.so.Registry.Watch()
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second)
			continue
		}

		// watch for events
		if err := self.watch(w); err != nil {
			log.Println(err)
			continue
		}
	}
}

// check cache and expire on each tick
func (self *cacheBlackSelector) tick() {
	t := time.NewTicker(time.Minute)

	for {
		select {
		case <-t.C:
			self.mu.Lock()
			for service, expiry := range self.ttls {
				if d := time.Since(expiry); d > self.cacheTTL {
					// TODO: maybe refresh the cache rather than blowing it away
					self.del(service)
				}
			}
			self.mu.Unlock()
		case <-self.exit:
			return
		}
	}
}

// watch loops the next event and calls update
// it returns if there's an error
func (self *cacheBlackSelector) watch(w registry.Watcher) error {
	defer w.Stop()

	// manage this loop
	go func() {
		// wait for exit or reload signal
		select {
		case <-self.exit:
		case <-self.reload:
		}

		// stop the watcher
		w.Stop()
	}()

	for {
		res, err := w.Next()
		if err != nil {
			return err
		}
		self.update(res)
	}
}

func (self *cacheBlackSelector) Select(service string, opts ...selector.SelectOption) (selector.Next, error) {
	sopts := selector.SelectOptions{
		Strategy: self.so.Strategy,
	}

	for _, opt := range opts {
		opt(&sopts)
	}

	// get the service
	// try the cache first
	// if that fails go directly to the registry
	services, err := self.get(service)
	if err != nil {
		return nil, err
	}

	// apply the filters
	for _, filter := range sopts.Filters {
		services = filter(services)
	}

	// apply the blacklist
	services, err = self.bl.Filter(services)
	if err != nil {
		return nil, err
	}

	// if there's nothing left, return
	if len(services) == 0 {
		return nil, selector.ErrNoneAvailable
	}

	return sopts.Strategy(services), nil
}

func (self *cacheBlackSelector) Mark(service string, node *registry.Node, err error) {
	self.bl.Mark(service, node, err)
}

func (self *cacheBlackSelector) Reset(service string) {
	self.bl.Reset(service)
}

// Close stops the watcher and destroys the cache
func (self *cacheBlackSelector) Close() error {
	self.mu.Lock()
	self.cache = make(map[string][]*registry.Service)
	self.mu.Unlock()

	select {
	case <-self.exit:
		return nil
	default:
		close(self.exit)
		self.bl.Close()
	}
	return nil
}

func (self *cacheBlackSelector) String() string {
	return "cacheblack"
}

func newOptionBlacklist(options *selector.Options) *blacklist {
	blackTTL := DefaultBlackTTL
	if options.Context != nil {
		if t, ok := options.Context.Value(blackTTLKey{}).(int); ok {
			blackTTL = t
		}
	}

	blackErrCount := DefaultBlackErrCount
	if options.Context != nil {
		if t, ok := options.Context.Value(blackErrCountKey{}).(int); ok {
			blackErrCount = t
		}
	}

	return newBlacklist(blackTTL, blackErrCount)
}

func NewSelector(opts ...selector.Option) selector.Selector {
	options := selector.Options{
		Strategy: selector.Random,
	}

	for _, opt := range opts {
		opt(&options)
	}

	if options.Registry == nil {
		options.Registry = registry.DefaultRegistry
	}

	cacheTTL := DefaultCacheTTL
	if options.Context != nil {
		if t, ok := options.Context.Value(cacheTTLKey{}).(time.Duration); ok {
			cacheTTL = t
		}
	}

	c := &cacheBlackSelector{
		so:       options,
		cacheTTL: cacheTTL,
		cache:    make(map[string][]*registry.Service),
		ttls:     make(map[string]time.Time),
		reload:   make(chan bool, 1),
		exit:     make(chan bool),
		bl:       newOptionBlacklist(&options),
	}

	go c.run()
	return c
}
