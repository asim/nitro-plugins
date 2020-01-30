// Package uuid is a micro plugin for adding a trace header
package uuid

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/micro/cli/v2"
	"github.com/micro/micro/v2/plugin"
)

type uuidPlugin struct{}

func (u *uuidPlugin) Flags() []cli.Flag {
	return nil
}

func (u *uuidPlugin) Commands() []cli.Command {
	return nil
}

func (u *uuidPlugin) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if td := r.Header.Get("X-Micro-Trace-Id"); len(td) == 0 {
				r.Header.Set("X-Micro-Trace-Id", uuid.New().String())
			}
			h.ServeHTTP(w, r)
		})
	}
}

func (u *uuidPlugin) Init(ctx *cli.Context) error {
	return nil
}

func (u *uuidPlugin) String() string {
	return "uuid"
}

func New() plugin.Plugin {
	return new(uuidPlugin)
}
