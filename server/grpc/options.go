package grpc

import (
	"crypto/tls"

	"github.com/micro/go-micro/server"
	grpcServer "github.com/micro/go-micro/server/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
)

// gRPC Codec to be used to encode/decode requests for a given content type
func Codec(contentType string, c encoding.Codec) server.Option {
	return grpcServer.Codec(contentType, c)
}

// AuthTLS should be used to setup a secure authentication using TLS
func AuthTLS(t *tls.Config) server.Option {
	return grpcServer.AuthTLS(t)
}

// Options to be used to configure gRPC options
func Options(opts ...grpc.ServerOption) server.Option {
	return grpcServer.Options(opts...)
}

//
// MaxMsgSize set the maximum message in bytes the server can receive and
// send.  Default maximum message size is 4 MB.
//
func MaxMsgSize(s int) server.Option {
	return grpcServer.MaxMsgSize(s)
}
