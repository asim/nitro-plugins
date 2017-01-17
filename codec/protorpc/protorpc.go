// Package protorpc provides a proto-rpc codec
package protorpc

import (
	"io"

	"github.com/micro/go-micro/codec"
	"github.com/micro/go-micro/codec/protorpc"
)

/*
	PROTO-RPC is one of the default codecs in go-micro.
	Content type used is application/protobuf or application/proto-rpc
	Implementation here https://godoc.org/github.com/micro/go-micro/codec/protorpc
	We have a link here for completeness
*/

func NewCodec(rwc io.ReadWriteCloser) codec.Codec {
	return protorpc.NewCodec(rwc)
}
