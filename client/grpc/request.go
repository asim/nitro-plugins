package grpc

import (
	"fmt"
	"github.com/micro/go-micro/client"
	"strings"
)

type grpcRequest struct {
	service     string
	method      string
	contentType string
	request     interface{}
	opts        client.RequestOptions
}

func methodToGRPC(pkg, method string, request interface{}) string {
	// no method or already grpc method
	if len(method) == 0 || method[0] == '/' {
		return method
	}
	// assume method is Foo.Bar
	mParts := strings.Split(method, ".")
	if len(mParts) != 2 {
		return method
	}
	// return /pkg.Foo/Bar
	return fmt.Sprintf("/%s.%s/%s", pkg, mParts[0], mParts[1])
}

func newGRPCRequest(service, method string, request interface{}, contentType string, reqOpts ...client.RequestOption) client.Request {
	var opts client.RequestOptions
	for _, o := range reqOpts {
		o(&opts)
	}

	return &grpcRequest{
		service:     service,
		method:      method,
		request:     request,
		contentType: contentType,
		opts:        opts,
	}
}

func (g *grpcRequest) ContentType() string {
	return g.contentType
}

func (g *grpcRequest) Service() string {
	return g.service
}

func (g *grpcRequest) Method() string {
	return g.method
}

func (g *grpcRequest) Request() interface{} {
	return g.request
}

func (g *grpcRequest) Stream() bool {
	return g.opts.Stream
}
