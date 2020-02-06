package router

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/memory"
)

func TestRouter(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	routes := []Route{
		{
			Request: Request{
				Method: "GET",
				Host:   l.Addr().String(),
				Path:   "/",
			},
			Response: Response{
				StatusCode: 302,
				Header: map[string]string{
					"Location": "http://example.com",
				},
			},
			Weight: 1.0,
		},
		{
			Request: Request{
				Method: "POST",
				Host:   l.Addr().String(),
				Path:   "/bar",
			},
			Response: Response{
				StatusCode: 301,
				Header: map[string]string{
					"Location": "http://foo.bar.com",
				},
			},
			Weight: 1.0,
		},
		{
			Request: Request{
				Method: "GET",
				Host:   l.Addr().String(),
				Path:   "/foobar",
			},
			ProxyURL: URL{
				Scheme: "https",
				Host:   "www.google.com",
				Path:   "/",
			},
			Weight: 1.0,
			Type:   "proxy",
		},
	}

	apiConfig := map[string]interface{}{
		"api": map[string]interface{}{
			"routes": routes,
		},
	}

	b, _ := json.Marshal(apiConfig)
	m := memory.NewSource(memory.WithJSON(b))
	conf := config.NewConfig(config.WithSource(m))
	r := NewRouter(Config(conf))

	wr := r.Handler()
	h := wr(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", 404)
	}))

	go http.Serve(l, h)

	ErrRedirect := errors.New("redirect")

	c := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return ErrRedirect
		},
	}

	for _, route := range routes {
		var rsp *http.Response
		var err error

		switch route.Request.Method {
		case "GET":
			rsp, err = c.Get("http://" + route.Request.Host + route.Request.Path)
		case "POST":
			rsp, err = c.Post("http://"+route.Request.Host+route.Request.Path, "application/json", bytes.NewBuffer(nil))
		}

		if err != nil {
			urlErr, ok := err.(*url.Error)
			if ok && urlErr.Err == ErrRedirect {
				// skip
			} else {
				t.Fatal(err)
			}
		}

		if route.Type == "proxy" {
			if rsp.StatusCode >= 400 {
				t.Fatalf("Expected healthy response got %d", rsp.StatusCode)
			}
			continue
		}

		if rsp.StatusCode != route.Response.StatusCode {
			t.Fatalf("Expected code %d got %d", route.Response.StatusCode, rsp.StatusCode)
		}

		loc := rsp.Header.Get("Location")
		if loc != route.Response.Header["Location"] {
			t.Fatalf("Expected Location %s got %s", route.Response.Header["Location"], loc)
		}
	}
}
