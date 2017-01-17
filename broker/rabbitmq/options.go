package rabbitmq

import (
	"github.com/micro/go-micro/broker"
	"golang.org/x/net/context"
)

type durableQueueKey struct{}

// DurableQueue creates a durable queue when subscribing.
func DurableQueue() broker.SubscribeOption {
	return func(o *broker.SubscribeOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, durableQueueKey{}, true)
	}
}
