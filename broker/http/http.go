// Package http provides a http broker
package http

import (
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/broker/http"
)

/*
	HTTP Broker is the default broker in go-micro to reduce the number of dependencies.
	Find the implementation at https://godoc.org/github.com/micro/go-micro/broker/http.
	We add a link here for completeness
*/

func NewBroker(opts ...broker.Option) broker.Broker {
	return http.NewBroker(opts...)
}
