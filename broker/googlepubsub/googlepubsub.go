// Package googlepubsub provides a Google cloud pubsub broker
package googlepubsub

import (
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/cmd"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
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
	ctx := context.Background()
	for {
		select {
		case <-s.exit:
			return
		default:
			iter, err := s.sub.Pull(ctx)
			if err != nil {
				time.Sleep(time.Second)
				continue
			}
			defer iter.Stop()

			// iterate over the messages
			for {
				pm, err := iter.Next()
				if err != nil {
					break
				}

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
		return s.sub.Delete(context.Background())
	}
	return nil
}

func (p *publication) Ack() error {
	p.pm.Done(true)
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
		ID:         "m-" + uuid.NewUUID().String(),
		Data:       msg.Body,
		Attributes: msg.Header,
	}

	_, err = t.Publish(ctx, m)
	return err
}

// Subscribe registers a subscription to the given topic against the google pubsub api
func (b *pubsubBroker) Subscribe(topic string, h broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	options := broker.SubscribeOptions{
		AutoAck: true,
		Queue:   "q-" + uuid.NewUUID().String(),
	}

	for _, o := range opts {
		o(&options)
	}

	ctx := context.Background()
	sub := b.client.Subscription(options.Queue)
	exists, err := sub.Exists(ctx)
	if err != nil {
		return nil, err
	}

	if !exists {
		tt := b.client.Topic(topic)
		subb, err := b.client.CreateSubscription(ctx, options.Queue, tt, time.Duration(0), nil)
		if err != nil {
			return nil, err
		}
		sub = subb
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
