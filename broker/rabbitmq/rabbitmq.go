// Package rabbitmq provides a RabbitMQ broker
package rabbitmq

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/cmd"
	"github.com/streadway/amqp"
)

type rbroker struct {
	conn  *rabbitMQConn
	addrs []string
	opts  broker.Options
	mtx   sync.Mutex
	wg    sync.WaitGroup
}

type subscriber struct {
	mtx          sync.Mutex
	mayRun       bool
	opts         broker.SubscribeOptions
	topic        string
	ch           *rabbitMQChannel
	durableQueue bool
	r            *rbroker
	fn           func(msg amqp.Delivery)
	headers      map[string]interface{}
}

type publication struct {
	d amqp.Delivery
	m *broker.Message
	t string
}

func init() {
	cmd.DefaultBrokers["rabbitmq"] = NewBroker
}

func (p *publication) Ack() error {
	return p.d.Ack(false)
}

func (p *publication) Topic() string {
	return p.t
}

func (p *publication) Message() *broker.Message {
	return p.m
}

func (s *subscriber) Options() broker.SubscribeOptions {
	return s.opts
}

func (s *subscriber) Topic() string {
	return s.topic
}

func (s *subscriber) Unsubscribe() error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.mayRun = false
	if s.ch != nil {
		return s.ch.Close()
	}
	return nil
}

func (s *subscriber) resubscribe() {
	minResubscribeDelay := 100 * time.Millisecond
	maxResubscribeDelay := 30 * time.Second
	expFactor := time.Duration(2)
	reSubscribeDelay := minResubscribeDelay
	//loop until unsubscribe
	for {
		s.mtx.Lock()
		mayRun := s.mayRun
		s.mtx.Unlock()
		if !mayRun {
			// we are unsubscribed, showdown routine
			return
		}

		select {
		//check shutdown case
		case <-s.r.conn.close:
			//yep, its shutdown case
			return
		//wait until we reconect to rabbit
		case <-s.r.conn.waitConnection:
		}

		// it may crash (panic) in case of Consume without connection, so recheck it
		s.r.mtx.Lock()
		if !s.r.conn.connected {
			s.r.mtx.Unlock()
			continue
		}

		s.r.conn.Channel.channel.Qos(1000,100*1024,true)
		ch, sub, err := s.r.conn.Consume(
			s.opts.Queue,
			s.topic,
			s.headers,
			s.opts.AutoAck,
			s.durableQueue,
		)

		s.r.mtx.Unlock()
		switch err {
		case nil:
			reSubscribeDelay = minResubscribeDelay
			s.mtx.Lock()
			s.ch = ch
			s.mtx.Unlock()
		default:
			if reSubscribeDelay > maxResubscribeDelay {
				reSubscribeDelay = maxResubscribeDelay
			}
			time.Sleep(reSubscribeDelay)
			reSubscribeDelay *= expFactor
			continue
		}
		for d := range sub {
			s.r.wg.Add(1)
			go func(d amqp.Delivery) {
				s.fn(d)
				s.r.wg.Done()
			}(d)
		}
	}
}

func (r *rbroker) Publish(topic string, msg *broker.Message, opts ...broker.PublishOption) error {
	opt := broker.PublishOptions{}

	for _, o := range opts {
		o(&opt)
	}

	//persistent := uint8(2)
	//if opt.Context != nil && opt.Context.Value(persistentDelivery{}).(bool) {
	//	persistent = 2
	//}

	m := amqp.Publishing{
		Body:    msg.Body,
		Headers: amqp.Table{},
		//DeliveryMode:persistent,
	}

	for k, v := range msg.Header {
		m.Headers[k] = v
	}

	if r.conn == nil {
		return errors.New("connection is nil")
	}

	return r.conn.Publish(r.conn.exchange, topic, m)
}

func (r *rbroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	opt := broker.SubscribeOptions{
		AutoAck: true,
	}

	for _, o := range opts {
		o(&opt)
	}

	durableQueue := false
	if opt.Context != nil {
		durableQueue, _ = opt.Context.Value(durableQueueKey{}).(bool)
	}

	var headers map[string]interface{}
	if opt.Context != nil {
		if h, ok := opt.Context.Value(headersKey{}).(map[string]interface{}); ok {
			headers = h
		}
	}

	if r.conn == nil {
		return nil, errors.New("connection is nil")
	}

	fn := func(msg amqp.Delivery) {
		header := make(map[string]string)
		for k, v := range msg.Headers {
			header[k], _ = v.(string)
		}
		m := &broker.Message{
			Header: header,
			Body:   msg.Body,
		}
		handler(&publication{d: msg, m: m, t: msg.RoutingKey})
	}

	sret := &subscriber{topic: topic, opts: opt, mayRun: true, r: r,
		durableQueue: durableQueue, fn: fn, headers: headers}

	go sret.resubscribe()

	return sret, nil
}

func (r *rbroker) Options() broker.Options {
	return r.opts
}

func (r *rbroker) String() string {
	return "rabbitmq"
}

func (r *rbroker) Address() string {
	if len(r.addrs) > 0 {
		return r.addrs[0]
	}
	return ""
}

func (r *rbroker) Init(opts ...broker.Option) error {
	for _, o := range opts {
		o(&r.opts)
	}
	return nil
}

func (r *rbroker) Connect() error {
	if r.conn == nil {
		r.conn = newRabbitMQConn(r.getExchange(), r.opts.Addrs)
	}
	return r.conn.Connect(r.opts.Secure, r.opts.TLSConfig)
}

func (r *rbroker) Disconnect() error {
	if r.conn == nil {
		return errors.New("connection is nil")
	}
	ret := r.conn.Close()
	r.wg.Wait() // wait all goroutines
	return ret
}

func NewBroker(opts ...broker.Option) broker.Broker {
	options := broker.Options{
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return &rbroker{
		addrs: options.Addrs,
		opts:  options,
	}
}

func (r *rbroker) getExchange() string {
	if e, ok := r.opts.Context.Value(exchangeKey{}).(string); ok {
		return e
	}
	return DefaultExchange
}
