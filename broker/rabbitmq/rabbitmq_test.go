package rabbitmq_test

import (
	"context"
	"testing"

	micro "github.com/micro/go-micro"
	broker "github.com/micro/go-micro/broker"
	server "github.com/micro/go-micro/server"
	rabbitmq "github.com/micro/go-plugins/broker/rabbitmq"
)

type Example struct{}

func (e *Example) Handler(ctx context.Context, r interface{}) error {
	panic("TE")
	return nil
}

func TestDurable(t *testing.T) {
	rabbitmq.DefaultRabbitURL = "amqp://rabbitmq:rabbitmq@172.18.0.2:5672"
	brkrSub := broker.NewSubscribeOptions(
		broker.Queue("queue.default"),
		broker.DisableAutoAck(),
		rabbitmq.DurableQueue(),
	)

	b := rabbitmq.NewBroker()
	b.Init()
	b.Connect()
	s := server.NewServer(server.Broker(b))

	service := micro.NewService(
		micro.Server(s),
		micro.Broker(b),
	)
	h := &Example{}
	// Register a subscriber
	micro.RegisterSubscriber(
		"topic",
		service.Server(),
		h.Handler,
		server.SubscriberContext(brkrSub.Context),
		server.SubscriberQueue("queue.default"),
	)

	//service.Init()

	if err := service.Run(); err != nil {
		t.Fatal(err)
	}

}
