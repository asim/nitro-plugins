// Package grpc provides a gRPC options
package grpc

import (
	"crypto/tls"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/client/grpc"
	"google.golang.org/grpc/encoding"
)

var (
	// DefaultMaxRecvMsgSize maximum message that client can receive
	// (4 MB).
	// Deprecated: use `github.com/micro/go-micro/client/grpc` instead
	DefaultMaxRecvMsgSize = 1024 * 1024 * 4

	// DefaultMaxSendMsgSize maximum message that client can send
	// (4 MB).
	// Deprecated: use `github.com/micro/go-micro/client/grpc` instead
	DefaultMaxSendMsgSize = 1024 * 1024 * 4
)

// gRPC Codec to be used to encode/decode requests for a given content type
// Deprecated: use `github.com/micro/go-micro/client/grpc` instead
func Codec(contentType string, c encoding.Codec) client.Option {
	return grpc.Codec(contentType, c)
}

// AuthTLS should be used to setup a secure authentication using TLS
// Deprecated: use `github.com/micro/go-micro/client/grpc` instead
func AuthTLS(t *tls.Config) client.Option {
	return grpc.AuthTLS(t)
}

//
// MaxRecvMsgSize set the maximum size of message that client can receive.
//
// Deprecated: use `github.com/micro/go-micro/client/grpc` instead
func MaxRecvMsgSize(s int) client.Option {
	return grpc.MaxRecvMsgSize(s)
}

//
// MaxSendMsgSize set the maximum size of message that client can send.
//
// Deprecated: use `github.com/micro/go-micro/client/grpc` instead
func MaxSendMsgSize(s int) client.Option {
	return grpc.MaxSendMsgSize(s)
}
