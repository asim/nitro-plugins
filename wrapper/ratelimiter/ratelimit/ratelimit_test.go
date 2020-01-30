package ratelimit

import (
	"fmt"
	"testing"
	"time"

	"context"

	"github.com/juju/ratelimit"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/registry/memory"
	"github.com/micro/go-micro/v2/server"
)

type testHandler struct{}
type TestRequest struct{}
type TestResponse struct{}

func (t *testHandler) Method(ctx context.Context, req *TestRequest, rsp *TestResponse) error {
	return nil
}

func TestRateClientLimit(t *testing.T) {
	// setup
	r := memory.NewRegistry()
	s := selector.NewSelector(selector.Registry(r))

	testRates := []int{1, 10, 20, 100}

	for _, limit := range testRates {
		b := ratelimit.NewBucketWithRate(float64(limit), int64(limit))

		c := client.NewClient(
			// set the selector
			client.Selector(s),
			// add the breaker wrapper
			client.Wrap(NewClientWrapper(b, false)),
		)

		req := c.NewRequest(
			"test.service",
			"Test.Method",
			&TestRequest{},
			client.WithContentType("application/json"),
		)
		rsp := TestResponse{}

		for j := 0; j < limit; j++ {
			err := c.Call(context.TODO(), req, &rsp)
			e := errors.Parse(err.Error())
			if e.Code == 429 {
				t.Errorf("Unexpected rate limit error: %v", err)
			}
		}

		err := c.Call(context.TODO(), req, rsp)
		e := errors.Parse(err.Error())
		if e.Code != 429 {
			t.Errorf("Expected rate limit error, got: %v", err)
		}
	}
}

func TestRateServerLimit(t *testing.T) {
	// setup
	r := memory.NewRegistry()
	s := selector.NewSelector(selector.Registry(r))

	testRates := []int{1, 10, 20}

	for _, limit := range testRates {
		b := ratelimit.NewBucketWithRate(float64(limit), int64(limit))
		c := client.NewClient(client.Selector(s))

		name := fmt.Sprintf("test.service.%d", limit)

		s := server.NewServer(
			server.Name(name),
			// add registry
			server.Registry(r),
			// add the breaker wrapper
			server.WrapHandler(NewHandlerWrapper(b, false)),
		)

		type Test struct {
			*testHandler
		}

		s.Handle(
			s.NewHandler(&Test{new(testHandler)}),
		)

		if err := s.Start(); err != nil {
			t.Fatalf("Unexpected error starting server: %v", err)
		}

		req := c.NewRequest(name, "Test.Method", &TestRequest{}, client.WithContentType("application/json"))
		rsp := TestResponse{}

		for j := 0; j < limit; j++ {
			if err := c.Call(context.TODO(), req, &rsp); err != nil {
				t.Fatalf("Unexpected request error: %v", err)
			}
		}

		err := c.Call(context.TODO(), req, &rsp)
		if err == nil {
			t.Fatalf("Expected rate limit error, got nil: rate %d, err %v", limit, err)
		}

		e := errors.Parse(err.Error())
		if e.Code != 429 {
			t.Fatalf("Expected rate limit error, got %v", err)
		}

		s.Stop()

		// artificial test delay
		time.Sleep(time.Millisecond * 20)
	}
}
