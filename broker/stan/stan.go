// Package stan provides a NATS Streaming broker
package stan

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/codec/json"
	stan "github.com/nats-io/go-nats-streaming"
)

type stanBroker struct {
	sync.RWMutex
	addrs []string
	conn  stan.Conn
	opts  broker.Options
	nopts stan.Options
}

type subscriber struct {
	t    string
	s    stan.Subscription
	dq   bool
	opts broker.SubscribeOptions
}

type publication struct {
	t   string
	msg *stan.Msg
	m   *broker.Message
}

func init() {
	cmd.DefaultBrokers["stan"] = NewBroker
}

func (n *publication) Topic() string {
	return n.t
}

func (n *publication) Message() *broker.Message {
	return n.m
}

func (n *publication) Ack() error {
	return n.msg.Ack()
}

func (n *subscriber) Options() broker.SubscribeOptions {
	return n.opts
}

func (n *subscriber) Topic() string {
	return n.t
}

func (n *subscriber) Unsubscribe() error {
	if n.s == nil {
		return nil
	}
	// go-micro server Unsubscribe can't handle durable queues, so close as stan suggested
	// from nats streaming readme:
	// When a client disconnects, the streaming server is not notified, hence the importance of calling Close()
	if !n.dq {
		err := n.s.Unsubscribe()
		if err != nil {
			return err
		}
	}
	return n.Close()
}

func (n *subscriber) Close() error {
	if n.s != nil {
		return n.s.Close()
	}
	return nil
}

func (n *stanBroker) Address() string {
	// stan does not support connected server info
	if len(n.addrs) > 0 {
		return n.addrs[0]
	}

	return ""
}

func setAddrs(addrs []string) []string {
	var cAddrs []string
	for _, addr := range addrs {
		if len(addr) == 0 {
			continue
		}
		if !strings.HasPrefix(addr, "nats://") {
			addr = "nats://" + addr
		}
		cAddrs = append(cAddrs, addr)
	}
	if len(cAddrs) == 0 {
		cAddrs = []string{stan.DefaultNatsURL}
	}
	return cAddrs
}

func (n *stanBroker) Connect() error {
	n.RLock()
	if n.conn != nil {
		n.RUnlock()
		return nil
	}
	n.RUnlock()

	opts := n.nopts
	opts.NatsURL = strings.Join(n.addrs, ",")

	clusterID, ok := n.opts.Context.Value(clusterIDKey{}).(string)
	if !ok || len(clusterID) == 0 {
		return errors.New("must specify ClusterID Option")
	}

	clientID := uuid.New().String()

	nopts := []stan.Option{
		stan.NatsURL(opts.NatsURL),
		stan.NatsConn(opts.NatsConn),
		stan.ConnectWait(opts.ConnectTimeout),
		stan.PubAckWait(opts.AckTimeout),
		stan.MaxPubAcksInflight(opts.MaxPubAcksInflight),
		stan.Pings(opts.PingIterval, opts.PingMaxOut),
		stan.SetConnectionLostHandler(opts.ConnectionLostCB),
	}

	c, err := stan.Connect(clusterID, clientID, nopts...)
	if err != nil {
		return err
	}
	n.Lock()
	n.conn = c
	n.Unlock()
	return nil
}

func (n *stanBroker) Disconnect() error {
	n.RLock()
	n.conn.Close()
	n.RUnlock()
	return nil
}

func (n *stanBroker) Init(opts ...broker.Option) error {
	for _, o := range opts {
		o(&n.opts)
	}
	n.addrs = setAddrs(n.opts.Addrs)
	return nil
}

func (n *stanBroker) Options() broker.Options {
	return n.opts
}

func (n *stanBroker) Publish(topic string, msg *broker.Message, opts ...broker.PublishOption) error {
	b, err := n.opts.Codec.Marshal(msg)
	if err != nil {
		return err
	}
	n.RLock()
	defer n.RUnlock()
	return n.conn.Publish(topic, b)
}

func (n *stanBroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	if n.conn == nil {
		return nil, errors.New("not connected")
	}
	var successAutoAck bool

	opt := broker.SubscribeOptions{
		AutoAck: true,
	}

	for _, o := range opts {
		o(&opt)
	}

	// Make sure context is setup
	if opt.Context == nil {
		opt.Context = context.Background()
	}

	ctx := opt.Context
	if subscribeContext, ok := ctx.Value(subscribeContextKey{}).(context.Context); ok && subscribeContext != nil {
		ctx = subscribeContext
	}

	var stanOpts []stan.SubscriptionOption
	if !opt.AutoAck {
		stanOpts = append(stanOpts, stan.SetManualAckMode())
	}

	if subOpts, ok := ctx.Value(subscribeOptionKey{}).([]stan.SubscriptionOption); ok && len(subOpts) > 0 {
		stanOpts = append(stanOpts, subOpts...)
	}

	if bval, ok := ctx.Value(successAutoAckKey{}).(bool); ok && bval {
		stanOpts = append(stanOpts, stan.SetManualAckMode())
		successAutoAck = true
	}

	bopts := stan.DefaultSubscriptionOptions
	for _, bopt := range stanOpts {
		if err := bopt(&bopts); err != nil {
			return nil, err
		}
	}

	opt.AutoAck = !bopts.ManualAcks

	fn := func(msg *stan.Msg) {
		var m broker.Message
		var err error
		if err = n.opts.Codec.Unmarshal(msg.Data, &m); err != nil {
			return
		}
		err = handler(&publication{m: &m, msg: msg, t: msg.Subject})
		if err == nil && successAutoAck && !opt.AutoAck {
			msg.Ack()
		} else if err != nil && opt.AutoAck {
			msg.Ack()
		}
	}

	var sub stan.Subscription
	var err error

	n.RLock()
	if len(opt.Queue) > 0 {
		sub, err = n.conn.QueueSubscribe(topic, opt.Queue, fn, stanOpts...)
	} else {
		sub, err = n.conn.Subscribe(topic, fn, stanOpts...)
	}
	n.RUnlock()
	if err != nil {
		return nil, err
	}
	return &subscriber{dq: len(bopts.DurableName) > 0, s: sub, opts: opt, t: topic}, nil
}

func (n *stanBroker) String() string {
	return "stan"
}

func NewBroker(opts ...broker.Option) broker.Broker {
	options := broker.Options{
		// Default codec
		Codec:   json.Marshaler{},
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	stanOpts := stan.DefaultOptions
	if n, ok := options.Context.Value(optionsKey{}).(stan.Options); ok {
		stanOpts = n
	}

	if len(options.Addrs) == 0 {
		options.Addrs = strings.Split(stanOpts.NatsURL, ",")
	}

	nb := &stanBroker{
		opts:  options,
		nopts: stanOpts,
		addrs: setAddrs(options.Addrs),
	}

	return nb
}
