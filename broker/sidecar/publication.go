package sidecar

import (
	"github.com/micro/go-micro/broker"
)

type publication struct {
	topic   string
	message *broker.Message
}

func (p *publication) Topic() string {
	return p.topic
}

func (p *publication) Message() *broker.Message {
	return p.message
}

func (p *publication) Ack() error {
	return nil
}
