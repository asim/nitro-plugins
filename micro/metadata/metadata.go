package metadata

import (
	"net/http"
	"strings"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/micro/v2/plugin"
)

type metadata struct {
	md map[string]string
}

func (m *metadata) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:   "metadata",
			Usage:  "A list of key-value pairs to be forwarded as metadata",
			EnvVars: []string{"METADATA"},
		},
	}
}

func (m *metadata) Commands() []*cli.Command {
	return nil
}

func (m *metadata) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for k, v := range m.md {
				// set metadata
				r.Header.Set(k, v)
			}
			// exec handler
			h.ServeHTTP(w, r)
		})
	}
}

func (m *metadata) Init(ctx *cli.Context) error {
	md := ctx.StringSlice("metadata")

	// no op
	if len(md) == 0 {
		return nil
	}

	// iterate the string slice
	for _, pair := range md {
		parts := strings.Split(pair, "=")
		if len(parts) < 2 {
			continue
		}
		// set key-vals
		m.md[parts[0]] = strings.Join(parts[1:], "=")
	}

	// wrap the client
	client.DefaultClient = newClient(m.md)

	return nil
}

func (m *metadata) String() string {
	return "metadata"
}

func NewPlugin() plugin.Plugin {
	return &metadata{
		md: make(map[string]string),
	}
}
