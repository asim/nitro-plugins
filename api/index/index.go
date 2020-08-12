// Package index is a micro plugin for stripping a path index
package index

import (
	"errors"
	"net/http"
	"strings"

	"github.com/micro/cli/v2"
	"github.com/micro/micro/v2/plugin"
)

type index struct {
	url string
	r   *response
}

type response struct {
	status int
	header http.Header
	body   []byte
}

func (i *index) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:   "index_service",
			Usage:  "Service/Method to route index to. Specified without namespace e.g greeter/say/hello",
			EnvVars: []string{"INDEX_SERVICE"},
		},
		// flags for response instead of service
		&cli.IntFlag{
			Name:   "index_status",
			Usage:  "HTTP status code for response",
			EnvVars: []string{"INDEX_STATUS"},
		},
		&cli.StringFlag{
			Name:   "index_header",
			Usage:  "Comma separated list of key-value pairs for response header",
			EnvVars: []string{"INDEX_HEADER"},
		},
		&cli.StringFlag{
			Name:   "index_body",
			Usage:  "Body of the response",
			EnvVars: []string{"INDEX_BODY"},
		},
	}
}

func (i *index) Commands() []*cli.Command {
	return nil
}

func (i *index) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// set path if index
			if r.URL.Path == "/" {
				// write response if we have one
				if i.r != nil {
					// write headers
					for k, v := range i.r.header {
						w.Header().Add(k, strings.Join(v, ","))
					}
					// write status
					w.WriteHeader(i.r.status)
					// write body
					w.Write(i.r.body)
					// we're done
					return
				}

				// no response, rewrite url to service
				r.URL.Path = i.url
			}

			// serve request
			h.ServeHTTP(w, r)
		})
	}
}

func (i *index) Init(ctx *cli.Context) error {
	// check if there's a service
	if service := ctx.String("index_service"); len(service) > 0 {
		i.url = "/" + service
	}

	r := new(response)

	// check if there's response content
	// check status
	if status := ctx.Int("index_status"); status > 0 {
		r.status = status
	}
	// check header
	if header := ctx.String("index_header"); len(header) > 0 {
		head := make(http.Header)
		for _, h := range strings.Split(header, ",") {
			if parts := strings.Split(h, ":"); len(parts) == 2 {
				head[parts[0]] = []string{parts[1]}
			}
		}
		r.header = head
	}
	// check body
	if body := ctx.String("index_body"); len(body) > 0 {
		r.body = []byte(body)
	}

	// if we have status then use the response
	if r.status > 0 {
		i.r = r
	}

	if len(i.url) == 0 && i.r == nil {
		return errors.New("neither index service or response specified")
	}
	return nil
}

func (i *index) String() string {
	return "index"
}

// NewPlugin creates a new plugin expecting the service specified via flag
func NewPlugin() plugin.Plugin {
	return &index{}
}

// WithService creates an index plugin with a service
func WithService(service string) plugin.Plugin {
	return &index{
		url: "/" + service,
	}
}

// WithContent will write the given status, header and body
func WithResponse(status int, header http.Header, body []byte) plugin.Plugin {
	return &index{
		r: &response{
			status: status,
			header: header,
			body:   body,
		},
	}
}
