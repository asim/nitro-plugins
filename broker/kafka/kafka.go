// Package kafka provides a kafka broker using sarama cluster
package kafka

import (
	"log"

	"github.com/Shopify/sarama"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/broker/codec/json"
	"github.com/micro/go-micro/cmd"
	"github.com/pborman/uuid"
	sc "gopkg.in/bsm/sarama-cluster.v2"
)

type kBroker struct {
	addrs []string

	c  sarama.Client
	p  sarama.SyncProducer
	sc *sc.Client

	opts broker.Options
}

type subscriber struct {
	s    *sc.Consumer
	t    string
	opts broker.SubscribeOptions
}

type publication struct {
	t  string
	c  *sc.Consumer
	km *sarama.ConsumerMessage
	m  *broker.Message
}

func init() {
	cmd.DefaultBrokers["kafka"] = NewBroker
}

func (p *publication) Topic() string {
	return p.t
}

func (p *publication) Message() *broker.Message {
	return p.m
}

func (p *publication) Ack() error {
	p.c.MarkOffset(p.km, "")
	return nil
}

func (s *subscriber) Options() broker.SubscribeOptions {
	return s.opts
}

func (s *subscriber) Topic() string {
	return s.t
}

func (s *subscriber) Unsubscribe() error {
	return s.s.Close()
}

func (k *kBroker) Address() string {
	if len(k.addrs) > 0 {
		return k.addrs[0]
	}
	return "127.0.0.1:9092"
}

func (k *kBroker) Connect() error {
	if k.c != nil {
		return nil
	}

	c, err := sarama.NewClient(k.addrs, sarama.NewConfig())
	if err != nil {
		return err
	}

	k.c = c

	p, err := sarama.NewSyncProducerFromClient(c)
	if err != nil {
		return err
	}

	k.p = p

	config := sc.NewConfig()
	// TODO: make configurable offset as SubscriberOption
	config.Config.Consumer.Offsets.Initial = sarama.OffsetNewest

	cs, err := sc.NewClient(k.addrs, config)
	if err != nil {
		return err
	}

	k.sc = cs
	// TODO: TLS
	/*
		opts.Secure = k.opts.Secure
		opts.TLSConfig = k.opts.TLSConfig

		// secure might not be set
		if k.opts.TLSConfig != nil {
			opts.Secure = true
		}
	*/
	return nil
}

func (k *kBroker) Disconnect() error {
	k.sc.Close()
	k.p.Close()
	return k.c.Close()
}

func (k *kBroker) Init(opts ...broker.Option) error {
	for _, o := range opts {
		o(&k.opts)
	}
	return nil
}

func (k *kBroker) Options() broker.Options {
	return k.opts
}

func (k *kBroker) Publish(topic string, msg *broker.Message, opts ...broker.PublishOption) error {
	b, err := k.opts.Codec.Marshal(msg)
	if err != nil {
		return err
	}
	_, _, err = k.p.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(b),
	})
	return err
}

func (k *kBroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	opt := broker.SubscribeOptions{
		AutoAck: true,
		Queue:   uuid.NewUUID().String(),
	}

	for _, o := range opts {
		o(&opt)
	}

	c, err := sc.NewConsumerFromClient(k.sc, opt.Queue, []string{topic})
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case err := <-c.Errors():
				log.Println("consumer error:", err)
			case sm := <-c.Messages():
				var m broker.Message
				if err := k.opts.Codec.Unmarshal(sm.Value, &m); err != nil {
					continue
				}
				if err := handler(&publication{
					m:  &m,
					t:  sm.Topic,
					c:  c,
					km: sm,
				}); err == nil && opt.AutoAck {
					c.MarkOffset(sm, "")
				}
			}
		}
	}()

	return &subscriber{s: c, opts: opt}, nil
}

func (k *kBroker) String() string {
	return "kafka"
}

func NewBroker(opts ...broker.Option) broker.Broker {
	options := broker.Options{
		// default to json codec
		Codec: json.NewCodec(),
	}

	for _, o := range opts {
		o(&options)
	}

	var cAddrs []string
	for _, addr := range options.Addrs {
		if len(addr) == 0 {
			continue
		}
		cAddrs = append(cAddrs, addr)
	}
	if len(cAddrs) == 0 {
		cAddrs = []string{"127.0.0.1:9092"}
	}

	return &kBroker{
		addrs: cAddrs,
	}
}
