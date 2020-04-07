// Package kafka provides a kafka broker using segmentio
package segmentio

import (
	"context"
	"errors"
	"io"
	"sync"

	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/codec/json"
	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/go-micro/v2/logger"
	kafka "github.com/segmentio/kafka-go"
)

type kBroker struct {
	addrs []string

	readerConfig kafka.ReaderConfig
	writerConfig kafka.WriterConfig

	writers map[string]*kafka.Writer
	readers map[string]*kafka.Reader

	connected bool
	sync.RWMutex
	opts broker.Options
}

type subscriber struct {
	reader *kafka.Reader
	t      string
	opts   broker.SubscribeOptions
}

type publication struct {
	t      string
	err    error
	reader *kafka.Reader
	// deprecate broker.Message and use kafka.Message directly?
	m   *broker.Message
	ctx context.Context
	km  kafka.Message
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
	return p.reader.CommitMessages(p.ctx, p.km)
}

func (p *publication) Error() error {
	return p.err
}

func (s *subscriber) Options() broker.SubscribeOptions {
	return s.opts
}

func (s *subscriber) Topic() string {
	return s.t
}

func (s *subscriber) Unsubscribe() error {
	return s.reader.Close()
}

func (k *kBroker) Address() string {
	if len(k.addrs) > 0 {
		return k.addrs[0]
	}
	return "127.0.0.1:9092"
}

func (k *kBroker) Connect() error {
	k.RLock()
	if k.connected {
		k.RUnlock()
		return nil
	}
	k.RUnlock()

	kaddrs := make([]string, 0, len(k.addrs))
	for _, addr := range k.addrs {
		conn, err := kafka.DialContext(k.opts.Context, "tcp", addr)
		if err != nil {
			continue
		}
		if _, err = conn.Brokers(); err != nil {
			conn.Close()
			continue
		}
		kaddrs = append(kaddrs, addr)
		conn.Close()
	}

	if len(kaddrs) == 0 {
		return errors.New("no available brokers")
	}

	k.Lock()
	k.addrs = kaddrs
	k.connected = true
	k.Unlock()

	return nil
}

func (k *kBroker) Disconnect() error {
	k.RLock()
	if !k.connected {
		k.RUnlock()
		return nil
	}
	k.RUnlock()

	k.Lock()
	defer k.Unlock()
	for _, reader := range k.readers {
		if err := reader.Close(); err != nil {
			return err
		}
	}
	for _, writer := range k.writers {
		if err := writer.Close(); err != nil {
			return err
		}
	}

	k.connected = false
	return nil
}

func (k *kBroker) Init(opts ...broker.Option) error {
	for _, o := range opts {
		o(&k.opts)
	}
	var cAddrs []string
	for _, addr := range k.opts.Addrs {
		if len(addr) == 0 {
			continue
		}
		cAddrs = append(cAddrs, addr)
	}
	if len(cAddrs) == 0 {
		cAddrs = []string{"127.0.0.1:9092"}
	}
	k.addrs = cAddrs
	return nil
}

func (k *kBroker) Options() broker.Options {
	return k.opts
}

func (k *kBroker) Publish(topic string, msg *broker.Message, opts ...broker.PublishOption) error {
	buf, err := k.opts.Codec.Marshal(msg)
	if err != nil {
		return err
	}

	k.Lock()
	writer, ok := k.writers[topic]
	if !ok {
		cfg := k.writerConfig
		cfg.Topic = topic
		if err = cfg.Validate(); err != nil {
			k.Unlock()
			return err
		}
		writer = kafka.NewWriter(cfg)
		k.writers[topic] = writer
	}
	k.Unlock()

	return writer.WriteMessages(k.opts.Context, kafka.Message{Value: buf})
}

func (k *kBroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	opt := broker.SubscribeOptions{
		AutoAck: true,
		Queue:   uuid.New().String(),
	}
	for _, o := range opts {
		o(&opt)
	}

	k.Lock()
	reader, ok := k.readers[topic]
	if !ok {
		cfg := k.readerConfig
		cfg.Topic = topic
		cfg.GroupID = opt.Queue
		if err := cfg.Validate(); err != nil {
			k.Unlock()
			return nil, err
		}
		reader = kafka.NewReader(cfg)
		k.readers[topic] = reader
	}
	k.Unlock()
	go func() {
		for {
			select {
			case <-k.opts.Context.Done():
				return
			default:
				msg, err := reader.FetchMessage(k.opts.Context)
				if err != nil && err != io.EOF {
					if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
						logger.Errorf("[kafka] subscribe error: %v", err)
						return
					}
				} else if err == io.EOF {
					// reader closed
					return
				}

				var m broker.Message
				p := &publication{m: &m, t: msg.Topic, km: msg, ctx: k.opts.Context, reader: reader}
				eh := k.opts.ErrorHandler

				if err := k.opts.Codec.Unmarshal(msg.Value, &m); err != nil {
					p.err = err
					p.m.Body = msg.Value
					if eh != nil {
						eh(p)
					} else {
						if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
							logger.Errorf("[kafka]: failed to unmarshal: %v", err)
						}
					}
					continue
				}
				err = handler(p)
				if err == nil && opt.AutoAck {
					if err = reader.CommitMessages(k.opts.Context, msg); err != nil {
						logger.Errorf("[kafka]: unable to commit msg: %v", err)
					}
				} else if err != nil {
					p.err = err
					if eh != nil {
						eh(p)
					} else {
						if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
							logger.Errorf("[kafka]: subscriber error: %v", err)
						}
					}
				}
			}
		}
	}()

	return &subscriber{reader: reader, opts: opt, t: topic}, nil
}

func (k *kBroker) String() string {
	return "kafka"
}

func NewBroker(opts ...broker.Option) broker.Broker {
	options := broker.Options{
		// default to json codec
		Codec:   json.Marshaler{},
		Context: context.Background(),
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

	readerConfig := kafka.ReaderConfig{}
	if cfg, ok := options.Context.Value(readerConfigKey{}).(kafka.ReaderConfig); ok {
		readerConfig = cfg
	}
	if len(readerConfig.Brokers) == 0 {
		readerConfig.Brokers = cAddrs
	}
	readerConfig.WatchPartitionChanges = true

	writerConfig := kafka.WriterConfig{}
	if cfg, ok := options.Context.Value(writerConfigKey{}).(kafka.WriterConfig); ok {
		writerConfig = cfg
	}
	if len(writerConfig.Brokers) == 0 {
		writerConfig.Brokers = cAddrs
	}
	writerConfig.BatchSize = 1

	return &kBroker{
		readerConfig: readerConfig,
		writerConfig: writerConfig,
		writers:      make(map[string]*kafka.Writer),
		readers:      make(map[string]*kafka.Reader),
		addrs:        cAddrs,
		opts:         options,
	}
}
