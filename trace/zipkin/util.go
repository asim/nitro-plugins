package zipkin

import (
	"github.com/micro/go-os/trace"
)

func SpanFromHeader(md map[string]string) (*trace.Span, bool) {
	var debug bool
	if md[SampleHeader] == "1" {
		debug = true
	}

	// can we get span header and trace header?
	if len(md[SpanHeader]) == 0 && len(md[TraceHeader]) == 0 {
		return nil, false
	}

	return &trace.Span{
		Id:       md[SpanHeader],
		TraceId:  md[TraceHeader],
		ParentId: md[ParentHeader],
		Debug:    debug,
	}, true
}

func HeaderWithSpan(md map[string]string, s *trace.Span) map[string]string {
	sample := "0"
	if s.Debug {
		sample = "1"
	}
	md[SpanHeader] = s.Id
	md[TraceHeader] = s.TraceId
	md[ParentHeader] = s.ParentId
	md[SampleHeader] = sample
	return md
}
