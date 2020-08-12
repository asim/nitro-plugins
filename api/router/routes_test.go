package router

import (
	"net/http"
	"net/url"
	"testing"
)

func TestRoutes(t *testing.T) {
	testData := []struct {
		Routes []Route
		Req    *http.Request
		Match  bool
	}{
		{
			Routes: []Route{
				{
					Request: Request{
						Method: "GET",
						Host:   "example.com",
						Path:   "/",
					},
					Weight: 1.0,
				},
				{
					Request: Request{
						Method: "POST",
						Host:   "foo.com",
						Path:   "/bar",
					},
					Weight: 1.0,
				},
				{
					Request: Request{
						Method: "GET",
						Host:   "bif.com",
						Path:   "/bif",
					},
					Weight: 2.0,
					Type:   "proxy",
				},
			},
			Req: &http.Request{
				Method: "GET",
				Host:   "example.com",
				URL: &url.URL{
					Path: "/bif",
				},
			},
			Match: true,
		},
		{
			Routes: []Route{
				{
					Request: Request{
						Method: "GET",
						Host:   "example.com",
						Path:   "/",
					},
					Weight: 1.0,
				},
				{
					Request: Request{
						Method: "POST",
						Host:   "foo.com",
						Path:   "/bar",
					},
					Weight: 1.0,
				},
				{
					Request: Request{
						Method: "GET",
						Host:   "bif.com",
						Path:   "/bif",
					},
					Weight: 2.0,
					Type:   "proxy",
				},
			},
			Req: &http.Request{
				Method: "POST",
				Host:   "foo.com",
				URL: &url.URL{
					Path: "/bar",
				},
			},
			Match: true,
		},
		{
			Routes: []Route{
				{
					Request: Request{
						Method: "GET",
						Host:   "example.com",
						Path:   "/",
					},
					Weight: 1.0,
				},
				{
					Request: Request{
						Method: "POST",
						Host:   "foo.com",
						Path:   "/bar",
					},
					Weight: 1.0,
				},
				{
					Request: Request{
						Method: "GET",
						Host:   "bif.com",
						Path:   "/bif",
					},
					Weight: 2.0,
					Type:   "proxy",
				},
			},
			Req: &http.Request{
				Method: "GET",
				Host:   "bif.com",
				URL: &url.URL{
					Path: "/bif",
				},
			},
			Match: true,
		},
		{
			Routes: []Route{
				{
					Request: Request{
						Method: "GET",
						Host:   "example.com",
						Path:   "/",
					},
					Weight: 1.0,
				},
				{
					Request: Request{
						Method: "POST",
						Host:   "foo.com",
						Path:   "/bar",
					},
					Weight: 1.0,
				},
				{
					Request: Request{
						Method: "GET",
						Host:   "bif.com",
						Path:   "/bif",
					},
					Weight: 2.0,
					Type:   "proxy",
				},
			},
			Req: &http.Request{
				Method: "GET",
				Host:   "bif.com",
				URL: &url.URL{
					Path: "/bif/bam",
				},
			},
			Match: true,
		},
	}

	for _, d := range testData {
		var match bool

		for _, r := range d.Routes {
			if r.Match(d.Req) {
				match = true
				break
			}
		}

		if match != d.Match {
			t.Fatalf("Expected match %t got %t", d.Match, match)
		}
	}
}
