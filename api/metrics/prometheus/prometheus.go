// Package prometheus provides prometheus metrics via a http handler
package prometheus

import (
	"net/http"

	p "github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct{}

func (m *Metrics) Handler(h http.Handler) http.Handler {
	ph := p.Handler()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// serve prometheus handler at /metrics
		if r.URL.Path == "/metrics" {
			ph.ServeHTTP(w, r)
			return
		}
		// otherwise serve everything
		h.ServeHTTP(w, r)
	})
}

func New() *Metrics {
	return new(Metrics)
}
