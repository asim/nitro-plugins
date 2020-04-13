// Package kafka provides a kafka broker using segmentio
package segmentio

import (
	"context"
	"errors"
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

	connected bool
	sync.RWMutex
	opts broker.Options
}

type subscriber struct {
	group *kafka.ConsumerGroup
	t     string
	opts  broker.SubscribeOptions
}

type publication struct {
	t   string
	err error
	m   *broker.Message
	ctx context.Context
	gen *kafka.Generation
	mp  map[string]map[int]int64 // for commit offsets
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
	return p.gen.CommitOffsets(p.mp)
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
	return s.group.Close()
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

	gcfg := kafka.ConsumerGroupConfig{
		ID:                    opt.Queue,
		WatchPartitionChanges: true,
		Brokers:               k.readerConfig.Brokers,
		Topics:                []string{topic},
	}
	if err := gcfg.Validate(); err != nil {
		return nil, err
	}

	group, err := kafka.NewConsumerGroup(gcfg)
	if err != nil {
		return nil, err
	}

	sub := &subscriber{group: group, opts: opt, t: topic}

	go func() {
		for {
			gen, err := group.Next(k.opts.Context)
			if err == kafka.ErrGroupClosed {
				return
			} else if err != nil {
				if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
					logger.Errorf("[kafka] subscribe error: %v", err)
				}
				return
			}
			assignments := gen.Assignments[topic]
			for _, assignment := range assignments {
				partition, offset := assignment.ID, assignment.Offset
				p := &publication{t: topic, ctx: k.opts.Context, gen: gen}
				p.mp = map[string]map[int]int64{p.t: {partition: offset}}

				gen.Start(func(ctx context.Context) {
					// create reader for this partition.
					cfg := k.readerConfig
					cfg.Topic = topic
					cfg.Partition = partition
					reader := kafka.NewReader(cfg)
					defer reader.Close()
					// seek to the last committed offset for this partition.
					reader.SetOffset(offset)
					for {
						msg, err := reader.ReadMessage(ctx)
						switch err {
						case kafka.ErrGenerationEnded:
							// generation has ended
							if logger.V(logger.DebugLevel, logger.DefaultLogger) {
								logger.Debug("[kafka] subscription closed")
							}
							return
						case nil:
							var m broker.Message
							eh := k.opts.ErrorHandler
							p.m = &m
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
								if err = p.Ack(); err != nil {
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
				})
			}
		}
	}()

	return sub, nil
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

	writerConfig := kafka.WriterConfig{CompressionCodec: nil}
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
		addrs:        cAddrs,
		opts:         options,
	}
}
