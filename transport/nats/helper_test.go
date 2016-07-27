package nats_test

import (
	"os"
	"testing"

	"github.com/micro/go-micro/transport"
	"github.com/micro/go-plugins/transport/nats"
)

func setUpTestTransport(tb testing.TB) (transport.Transport, error) {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		tb.Skip("NATS_URL is not set")
		return nil, nil
	}

	return nats.NewTransport(transport.Addrs(natsURL)), nil
}
