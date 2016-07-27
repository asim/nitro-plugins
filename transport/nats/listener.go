package nats

import (
	"sync"
	"time"

	microtransport "github.com/micro/go-micro/transport"
	"github.com/nats-io/nats"
)

type listener struct {
	conn *nats.Conn
	addr string
	exit chan bool

	sync.RWMutex
	so map[string]*socket
}

func (l *listener) Addr() string {
	return l.addr
}

func (l *listener) Close() error {
	l.exit <- true
	l.conn.Close()
	return nil
}

func (l *listener) Accept(fn func(microtransport.Socket)) error {
	s, err := l.conn.SubscribeSync(l.addr)
	if err != nil {
		return err
	}

	var lerr error

	go func() {
		<-l.exit
		lerr = s.Unsubscribe()
	}()

	for {
		m, err := s.NextMsg(time.Minute)
		if err != nil {
			if err == nats.ErrTimeout {
				continue
			}
			return err
		}

		l.RLock()
		sock, ok := l.so[m.Reply]
		l.RUnlock()

		if !ok {
			var once sync.Once
			sock = &socket{
				conn:  l.conn,
				once:  once,
				m:     m,
				r:     make(chan *nats.Msg, 1),
				close: make(chan bool),
			}
			l.Lock()
			l.so[m.Reply] = sock
			l.Unlock()

			go func() {
				// TODO: think of a better error response strategy
				defer func() {
					if r := recover(); r != nil {
						sock.Close()
					}
				}()
				fn(sock)
			}()

			go func() {
				<-sock.close
				l.Lock()
				delete(l.so, sock.m.Reply)
				l.Unlock()
			}()
		}

		select {
		case <-sock.close:
			continue
		default:
		}

		sock.Lock()
		sock.bl = append(sock.bl, m)
		select {
		case sock.r <- sock.bl[0]:
			sock.bl = sock.bl[1:]
		default:
		}
		sock.Unlock()

	}
	return lerr
}
