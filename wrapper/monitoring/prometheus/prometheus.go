package prometheus

import (
	"context"
	"fmt"

	"github.com/micro/go-micro/v2/server"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	defaultMetricPrefix = "micro"
)

func NewHandlerWrapper(opts ...server.Option) server.HandlerWrapper {
	md := make(map[string]string)
	sopts := server.Options{}

	for _, opt := range opts {
		opt(&sopts)
	}

	for k, v := range sopts.Metadata {
		md[fmt.Sprintf("%s_%s", defaultMetricPrefix, k)] = v
	}
	if len(sopts.Name) > 0 {
		md[fmt.Sprintf("%s_%s", defaultMetricPrefix, "name")] = sopts.Name
	}
	if len(sopts.Id) > 0 {
		md[fmt.Sprintf("%s_%s", defaultMetricPrefix, "id")] = sopts.Id
	}
	if len(sopts.Version) > 0 {
		md[fmt.Sprintf("%s_%s", defaultMetricPrefix, "version")] = sopts.Version
	}

	opsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "micro",
			Name:      "request_total",
			Help:      "How many go-micro requests processed, partitioned by method and status",
		},
		[]string{"method", "status"},
	)

	timeCounterSummary := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "micro",
			Name:      "upstream_latency_microseconds",
			Help:      "Service backend method request latencies in microseconds",
		},
		[]string{"method"},
	)

	timeCounterHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "micro",
			Name:      "request_duration_seconds",
			Help:      "Service method request time in seconds",
		},
		[]string{"method"},
	)

	wrapreg := prometheus.WrapRegistererWith(md, prometheus.DefaultRegisterer)
	wrapreg.MustRegister(
		//prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		//prometheus.NewGoCollector(),
		opsCounter,
		timeCounterSummary,
		timeCounterHistogram,
	)

	prometheus.DefaultRegisterer = wrapreg

	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			name := req.Endpoint()

			timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
				us := v * 1000000 // make microseconds
				timeCounterSummary.WithLabelValues(name).Observe(us)
				timeCounterHistogram.WithLabelValues(name).Observe(v)
			}))
			defer timer.ObserveDuration()

			err := fn(ctx, req, rsp)
			if err == nil {
				opsCounter.WithLabelValues(name, "success").Inc()
			} else {
				opsCounter.WithLabelValues(name, "fail").Inc()
			}

			return err
		}
	}
}
