package nsq

import (
	"context"

	"github.com/micro/go-micro/broker"
)

type contextKeyT string

var (
	concurrentHandlerKey = contextKeyT("github.com/saturn-sonic/go-plugins/broker/nsq/concurrentHandlers")
)

func ConcurrentHandlers(n int) broker.SubscribeOption {
	return func(o *broker.SubscribeOptions) {
		o.Context = context.WithValue(o.Context, concurrentHandlerKey, n)
	}
}
