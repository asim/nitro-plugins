package victoriametrics

import (
	"context"
	"fmt"
	"strings"
	"time"

	metrics "github.com/VictoriaMetrics/metrics"
	"github.com/micro/go-micro/server"
)

var (
	defaultMetricPrefix = "micro"
	metaLabels          []string
)

func getName(name string, md map[string]interface{}) string {
	labels := make([]string, 0, len(metaLabels)+len(md))
	labels = append(labels, metaLabels...)

	for k, v := range md {
		labels = append(labels, fmt.Sprintf(`%s="%v"`, k, v))
	}

	if len(labels) > 0 {
		return fmt.Sprintf(`%s_%s{%s}`, defaultMetricPrefix, name, strings.Join(labels, ","))
	}
	return fmt.Sprintf(`%s_%s`, defaultMetricPrefix, name)
}

func NewHandlerWrapper(opts ...server.Option) server.HandlerWrapper {
	sopts := server.Options{}

	for _, opt := range opts {
		opt(&sopts)
	}

	metadata := make(map[string]string, len(sopts.Metadata))
	for k, v := range sopts.Metadata {
		metadata[fmt.Sprintf("%s_%s", defaultMetricPrefix, k)] = v
	}
	if len(sopts.Name) > 0 {
		metadata[fmt.Sprintf("%s_%s", defaultMetricPrefix, "name")] = sopts.Name
	}
	if len(sopts.Id) > 0 {
		metadata[fmt.Sprintf("%s_%s", defaultMetricPrefix, "id")] = sopts.Id
	}
	if len(sopts.Version) > 0 {
		metadata[fmt.Sprintf("%s_%s", defaultMetricPrefix, "version")] = sopts.Version
	}
	metaLabels = make([]string, 0, len(metadata))
	for k, v := range metadata {
		metaLabels = append(metaLabels, fmt.Sprintf(`%s="%v"`, k, v))
	}

	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			name := req.Endpoint()
			timeCounterSummary := metrics.GetOrCreateSummary(
				getName("upstream_latency_seconds", map[string]interface{}{"method": name}),
			)
			timeCounterHistogram := metrics.GetOrCreateSummary(
				getName("request_duration_seconds", map[string]interface{}{"method": name}),
			)

			ts := time.Now()
			err := fn(ctx, req, rsp)
			te := time.Since(ts)

			timeCounterSummary.Update(float64(te.Seconds()))
			timeCounterHistogram.Update(te.Seconds())
			if err == nil {
				metrics.GetOrCreateCounter(getName("request_total", map[string]interface{}{"method": name, "status": "success"})).Inc()
			} else {
				metrics.GetOrCreateCounter(getName("request_total", map[string]interface{}{"method": name, "status": "failure"})).Inc()
			}

			return err
		}
	}
}
