// Package grpc provides a gRPC client
package grpc

import (
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/client/grpc"
)

func NewClient(opts ...client.Option) client.Client {
	return grpc.NewClient(opts...)
}
