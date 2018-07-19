package rabbitmq

//
// All credit to Mondo
//

import (
	"crypto/tls"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/streadway/amqp"
	"log"
)

var (
	DefaultExchange  = "micro"
	DefaultRabbitURL = "amqp://guest:guest@127.0.0.1:5672"

	dial    = amqp.Dial
	dialTLS = amqp.DialTLS
)

type rabbitMQConn struct {
	Connection      *amqp.Connection
	Channel         *rabbitMQChannel
	ExchangeChannel *rabbitMQChannel
	exchange        string
	url             string

	sync.Mutex
	connected bool
	close     chan bool

	waitConnection chan struct{}
}

func newRabbitMQConn(exchange string, urls []string) *rabbitMQConn {
	var url string

	if len(urls) > 0 && regexp.MustCompile("^amqp(s)?://.*").MatchString(urls[0]) {
		url = urls[0]
	} else {
		url = DefaultRabbitURL
	}

	if len(exchange) == 0 {
		exchange = DefaultExchange
	}

	ret := &rabbitMQConn{
		exchange:       exchange,
		url:            url,
		close:          make(chan bool),
		waitConnection: make(chan struct{}),
	}
	// its bad case of nil == waitConnection, so close it at start
	close(ret.waitConnection)
	return ret
}

func (r *rabbitMQConn) connect(secure bool, config *tls.Config) error {
	// try connect
	if err := r.tryConnect(secure, config); err != nil {
		return err
	}

	// connected
	r.Lock()
	r.connected = true
	r.Unlock()

	// create reconnect loop
	go r.reconnect(secure, config)
	return nil
}

func (r *rabbitMQConn) reconnect(secure bool, config *tls.Config) {
	// skip first connect
	var connect bool

	for {
		if connect {
			// try reconnect
			if err := r.tryConnect(secure, config); err != nil {
				time.Sleep(1 * time.Second)
				continue
			}

			// connected
			r.Lock()
			r.connected = true
			r.Unlock()
			//unblock resubscribe cycle - close channel
			//at this point channel is created and unclosed - close it without any additional checks
			close(r.waitConnection)
		}

		connect = true
		notifyClose := make(chan *amqp.Error)
		r.Connection.NotifyClose(notifyClose)

		// block until closed
		select {
		case <-notifyClose:
			// block all resubscribe attempt - they are useless because there is no connection to rabbitmq
			// create channel 'waitConnection' (at this point channel is nil or closed, create it without unnecessary checks)
			r.Lock()
			r.connected = false
			r.waitConnection = make(chan struct{})
			r.Unlock()
		case <-r.close:
			return
		}
	}
}

func (r *rabbitMQConn) Connect(secure bool, config *tls.Config) error {
	r.Lock()

	// already connected
	if r.connected {
		r.Unlock()
		return nil
	}

	// check it was closed
	select {
	case <-r.close:
		r.close = make(chan bool)
	default:
		// no op
		// new conn
	}

	r.Unlock()

	return r.connect(secure, config)
}

func (r *rabbitMQConn) Close() error {
	r.Lock()
	defer r.Unlock()

	select {
	case <-r.close:
		return nil
	default:
		close(r.close)
		r.connected = false
	}

	return r.Connection.Close()
}

func (r *rabbitMQConn) tryConnect(secure bool, config *tls.Config) error {
	var err error

	if secure || config != nil || strings.HasPrefix(r.url, "amqps://") {
		if config == nil {
			config = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		url := strings.Replace(r.url, "amqp://", "amqps://", 1)
		r.Connection, err = dialTLS(url, config)
	} else {
		r.Connection, err = dial(r.url)
	}

	if err != nil {
		return err
	}

	if r.Channel, err = newRabbitChannel(r.Connection); err != nil {
		return err
	}

	r.Channel.DeclareDurableExchange(r.exchange)
	r.ExchangeChannel, err = newRabbitChannel(r.Connection)

	return err
}

func (r *rabbitMQConn) Consume(queue, key string, headers amqp.Table, autoAck, durableQueue bool) (*rabbitMQChannel, <-chan amqp.Delivery, error) {
	consumerChannel, err := newRabbitChannel(r.Connection)
	if err != nil {
		return nil, nil, err
	}

	if durableQueue {
		err = consumerChannel.DeclareDurableQueue(queue)
	} else {
		err = consumerChannel.DeclareQueue(queue)
	}

	if err != nil {
		return nil, nil, err
	}

	if err=consumerChannel.channel.Qos(1000,0,false); err!= nil {
		log.Println(err)
	}

	deliveries, err := consumerChannel.ConsumeQueue(queue, autoAck)
	if err != nil {
		return nil, nil, err
	}

	err = consumerChannel.BindQueue(queue, key, r.exchange, headers)
	if err != nil {
		return nil, nil, err
	}

	return consumerChannel, deliveries, nil
}

func (r *rabbitMQConn) Publish(exchange, key string, msg amqp.Publishing) error {
	return r.ExchangeChannel.Publish(exchange, key, msg)
}
