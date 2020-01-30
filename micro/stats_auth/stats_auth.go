// stats_auth enables basic auth on the /stats endpoint
package stats_auth

import (
	"net/http"

	"github.com/micro/cli/v2"
	"github.com/micro/micro/v2/plugin"
)

const (
	defaultRealm = "Access to stats is restricted"
)

type stats_auth struct {
	User  string
	Pass  string
	Realm string
}

func (sa *stats_auth) Flags() []cli.Flag {
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
		cli.StringFlag{
			Name:   "stats_auth_realm",
			Usage:  "Realm used for basic auth for /stats endpoint. Escape spaces to add multiple words. Optional. Defaults to " + defaultRealm,
			EnvVar: "STATS_AUTH_REALM",
		},
	}
}

func (sa *stats_auth) Commands() []cli.Command {
	return nil
}

func (sa *stats_auth) Handler() plugin.Handler {
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

func (sa *stats_auth) Init(ctx *cli.Context) error {
	sa.User = ctx.String("stats_auth_user")
	sa.Pass = ctx.String("stats_auth_pass")
	if ctx.IsSet("stats_auth_realm") {
		sa.Realm = ctx.String("stats_auth_realm")
	} else {
		sa.Realm = defaultRealm
	}
	return nil
}

func (sa *stats_auth) String() string {
	return "stats_auth"
}

func NewPlugin() plugin.Plugin {
	return &stats_auth{}
}
