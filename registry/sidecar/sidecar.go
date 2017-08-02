// Package sidecar is a registry plugin for the micro sidecar
package sidecar

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/registry"
)

type sidecar struct {
	opts registry.Options
}

func init() {
	cmd.DefaultRegistries["sidecar"] = NewRegistry
}

func newRegistry(opts ...registry.Option) registry.Registry {
	var options registry.Options
	for _, o := range opts {
		o(&options)
	}
	var addrs []string
	for _, addr := range options.Addrs {
		if len(addr) > 0 {
			addrs = append(addrs, addr)
		}
	}
	if len(addrs) == 0 {
		addrs = []string{"localhost:8081"}
	}
	registry.Addrs(addrs...)(&options)

	return &sidecar{
		opts: options,
	}
}

func (s *sidecar) Register(service *registry.Service, opts ...registry.RegisterOption) error {
	b, err := json.Marshal(service)
	if err != nil {
		return err
	}

	var gerr error
	for _, addr := range s.opts.Addrs {
		scheme := "http"
		if s.opts.Secure {
			scheme = "https"
		}
		url := fmt.Sprintf("%s://%s/registry", scheme, addr)
		rsp, err := http.Post(url, "application/json", bytes.NewReader(b))
		if err != nil {
			gerr = err
			continue
		}
		if rsp.StatusCode != 200 {
			b, err := ioutil.ReadAll(rsp.Body)
			if err != nil {
				return err
			}
			rsp.Body.Close()
			gerr = errors.New(string(b))
			continue
		}
		io.Copy(ioutil.Discard, rsp.Body)
		rsp.Body.Close()
		return nil
	}
	return gerr
}

func (s *sidecar) Deregister(service *registry.Service) error {
	b, err := json.Marshal(service)
	if err != nil {
		return err
	}

	var gerr error
	for _, addr := range s.opts.Addrs {
		scheme := "http"
		if s.opts.Secure {
			scheme = "https"
		}
		url := fmt.Sprintf("%s://%s/registry", scheme, addr)

		req, err := http.NewRequest("DELETE", url, bytes.NewReader(b))
		if err != nil {
			gerr = err
			continue
		}

		rsp, err := http.DefaultClient.Do(req)
		if err != nil {
			gerr = err
			continue
		}

		if rsp.StatusCode != 200 {
			b, err := ioutil.ReadAll(rsp.Body)
			if err != nil {
				return err
			}
			rsp.Body.Close()
			gerr = errors.New(string(b))
			continue
		}

		io.Copy(ioutil.Discard, rsp.Body)
		rsp.Body.Close()
		return nil
	}
	return gerr
}

func (s *sidecar) GetService(service string) ([]*registry.Service, error) {
	var gerr error
	for _, addr := range s.opts.Addrs {
		scheme := "http"
		if s.opts.Secure {
			scheme = "https"
		}

		url := fmt.Sprintf("%s://%s/registry?service=%s", scheme, addr, url.QueryEscape(service))
		rsp, err := http.Get(url)
		if err != nil {
			gerr = err
			continue
		}

		if rsp.StatusCode != 200 {
			b, err := ioutil.ReadAll(rsp.Body)
			if err != nil {
				return nil, err
			}
			rsp.Body.Close()
			gerr = errors.New(string(b))
			continue
		}

		b, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			gerr = err
			continue
		}
		rsp.Body.Close()
		var services []*registry.Service
		if err := json.Unmarshal(b, &services); err != nil {
			gerr = err
			continue
		}
		return services, nil
	}
	return nil, gerr
}

func (s *sidecar) ListServices() ([]*registry.Service, error) {
	var gerr error
	for _, addr := range s.opts.Addrs {
		scheme := "http"
		if s.opts.Secure {
			scheme = "https"
		}
		url := fmt.Sprintf("%s://%s/registry", scheme, addr)
		rsp, err := http.Get(url)
		if err != nil {
			gerr = err
			continue
		}

		if rsp.StatusCode != 200 {
			b, err := ioutil.ReadAll(rsp.Body)
			if err != nil {
				return nil, err
			}
			rsp.Body.Close()
			gerr = errors.New(string(b))
			continue
		}

		b, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			gerr = err
			continue
		}
		rsp.Body.Close()
		var services []*registry.Service
		if err := json.Unmarshal(b, &services); err != nil {
			gerr = err
			continue
		}
		return services, nil
	}
	return nil, gerr
}

func (s *sidecar) Watch() (registry.Watcher, error) {
	watch := func(addr string) (registry.Watcher, error) {
		scheme := "ws"
		if s.opts.Secure {
			scheme = "wss"
		}
		url := fmt.Sprintf("%s://%s/registry", scheme, addr)
		return newWatcher(url)
	}

	var gerr error
	for _, addr := range s.opts.Addrs {
		w, err := watch(addr)
		if err != nil {
			gerr = err
			continue
		}
		return w, nil
	}
	return nil, gerr
}

func (s *sidecar) String() string {
	return "sidecar"
}

func NewRegistry(opts ...registry.Option) registry.Registry {
	return newRegistry(opts...)
}
