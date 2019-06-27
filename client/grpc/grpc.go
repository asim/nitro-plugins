// Package grpc provides a gRPC client
// Deprecated: use `github.com/micro/go-micro/client/grpc` instead
package grpc

import (
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/client/grpc"
)

// Deprecated: use `github.com/micro/go-micro/client/grpc` instead
func NewClient(opts ...client.Option) client.Client {
	return grpc.NewClient(opts...)
}
