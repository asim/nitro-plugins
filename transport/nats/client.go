package nats

import (
	"encoding/json"
	"time"

	microtransport "github.com/micro/go-micro/transport"
	"github.com/nats-io/nats"
)

type client struct {
	conn *nats.Conn
	addr string
	id   string
	sub  *nats.Subscription
}

func (c *client) Send(m *microtransport.Message) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return c.conn.PublishRequest(c.addr, c.id, b)
}

func (c *client) Recv(m *microtransport.Message) error {
	rsp, err := c.sub.NextMsg(time.Second * 10)
	if err != nil {
		return err
	}

	var mr *microtransport.Message
	if err := json.Unmarshal(rsp.Data, &mr); err != nil {
		return err
	}

	*m = *mr
	return nil
}

func (c *client) Close() error {
	if err := c.Send(&microtransport.Message{Header: map[string]string{"Close": "true"}}); err != nil {
		return err
	}
	c.sub.Unsubscribe()
	c.conn.Close()
	return nil
}
