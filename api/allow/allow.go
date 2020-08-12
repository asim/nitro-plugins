// Package allow is a micro plugin for allowing service requests
package allow

import (
	"net/http"
	"strings"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/micro/v2/plugin"
)

type allow struct {
	services map[string]bool
}

func (w *allow) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:   "rpc_allow",
			Usage:  "Comma separated allow of allowed services for RPC calls",
			EnvVars: []string{"RPC_ALLOW"},
		},
	}
}

func (w *allow) Commands() []*cli.Command {
	return nil
}

func (w *allow) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return h
	}
}

func (w *allow) Init(ctx *cli.Context) error {
	if allow := ctx.String("rpc_allow"); len(allow) > 0 {
		client.DefaultClient = newClient(strings.Split(allow, ",")...)
	}
	return nil
}

func (w *allow) String() string {
	return "allow"
}

func NewPlugin() plugin.Plugin {
	return NewRPCAllow()
}

func NewRPCAllow(services ...string) plugin.Plugin {
	list := make(map[string]bool)

	for _, service := range services {
		list[service] = true
	}

	return &allow{
		services: list,
	}
}
