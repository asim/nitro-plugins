package whitelist

import (
	"context"

	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/errors"
)

type wrapper struct {
	client.Client
	whitelist map[string]bool
}

func (w *wrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	if !w.whitelist[req.Service()] {
		return errors.Forbidden("go.micro.rpc", "forbidden")
	}

	return w.Client.Call(ctx, req, rsp, opts...)
}

func newClient(services ...string) client.Client {
	whitelist := make(map[string]bool)

	for _, service := range services {
		whitelist[service] = true
	}

	return &wrapper{client.DefaultClient, whitelist}
}
