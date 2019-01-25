package rabbitmq

import (
	"context"

	"github.com/micro/go-micro/broker"
)

type durableQueueKey struct{}
type headersKey struct{}
type prefetchCountKey struct{}
type prefetchGlobalKey struct{}
type exchangeKey struct{}
type requeueOnErrorKey struct{}
type deliveryMode struct{}

// DurableQueue creates a durable queue when subscribing.
func DurableQueue() broker.SubscribeOption {
	return setSubscribeOption(durableQueueKey{}, true)
}

// Headers adds headers used by the headers exchange
func Headers(h map[string]interface{}) broker.SubscribeOption {
	return setSubscribeOption(headersKey{}, h)
}

// RequeueOnError calls Nack(muliple:false, requeue:true) on amqp delivery when handler returns error
func RequeueOnError() broker.SubscribeOption {
	return setSubscribeOption(requeueOnErrorKey{}, true)
}

// Exchange is an option to set the Exchange
func Exchange(e string) broker.Option {
	return setBrokerOption(exchangeKey{}, e)
}

// PrefetchCount ...
func PrefetchCount(c int) broker.Option {
	return setBrokerOption(prefetchCountKey{}, c)
}

// PrefetchGlobal creates a durable queue when subscribing.
func PrefetchGlobal() broker.Option {
	return setBrokerOption(prefetchGlobalKey{}, true)
}

// DeliveryMode sets a delivery mode for publishing
func DeliveryMode(value uint8) broker.PublishOption {
	return setPublishOption(deliveryMode{}, value)
}

type subscribeContextKey struct{}

// SubscribeContext set the context for broker.SubscribeOption
func SubscribeContext(ctx context.Context) broker.SubscribeOption {
	return setSubscribeOption(subscribeContextKey{}, ctx)
}

type successAutoAckKey struct{}

// SuccessAutoAck allow to AutoAck messages when handler returns without error
func SuccessAutoAck() broker.SubscribeOption {
	return setSubscribeOption(successAutoAckKey{}, true)
}
