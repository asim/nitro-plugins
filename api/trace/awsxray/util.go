package awsxray

import (
	"crypto/rand"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/asim/go-awsxray"
)

// complete sets the response status and end time
func complete(s *awsxray.Segment, status int) {
	switch {
	case status >= 500:
		s.Fault = true
	case status >= 400:
		s.Error = true
	}
	s.HTTP.Response.Status = status
	s.EndTime = float64(time.Now().Truncate(time.Millisecond).UnixNano()) / 1e9
}

// getIp naively returns an ip for the request
func getIp(r *http.Request) string {
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		for _, ip := range strings.Split(r.Header.Get(h), ",") {
			if len(ip) == 0 {
				continue
			}
			realIP := net.ParseIP(strings.Replace(ip, " ", "", -1))
			return realIP.String()
		}
	}

	// not found in header
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// just return remote addr
		return r.RemoteAddr
	}

	return host
}

// getRandom generates a random byte slice
func getRandom(i int) string {
	b := make([]byte, i)
	for {
		// keep trying till we get it
		if _, err := rand.Read(b); err != nil {
			continue
		}
		return fmt.Sprintf("%x", b)
	}
}

// getTraceId returns trace header or generates a new one
func getTraceId(hdr http.Header) string {
	if h := hdr.Get(awsxray.TraceHeader); len(h) > 0 {
		return awsxray.GetTraceId(h)
	}

	// generate new one, probably a bad idea...
	return fmt.Sprintf("%d-%x-%s", 1, time.Now().Unix(), getRandom(12))
}

// getParentId returns parent header or blank
func getParentId(hdr http.Header) string {
	if h := hdr.Get(awsxray.TraceHeader); len(h) > 0 {
		return awsxray.GetTraceId(h)
	}

	// return nothing
	return ""
}

// newHTTP returns a http struct
func newHTTP(r *http.Request) *awsxray.HTTP {
	scheme := "http"
	host := r.Host

	if len(r.URL.Scheme) > 0 {
		scheme = r.URL.Scheme
	}

	if len(r.URL.Host) > 0 {
		host = r.URL.Host
	}

	return &awsxray.HTTP{
		Request: &awsxray.Request{
			Method:    r.Method,
			URL:       fmt.Sprintf("%s://%s%s", scheme, host, r.URL.Path),
			ClientIP:  getIp(r),
			UserAgent: r.UserAgent(),
		},
		Response: &awsxray.Response{
			Status: 200,
		},
	}
}

// newSegment creates a new segment based on whether we're part of an existing flow
func newSegment(name string, r *http.Request) *awsxray.Segment {
	// attempt to get IDs first
	parentId := getParentId(r.Header)
	traceId := getTraceId(r.Header)

	// now set the trace ID
	traceHdr := r.Header.Get(awsxray.TraceHeader)
	traceHdr = awsxray.SetTraceId(traceHdr, traceId)

	// create segment
	s := &awsxray.Segment{
		Id:        getRandom(8),
		HTTP:      newHTTP(r),
		Name:      name,
		TraceId:   traceId,
		StartTime: float64(time.Now().Truncate(time.Millisecond).UnixNano()) / 1e9,
	}

	// if we have a parent then we are a subsegment
	if len(parentId) > 0 {
		s.ParentId = parentId
		s.Type = "subsegment"
	} else {
		// set a new parent Id
		traceHdr = awsxray.SetParentId(traceHdr, s.Id)
	}

	// now save the header for the future context
	r.Header.Set(awsxray.TraceHeader, traceHdr)

	return s
}
