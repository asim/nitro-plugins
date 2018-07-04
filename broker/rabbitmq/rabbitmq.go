// Package rabbitmq provides a RabbitMQ broker
package rabbitmq

import (
	"context"
	"errors"

	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/cmd"
	"github.com/streadway/amqp"
)

type rbroker struct {
	conn  *rabbitMQConn
	addrs []string
	opts  broker.Options
}

type subscriber struct {
	opts  broker.SubscribeOptions
	topic string
	ch    *rabbitMQChannel
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
	return s.ch.Close()
}

func (r *rbroker) Publish(topic string, msg *broker.Message, opts ...broker.PublishOption) error {
	m := amqp.Publishing{
		Body:    msg.Body,
		Headers: amqp.Table{},
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

	ch, sub, err := r.conn.Consume(
		opt.Queue,
		topic,
		headers,
		opt.AutoAck,
		durableQueue,
	)
	if err != nil {
		return nil, err
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

	go func() {
		for d := range sub {
			fn(d)
		}
	}()

	return &subscriber{ch: ch, topic: topic, opts: opt}, nil
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
	return r.conn.Close()
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
