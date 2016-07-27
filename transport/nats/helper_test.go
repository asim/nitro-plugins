package nats_test

import (
	"fmt"
	"os"

	"github.com/micro/go-micro/transport"
	"github.com/micro/go-plugins/transport/nats"
)

func setUpTestTransport() (transport.Transport, error) {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		return nil, fmt.Errorf("NATS_URL is not set")
	}

	return nats.NewTransport(transport.Addrs(natsURL)), nil
}
