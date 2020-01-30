// Package disable_rpc disables the /rpc endpoint
package disable_rpc

import (
	"net/http"

	"github.com/micro/cli/v2"
	"github.com/micro/micro/v2/plugin"
)

type disable_rpc struct{}

func (i *disable_rpc) Flags() []cli.Flag {
	return nil
}

func (r *disable_rpc) Commands() []cli.Command {
	return nil
}

func (r *disable_rpc) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/rpc" {
				http.Error(w, "forbidden", 403)
				return
			}
			// serve request
			h.ServeHTTP(w, r)
		})
	}
}

func (r *disable_rpc) Init(ctx *cli.Context) error {
	return nil
}

func (r *disable_rpc) String() string {
	return "disable_rpc"
}

// NewPlugin creates a new plugin expecting the service specified via flag
func NewPlugin() plugin.Plugin {
	return &disable_rpc{}
}
