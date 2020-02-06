package header

import (
	"net/http"
	"strings"

	"github.com/micro/cli/v2"
	"github.com/micro/micro/v2/plugin"
)

type header struct {
	hd map[string]string
}

func (h *header) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:   "header",
			Usage:  "Headers to be set in the http response",
			EnvVars: []string{"HEADER"},
		},
	}
}

func (h *header) Commands() []*cli.Command {
	return nil
}

func (h *header) Handler() plugin.Handler {
	return func(ha http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for k, v := range h.hd {
				// set header
				r.Header.Set(k, v)
			}
			// exec handler
			ha.ServeHTTP(w, r)
		})
	}
}

func (h *header) Init(ctx *cli.Context) error {
	header := ctx.StringSlice("header")

	// no op
	if len(header) == 0 {
		return nil
	}

	// iterate the string slice
	for _, pair := range header {
		parts := strings.Split(pair, "=")
		if len(parts) < 2 {
			continue
		}
		// set key-vals
		h.hd[parts[0]] = strings.Join(parts[1:], "=")
	}

	return nil
}

func (h *header) String() string {
	return "header"
}

func NewPlugin() plugin.Plugin {
	return &header{
		hd: make(map[string]string),
	}
}
