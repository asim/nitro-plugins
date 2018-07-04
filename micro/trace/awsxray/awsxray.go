// Package awsxray is a micro plugin for whitelisting service requests
package awsxray

import (
	"net/http"

	_ "github.com/asim/go-awsxray"
	"github.com/micro/cli"
	"github.com/micro/go-micro/client"
	"github.com/micro/micro/plugin"
	xray "github.com/saturn-sonic/go-plugins/wrapper/trace/awsxray"
)

type awsXRay struct {
	opts Options
	x    *awsxray.AWSXRay
}

func (x *awsXRay) Flags() []cli.Flag {
	return nil
}

func (x *awsXRay) Commands() []cli.Command {
	return nil
}

func (x *awsXRay) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s := newSegment(x.opts.Name, r)
			// use our own writer
			xw := &writer{w, 200}
			// serve request
			h.ServeHTTP(xw, r)
			// set status
			complete(s, xw.status)
			// send segment asynchronously
			go x.x.Record(s)
		})
	}
}

func (x *awsXRay) Init(ctx *cli.Context) error {
	opts := []xray.Option{
		xray.WithName(x.opts.Name),
		xray.WithClient(x.opts.Client),
		xray.WithDaemon(x.opts.Daemon),
	}

	// setup client
	c := client.DefaultClient
	c = xray.NewClientWrapper(opts...)(c)
	c.Init(client.WrapCall(xray.NewCallWrapper(opts...)))
	client.DefaultClient = c
	return nil
}

func (x *awsXRay) String() string {
	return "awsxray"
}

func NewXRayPlugin(opts ...Option) plugin.Plugin {
	options := Options{
		Name:   "go.micro.http",
		Daemon: "localhost:2000",
	}

	for _, o := range opts {
		o(&options)
	}

	return &awsXRay{
		opts: options,
		x: awsxray.New(
			awsxray.WithDaemon(options.Daemon),
			awsxray.WithClient(options.Client),
		),
	}
}
