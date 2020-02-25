// Package metrics provides metrics for micro
package metrics

import (
	"net/http"

	"github.com/micro/cli/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/micro/v2/plugin"

	// prometheus metrics
	"github.com/micro/go-plugins/micro/metrics/prometheus/v2"
)

type Metrics struct {
	Path string

	Provider Provider
}

// Provider is a metrics provider
type Provider interface {
	Handler(h http.Handler) http.Handler
}

func (m *Metrics) Handler(h http.Handler) http.Handler {
	return m.Provider.Handler(h)
}

// NewPlugin returns a new metrics plugin
func NewPlugin() plugin.Plugin {
	metrics := new(Metrics)

	return plugin.NewPlugin(
		plugin.WithName("metrics"),
		plugin.WithFlag(
			&cli.StringFlag{
				Name:  "metrics",
				Usage: "Specify the type of metrics provider e.g prometheus",
			},
		),
		plugin.WithHandler(metrics.Handler),
		plugin.WithInit(func(ctx *cli.Context) error {
			provider := ctx.String("metrics")

			switch provider {
			case "prometheus":
				metrics.Provider = prometheus.New()
				log.Info("Loaded prometheus metrics at /metrics")
			}

			return nil
		}),
	)
}
