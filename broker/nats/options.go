package nats

import (
	"github.com/micro/go-micro/broker"
	nats "github.com/nats-io/go-nats"
)

type optionsKey struct{}

// Options accepts nats.Options
func Options(opts nats.Options) broker.Option {
	return setBrokerOption(optionsKey{}, opts)
}
