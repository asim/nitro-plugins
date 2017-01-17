// Package opentracing provides wrappers for OpenTracing
package opentracing

import (
	"fmt"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
)

type otWrapper struct {
	ot opentracing.Tracer
	client.Client
}

func (o *otWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	md, _ := metadata.FromContext(ctx)
	name := fmt.Sprintf("%s.%s", req.Service(), req.Method())
	var sp opentracing.Span
	wireContext, err := o.ot.Extract(opentracing.TextMap, opentracing.TextMapCarrier(md))
	if err != nil {
		sp = o.ot.StartSpan(name)
	} else {
		sp = o.ot.StartSpan(name, opentracing.ChildOf(wireContext))
	}
	defer sp.Finish()
	if err := sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.TextMapCarrier(md)); err != nil {
		return err
	}
	ctx = metadata.NewContext(ctx, md)
	return o.Client.Call(ctx, req, rsp, opts...)
}

// NewClientWrapper accepts an open tracing Trace and returns a Client Wrapper
func NewClientWrapper(ot opentracing.Tracer) client.Wrapper {
	return func(c client.Client) client.Client {
		return &otWrapper{ot, c}
	}
}

// NewHandlerWrapper accepts an opentracing Tracer and returns a Handler Wrapper
func NewHandlerWrapper(ot opentracing.Tracer) server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			md, _ := metadata.FromContext(ctx)
			name := fmt.Sprintf("%s.%s", req.Service(), req.Method())
			var sp opentracing.Span
			wireContext, err := ot.Extract(opentracing.TextMap, opentracing.TextMapCarrier(md))
			if err != nil {
				sp = ot.StartSpan(name)
			} else {
				sp = ot.StartSpan(name, opentracing.ChildOf(wireContext))
			}
			defer sp.Finish()
			if err := sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.TextMapCarrier(md)); err != nil {
				return err
			}
			ctx = metadata.NewContext(ctx, md)
			return h(ctx, req, rsp)
		}
	}
}
