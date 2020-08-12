// Package ip_allow is a micro plugin for allowing ip addresses
package ip_allow

import (
	"net"
	"net/http"
	"strings"

	"github.com/micro/cli/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/micro/v2/plugin"
)

type allow struct {
	cidrs map[string]*net.IPNet
	ips   map[string]bool
}

func (w *allow) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "ip_allow",
			Usage:   "Comma separated list of allowed IPs",
			EnvVars: []string{"IP_ALLOW"},
		},
	}
}

func (w *allow) load(ips ...string) {
	for _, ip := range ips {
		parts := strings.Split(ip, "/")

		switch len(parts) {
		// assume just an ip
		case 1:
			w.ips[ip] = true
		case 2:
			// parse cidr
			_, ipnet, err := net.ParseCIDR(ip)
			if err != nil {
				log.Fatalf("[ip_allow] failed to parse %v: %v", ip, err)
			}
			w.cidrs[ipnet.String()] = ipnet
		default:
			log.Fatalf("[ip_allow] failed to parse %v", ip)
		}
	}

}

func (w *allow) match(ip string) bool {
	// make ip
	nip := net.ParseIP(ip)

	// check ips
	if ok := w.ips[nip.String()]; ok {
		return true
	}

	// check cidrs
	for _, cidr := range w.cidrs {
		if cidr.Contains(nip) {
			return true
		}
	}

	// no match
	return false
}

func (w *allow) Commands() []*cli.Command {
	return nil
}

func (w *allow) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			// check remote addr; if we can't parse it passes through
			if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
				// reject if no match
				if !w.match(ip) {
					http.Error(rw, "forbidden", 403)
					return
				}
			}

			// serve the request
			h.ServeHTTP(rw, r)
		})
	}
}

func (w *allow) Init(ctx *cli.Context) error {
	if allow := ctx.String("ip_allow"); len(allow) > 0 {
		w.load(strings.Split(allow, ",")...)
	}
	return nil
}

func (w *allow) String() string {
	return "ip_allow"
}

func NewPlugin() plugin.Plugin {
	return NewIPAllow()
}

func NewIPAllow(ips ...string) plugin.Plugin {
	// create plugin
	w := &allow{
		cidrs: make(map[string]*net.IPNet),
		ips:   make(map[string]bool),
	}

	// load ips
	w.load(ips...)

	return w
}
