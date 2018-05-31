// Package nsq provides an NSQ broker
package nsq

import (
	"math/rand"
	"sync"
	"time"

	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/broker/codec/json"
	"github.com/micro/go-micro/cmd"
	"github.com/nsqio/go-nsq"
	"github.com/pborman/uuid"
)

type nsqBroker struct {
	addrs  []string
	opts   broker.Options
	config *nsq.Config

	sync.Mutex
	running bool
	p       []*nsq.Producer
	c       []*subscriber
}

type publication struct {
	topic string
	m     *broker.Message
	nm    *nsq.Message
	opts  broker.PublishOptions
}

type subscriber struct {
	topic string
	opts  broker.SubscribeOptions

	c *nsq.Consumer

	// handler so we can resubcribe
	h nsq.HandlerFunc
	// concurrency
	n int
}

var (
	DefaultConcurrentHandlers = 1
)

func init() {
	rand.Seed(time.Now().UnixNano())
	cmd.DefaultBrokers["nsq"] = NewBroker
}

func (n *nsqBroker) Init(opts ...broker.Option) error {
	for _, o := range opts {
		o(&n.opts)
	}
	return nil
}

func (n *nsqBroker) Options() broker.Options {
	return n.opts
}

func (n *nsqBroker) Address() string {
	return n.addrs[rand.Int()%len(n.addrs)]
}

func (n *nsqBroker) Connect() error {
	n.Lock()
	defer n.Unlock()

	if n.running {
		return nil
	}

	var producers []*nsq.Producer

	// create producers
	for _, addr := range n.addrs {
		p, err := nsq.NewProducer(addr, n.config)
		if err != nil {
			return err
		}

		producers = append(producers, p)
	}

	// create consumers
	for _, c := range n.c {
		channel := c.opts.Queue
		if len(channel) == 0 {
			channel = uuid.NewUUID().String()
		}

		cm, err := nsq.NewConsumer(c.topic, channel, n.config)
		if err != nil {
			return err
		}

		cm.AddConcurrentHandlers(c.h, c.n)

		c.c = cm

		err = c.c.ConnectToNSQDs(n.addrs)
		if err != nil {
			return err
		}
	}

	n.p = producers
	n.running = true
	return nil
}

func (n *nsqBroker) Disconnect() error {
	n.Lock()
	defer n.Unlock()

	if !n.running {
		return nil
	}

	// stop the producers
	for _, p := range n.p {
		p.Stop()
	}

	// stop the consumers
	for _, c := range n.c {
		c.c.Stop()

		// disconnect from all nsq brokers
		for _, addr := range n.addrs {
			c.c.DisconnectFromNSQD(addr)
		}
	}

	n.p = nil
	n.running = false
	return nil
}

func (n *nsqBroker) Publish(topic string, message *broker.Message, opts ...broker.PublishOption) error {
	p := n.p[rand.Int()%len(n.p)]

	b, err := n.opts.Codec.Marshal(message)
	if err != nil {
		return err
	}

	if ttl, ok := message.Header["ttl"]; ok {

		delay, err := time.ParseDuration(ttl)

		if err != nil {
			return p.Publish(topic, b)
		}

		if delay <= 0 {
			return p.Publish(topic, b)
		}

		return p.DeferredPublish(topic, delay, b)

	} else {
		return p.Publish(topic, b)
	}
}

func (n *nsqBroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	options := broker.SubscribeOptions{
		AutoAck: true,
	}

	for _, o := range opts {
		o(&options)
	}

	var concurrency int

	if options.Context != nil {
		var ok bool
		concurrency, ok = options.Context.Value(concurrentHandlerKey).(int)
		if !ok {
			concurrency = DefaultConcurrentHandlers
		}
	} else {
		concurrency = DefaultConcurrentHandlers

	}

	channel := options.Queue
	if len(channel) == 0 {
		channel = uuid.NewUUID().String()
	}

	c, err := nsq.NewConsumer(topic, channel, n.config)
	if err != nil {
		return nil, err
	}

	h := nsq.HandlerFunc(func(nm *nsq.Message) error {
		if !options.AutoAck {
			nm.DisableAutoResponse()
		}

		var m broker.Message

		if err := n.opts.Codec.Unmarshal(nm.Body, &m); err != nil {
			return err
		}

		return handler(&publication{
			topic: topic,
			m:     &m,
			nm:    nm,
		})

	})

	c.AddConcurrentHandlers(h, concurrency)

	err = c.ConnectToNSQDs(n.addrs)
	if err != nil {
		return nil, err
	}

	return &subscriber{
		topic: topic,
		c:     c,
		h:     h,
		n:     concurrency,
	}, nil
}

func (n *nsqBroker) String() string {
	return "nsq"
}

func (p *publication) Topic() string {
	return p.topic
}

func (p *publication) Message() *broker.Message {
	return p.m
}

func (p *publication) Ack() error {
	p.nm.Finish()
	return nil
}

func (s *subscriber) Options() broker.SubscribeOptions {
	return s.opts
}

func (s *subscriber) Topic() string {
	return s.topic
}

func (s *subscriber) Unsubscribe() error {
	s.c.Stop()
	return nil
}

func NewBroker(opts ...broker.Option) broker.Broker {
	options := broker.Options{
		// Default codec
		Codec: json.NewCodec(),
	}

	for _, o := range opts {
		o(&options)
	}

	var cAddrs []string

	for _, addr := range options.Addrs {
		if len(addr) > 0 {
			cAddrs = append(cAddrs, addr)
		}
	}

	if len(cAddrs) == 0 {
		cAddrs = []string{"127.0.0.1:4150"}
	}

	return &nsqBroker{
		addrs:  cAddrs,
		opts:   options,
		config: nsq.NewConfig(),
	}
}
