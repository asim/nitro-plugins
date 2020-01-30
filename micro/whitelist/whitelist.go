// Package whitelist is a micro plugin for whitelisting service requests
package whitelist

import (
	"net/http"
	"strings"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/micro/v2/plugin"
)

type whitelist struct {
	services map[string]bool
}

func (w *whitelist) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "rpc_whitelist",
			Usage:  "Comma separated whitelist of allowed services for RPC calls",
			EnvVar: "RPC_WHITELIST",
		},
	}
}

func (w *whitelist) Commands() []cli.Command {
	return nil
}

func (w *whitelist) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return h
	}
}

func (w *whitelist) Init(ctx *cli.Context) error {
	if whitelist := ctx.String("rpc_whitelist"); len(whitelist) > 0 {
		client.DefaultClient = newClient(strings.Split(whitelist, ",")...)
	}
	return nil
}

func (w *whitelist) String() string {
	return "whitelist"
}

func NewPlugin() plugin.Plugin {
	return NewRPCWhitelist()
}

func NewRPCWhitelist(services ...string) plugin.Plugin {
	list := make(map[string]bool)

	for _, service := range services {
		list[service] = true
	}

	return &whitelist{
		services: list,
	}
}
