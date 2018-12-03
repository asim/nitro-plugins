// Package googlepubsub provides a Google cloud pubsub broker
package googlepubsub

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/cmd"
	"google.golang.org/api/option"
)

type pubsubBroker struct {
	client  *pubsub.Client
	options broker.Options
}

// A pubsub subscriber that manages handling of messages
type subscriber struct {
	options broker.SubscribeOptions
	topic   string
	exit    chan bool
	sub     *pubsub.Subscription
}

// A single publication received by a handler
type publication struct {
	pm    *pubsub.Message
	m     *broker.Message
	topic string
}

func init() {
	cmd.DefaultBrokers["googlepubsub"] = NewBroker
}

func (s *subscriber) run(hdlr broker.Handler) {
	if s.options.Context != nil {
		if max, ok := s.options.Context.Value(maxOutstandingMessagesKey{}).(int); ok {
			s.sub.ReceiveSettings.MaxOutstandingMessages = max
		}
		if max, ok := s.options.Context.Value(maxExtensionKey{}).(time.Duration); ok {
			s.sub.ReceiveSettings.MaxExtension = max
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	for {
		select {
		case <-s.exit:
			cancel()
			return
		default:
			if err := s.sub.Receive(ctx, func(ctx context.Context, pm *pubsub.Message) {
				// create broker message
				m := &broker.Message{
					Header: pm.Attributes,
					Body:   pm.Data,
				}

				// create publication
				p := &publication{
					pm:    pm,
					m:     m,
					topic: s.topic,
				}

				// If the error is nil lets check if we should auto ack
				if err := hdlr(p); err == nil {
					// auto ack?
					if s.options.AutoAck {
						p.Ack()
					}
				}
			}); err != nil {
				time.Sleep(time.Second)
				continue
			}
		}
	}
}

func (s *subscriber) Options() broker.SubscribeOptions {
	return s.options
}

func (s *subscriber) Topic() string {
	return s.topic
}

func (s *subscriber) Unsubscribe() error {
	select {
	case <-s.exit:
		return nil
	default:
		close(s.exit)
		if deleteSubscription, ok := s.options.Context.Value(deleteSubscription{}).(bool); !ok || deleteSubscription {
			return s.sub.Delete(context.Background())
		}
		return nil
	}
}

func (p *publication) Ack() error {
	p.pm.Ack()
	return nil
}

func (p *publication) Topic() string {
	return p.topic
}

func (p *publication) Message() *broker.Message {
	return p.m
}

func (b *pubsubBroker) Address() string {
	return ""
}

func (b *pubsubBroker) Connect() error {
	return nil
}

func (b *pubsubBroker) Disconnect() error {
	return b.client.Close()
}

// Init not currently implemented
func (b *pubsubBroker) Init(opts ...broker.Option) error {
	return nil
}

func (b *pubsubBroker) Options() broker.Options {
	return b.options
}

// Publish checks if the topic exists and then publishes via google pubsub
func (b *pubsubBroker) Publish(topic string, msg *broker.Message, opts ...broker.PublishOption) error {
	t := b.client.Topic(topic)
	ctx := context.Background()

	exists, err := t.Exists(ctx)
	if err != nil {
		return err
	}

	if !exists {
		tt, err := b.client.CreateTopic(ctx, topic)
		if err != nil {
			return err
		}
		t = tt
	}

	m := &pubsub.Message{
		ID:         "m-" + uuid.New().String(),
		Data:       msg.Body,
		Attributes: msg.Header,
	}

	pr := t.Publish(ctx, m)
	_, err = pr.Get(ctx)
	return err
}

// Subscribe registers a subscription to the given topic against the google pubsub api
func (b *pubsubBroker) Subscribe(topic string, h broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	queueName := "q-" + uuid.New().String()
	if subscriptionName, ok := b.options.Context.Value(subscriptionName{}).(string); ok {
		queueName = subscriptionName
	}

	options := broker.SubscribeOptions{
		AutoAck: true,
		Queue:   queueName,
		Context: b.options.Context,
	}

	for _, o := range opts {
		o(&options)
	}

	ctx := context.Background()
	sub := b.client.Subscription(options.Queue)

	if createSubscription, ok := b.options.Context.Value(createSubscription{}).(bool); !ok || createSubscription {
		exists, err := sub.Exists(ctx)
		if err != nil {
			return nil, err
		}

		if !exists {
			tt := b.client.Topic(topic)
			subb, err := b.client.CreateSubscription(ctx, options.Queue, pubsub.SubscriptionConfig{
				Topic:       tt,
				AckDeadline: time.Duration(0),
			})
			if err != nil {
				return nil, err
			}
			sub = subb
		}
	}

	subscriber := &subscriber{
		options: options,
		topic:   topic,
		exit:    make(chan bool),
		sub:     sub,
	}

	go subscriber.run(h)

	return subscriber, nil
}

func (b *pubsubBroker) String() string {
	return "googlepubsub"
}

// NewBroker creates a new google pubsub broker
func NewBroker(opts ...broker.Option) broker.Broker {
	options := broker.Options{
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	// retrieve project id
	prjID, _ := options.Context.Value(projectIDKey{}).(string)
	// retrieve client opts
	cOpts, _ := options.Context.Value(clientOptionKey{}).([]option.ClientOption)

	// create pubsub client
	c, err := pubsub.NewClient(context.Background(), prjID, cOpts...)
	if err != nil {
		panic(err.Error())
	}

	return &pubsubBroker{
		client:  c,
		options: options,
	}
}
