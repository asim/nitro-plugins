package metadata

import (
	"context"

	"github.com/micro/go-micro/v2/client"
	meta "github.com/micro/go-micro/v2/metadata"
)

type wrapper struct {
	client.Client
	md map[string]string
}

func (w *wrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	md, ok := meta.FromContext(ctx)
	if !ok {
		md = make(meta.Metadata)
	}

	// set our meta
	for k, v := range w.md {
		md[k] = v
	}

	ctx = meta.NewContext(ctx, md)
	return w.Client.Call(ctx, req, rsp, opts...)
}

func newClient(md map[string]string) client.Client {
	return &wrapper{client.DefaultClient, md}
}
