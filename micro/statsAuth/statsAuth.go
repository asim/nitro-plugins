// statsAuth enables basic auth on the /stats endpoint
package statsAuth

import (
	"net/http"
	"strings"

	"github.com/micro/cli"
	"github.com/micro/micro/plugin"
)

const (
	defaultRealm = "Access to stats is restricted"
)

type statsAuth struct {
	User  string
	Pass  string
	Realm string
}

func (sa *statsAuth) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "stats_auth_user",
			Usage:  "Username used for basic auth for /stats endpoint",
			EnvVar: "STATS_AUTH_USER",
		},
		cli.StringFlag{
			Name:   "stats_auth_pass",
			Usage:  "Password used for basic auth for /stats endpoint",
			EnvVar: "STATS_AUTH_PASS",
		},
		cli.StringSliceFlag{
			Name:   "stats_auth_realm",
			Usage:  "Realm used for basic auth for /stats endpoint. Optional. Defaults to " + defaultRealm,
			EnvVar: "STATS_AUTH_REALM",
		},
	}
}

func (sa *statsAuth) Commands() []cli.Command {
	return nil
}

func (sa *statsAuth) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/stats" {
				h.ServeHTTP(w, r)
				return
			}
			if u, p, ok := r.BasicAuth(); ok {
				if u == sa.User && p == sa.Pass {
					h.ServeHTTP(w, r)
					return
				}
			}
			w.Header().Add("WWW-Authenticate", sa.Realm)
			w.WriteHeader(http.StatusUnauthorized)
			return
		})
	}
}

func (sa *statsAuth) Init(ctx *cli.Context) error {
	sa.User = ctx.String("stats_auth_user")
	sa.Pass = ctx.String("stats_auth_pass")
	if ctx.IsSet("stats_auth_realm") {
		sa.Realm = strings.Join(ctx.StringSlice("stats_auth_realm"), " ")
	} else {
		sa.Realm = defaultRealm
	}
	return nil
}

func (sa *statsAuth) String() string {
	return "statsAuth"
}

func New() plugin.Plugin {
	return &statsAuth{}
}
