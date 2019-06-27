// Package grpc provides a grpc transport
package grpc

import (
	"github.com/micro/go-micro/transport"
	"github.com/micro/go-micro/transport/grpc"
)

func NewTransport(opts ...transport.Option) transport.Transport {
	return grpc.NewTransport(opts...)
}
