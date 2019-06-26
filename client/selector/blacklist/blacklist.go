// Package blacklist is a selector which includes blacklisting of nodes when they fail
package blacklist

import (
	"math/rand"
	"time"

	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/selector"
)

type blacklistSelector struct {
	so   selector.Options
	exit chan bool
	bl   *blacklist
}

func init() {
	rand.Seed(time.Now().Unix())
}

func (r *blacklistSelector) Init(opts ...selector.Option) error {
	for _, o := range opts {
		o(&r.so)
	}
	return nil
}

func (r *blacklistSelector) Options() selector.Options {
	return r.so
}

func (r *blacklistSelector) Select(service string, opts ...selector.SelectOption) (selector.Next, error) {
	sopts := selector.SelectOptions{
		Strategy: r.so.Strategy,
	}

	for _, opt := range opts {
		opt(&sopts)
	}

	// get the service
	services, err := r.so.Registry.GetService(service)
	if err != nil {
		return nil, err
	}

	// apply the filters
	for _, filter := range sopts.Filters {
		services = filter(services)
	}

	// apply the blacklist
	services, err = r.bl.Filter(services)
	if err != nil {
		return nil, err
	}

	// if there's nothing left, return
	if len(services) == 0 {
		return nil, selector.ErrNoneAvailable
	}

	return sopts.Strategy(services), nil
}

func (r *blacklistSelector) Mark(service string, node *registry.Node, err error) {
	r.bl.Mark(service, node, err)
}

func (r *blacklistSelector) Reset(service string) {
	r.bl.Reset(service)
}

func (r *blacklistSelector) Close() error {
	select {
	case <-r.exit:
		return nil
	default:
		close(r.exit)
		r.bl.Close()
	}
	return nil
}

func (r *blacklistSelector) String() string {
	return "blacklist"
}

func newSelector(opts ...selector.Option) selector.Selector {
	sopts := selector.Options{
		Strategy: selector.Random,
	}

	for _, opt := range opts {
		opt(&sopts)
	}

	if sopts.Registry == nil {
		sopts.Registry = registry.DefaultRegistry
	}

	return &blacklistSelector{
		so:   sopts,
		exit: make(chan bool),
		bl:   newBlacklist(),
	}
}

func NewSelector(opts ...selector.Option) selector.Selector {
	return newSelector(opts...)
}
