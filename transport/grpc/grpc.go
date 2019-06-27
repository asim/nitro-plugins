// Package grpc provides a grpc transport
// Deprecated: use `github.com/micro/go-micro/transport/grpc` instead
package grpc

import (
	"github.com/micro/go-micro/transport"
	"github.com/micro/go-micro/transport/grpc"
)

// Deprecated: use `github.com/micro/go-micro/transport/grpc` instead
func NewTransport(opts ...transport.Option) transport.Transport {
	return grpc.NewTransport(opts...)
}
