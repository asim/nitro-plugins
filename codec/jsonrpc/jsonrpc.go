// Package jsonrpc provides a json-rpc 2.0 codec
package jsonrpc

import (
	"io"

	"github.com/micro/go-micro/codec"
	"github.com/micro/go-micro/codec/jsonrpc"
)

/*
	JSON-RPC is one of the default codecs in go-micro.
	Content type used is application/json or application/json-rpc
	Implementation here https://godoc.org/github.com/micro/go-micro/codec/jsonrpc
	We have a link here for completeness
*/

func NewCodec(rwc io.ReadWriteCloser) codec.Codec {
	return jsonrpc.NewCodec(rwc)
}
