// Package mdns provides a multicast DNS registry
package mdns

import (
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/mdns"
)

/*
	MDNS registry is one of the default registries in go-micro.
	It supports multicast DNS service discovery.
	Implementation here https://godoc.org/github.com/micro/go-micro/registry/mdns
	We add a link for completeness
*/

func NewRegistry(opts ...registry.Option) registry.Registry {
	return mdns.NewRegistry(opts...)
}
