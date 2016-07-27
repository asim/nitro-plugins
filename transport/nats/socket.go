package nats

import (
	"encoding/json"
	"errors"
	"io"
	"sync"

	microtransport "github.com/micro/go-micro/transport"
	"github.com/nats-io/nats"
)

type socket struct {
	conn *nats.Conn
	m    *nats.Msg
	r    chan *nats.Msg

	once  sync.Once
	close chan bool

	sync.Mutex
	bl []*nats.Msg
}

func (s *socket) Recv(m *microtransport.Message) error {
	if m == nil {
		return errors.New("message passed in is nil")
	}

	r, ok := <-s.r
	if !ok {
		return io.EOF
	}
	s.Lock()
	if len(s.bl) > 0 {
		select {
		case s.r <- s.bl[0]:
			s.bl = s.bl[1:]
		default:
		}
	}
	s.Unlock()

	if err := json.Unmarshal(r.Data, &m); err != nil {
		return err
	}
	return nil
}

func (s *socket) Send(m *microtransport.Message) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return s.conn.Publish(s.m.Reply, b)
}

func (s *socket) Close() error {
	s.once.Do(func() {
		close(s.close)
	})
	return nil
}
