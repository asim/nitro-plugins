// Package grpc provides a grpc server
// Deprecated: use `github.com/micro/go-micro/server/grpc` instead
package grpc

import (
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/server/grpc"
)

// We use this to wrap any debug handlers so we preserve the signature Debug.{Method}
// Deprecated: use `github.com/micro/go-micro/server/grpc` instead
type Debug = grpc.Debug

var (
	// DefaultMaxMsgSize define maximum message size that server can send
	// or receive.  Default value is 4MB.
	// Deprecated: use `github.com/micro/go-micro/server/grpc` instead
	DefaultMaxMsgSize = grpc.DefaultMaxMsgSize
)

// Deprecated: use `github.com/micro/go-micro/server/grpc` instead
func NewServer(opts ...server.Option) server.Server {
	return grpc.NewServer(opts...)
}
