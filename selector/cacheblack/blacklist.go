package cacheblack

import (
	"math/rand"
	"sync"
	"time"

	"github.com/micro/go-micro/registry"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type node struct {
	age     time.Time
	id      string
	service string
	count   int
}

type blacklist struct {
	ttl   int // the ttl to blacklist for
	count int // number of times we see an error before blacklisting
	exit  chan bool

	sync.RWMutex
	bl map[string]node
}

func (self *blacklist) purge() {
	now := time.Now()
	self.Lock()
	for k, v := range self.bl {
		if d := v.age.Sub(now); d.Seconds() < 0 {
			delete(self.bl, k)
		}
	}
	self.Unlock()
}

func (self *blacklist) run() {
	t := time.NewTicker(time.Duration(self.ttl) * time.Second)

	for {
		select {
		case <-self.exit:
			t.Stop()
			return
		case <-t.C:
			self.purge()
		}
	}
}

func (self *blacklist) Filter(services []*registry.Service) ([]*registry.Service, error) {
	var viableServices []*registry.Service

	self.RLock()

	for _, service := range services {
		var viableNodes []*registry.Node

		for _, node := range service.Nodes {
			n, ok := self.bl[node.Id]
			if !ok {
				// blacklist miss so add it
				viableNodes = append(viableNodes, node)
				continue
			}

			// got some blacklist info
			// skip the node if it exceeds count
			if n.count >= self.count {
				continue
			}

			// doesn't exceed count, still viable
			viableNodes = append(viableNodes, node)
		}

		if len(viableNodes) == 0 {
			continue
		}

		viableService := new(registry.Service)
		*viableService = *service
		viableService.Nodes = viableNodes
		viableServices = append(viableServices, viableService)
	}

	self.RUnlock()

	return viableServices, nil
}

func (self *blacklist) Mark(service string, nod *registry.Node, err error) {
	self.Lock()
	defer self.Unlock()

	// reset when error is nil
	// basically closing the circuit
	if err == nil {
		delete(self.bl, nod.Id)
		return
	}

	n, ok := self.bl[nod.Id]
	if !ok {
		n = node{
			id:      nod.Id,
			service: service,
		}
	}

	// mark it
	n.count++

	// set age to ttl seconds in future
	n.age = time.Now().Add(time.Duration(self.ttl) * time.Second)

	// save
	self.bl[nod.Id] = n
}

func (self *blacklist) Reset(service string) {
	self.Lock()
	defer self.Unlock()

	for k, v := range self.bl {
		// delete every node that matches the service
		if v.service == service {
			delete(self.bl, k)
		}
	}
}

func (self *blacklist) Close() error {
	select {
	case <-self.exit:
		return nil
	default:
		close(self.exit)
	}
	return nil
}

func newBlacklist(ttl, errCount int) *blacklist {
	bl := &blacklist{
		ttl:   ttl,
		count: errCount,
		bl:    make(map[string]node),
		exit:  make(chan bool),
	}

	go bl.run()
	return bl
}
