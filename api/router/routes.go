package router

import (
	"crypto/tls"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

// Routes is the config expected to be loaded
type Routes struct {
	Routes []Route `json:"routes"`
	// TODO: default route
}

// Route describes a single route which is matched
// on Request and if so, will return the Response
type Route struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
	ProxyURL URL      `json:"proxy_url"`
	Priority int      `json:"priority"` // 0 is highest. Used for ordering routes
	Weight   float64  `json:"weight"`   // percentage weight between 0 and 1.0
	Type     string   `json:"type"`     // proxy or response. Response is default
	Insecure bool     `json:"bool"`     // allow insecure certificates
}

// Request describes the expected request and will
// attempt to match all fields specified
type Request struct {
	Method string            `json:"method"`
	Header map[string]string `json:"header"`
	Host   string            `json:"host"`
	Path   string            `json:"path"`
	Query  map[string]string `json:"query"`
	// TODO: RemoteAddr, Body
}

// Response is put into the http.Response for a Request
type Response struct {
	Status     string            `json:"status"`
	StatusCode int               `json:"status_code"`
	Header     map[string]string `json:"header"`
	Body       []byte            `json:"body"`
}

type URL struct {
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
	Path   string `json:"path"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (r Route) Match(req *http.Request) bool {
	// bail on nil
	if len(r.Request.Method) == 0 || len(r.Request.Path) == 0 {
		return false
	}

	// just for ease
	rq := r.Request

	// first level match, quick and dirty
	if (rq.Method == req.Method) && (rq.Host == req.Host) && strings.HasPrefix(req.URL.Path, rq.Path) {
		// skip
	} else {
		return false
	}

	// match headers
	for k, v := range rq.Header {
		// does it match?
		if rv := req.Header.Get(k); rv != v {
			return false
		}
	}

	// match query
	vals := req.URL.Query()
	for k, v := range rq.Query {
		// does it match?
		if rv := vals.Get(k); rv != v {
			return false
		}
	}

	// Now weight it. If already set to 0.0 then return
	// Otherwise rand.Float64
	if r.Weight == 0.0 || r.Weight < rand.Float64() {
		return false
	}

	// we got a match!
	return true
}

func (r Route) Write(w http.ResponseWriter, req *http.Request) {
	// Type: proxy then proxy the request to whatever response is
	if r.Type == "proxy" {
		p := &httputil.ReverseProxy{
			Director: func(rr *http.Request) {
				rr.URL.Host = r.ProxyURL.Host
				rr.URL.Scheme = r.ProxyURL.Scheme
				rr.URL.Path = strings.Replace(rr.URL.Path, r.Request.Path, r.ProxyURL.Path, 1)
				rr.Host = r.ProxyURL.Host
			},
		}
		if r.Insecure {
			p.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}
		}
		p.ServeHTTP(w, req)
		return
	}

	// Type: response or none then set the response

	// set headers
	for k, v := range r.Response.Header {
		w.Header().Set(k, v)
	}
	// set status code
	w.WriteHeader(r.Response.StatusCode)

	// set response
	if len(r.Response.Body) > 0 {
		w.Write(r.Response.Body)
	} else {
		w.Write([]byte(r.Response.Status))
	}
}
