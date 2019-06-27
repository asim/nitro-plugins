// Package grpc provides a grpc server
package grpc

import (
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/server/grpc"
)

// We use this to wrap any debug handlers so we preserve the signature Debug.{Method}
type Debug = grpc.Debug

var (
	// DefaultMaxMsgSize define maximum message size that server can send
	// or receive.  Default value is 4MB.
	DefaultMaxMsgSize = grpc.DefaultMaxMsgSize
)

func NewServer(opts ...server.Option) server.Server {
	return grpc.NewServer(opts...)
}
