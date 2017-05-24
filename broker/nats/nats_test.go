package nats

import (
	"testing"

	"fmt"

	"github.com/micro/go-micro/broker"
)

// TestInitAddrs tests issue #100. Ensures that if the addrs is set by an option in init it will be used.
func TestInitAddrs(t *testing.T) {
	nb := NewBroker()

	addr1, addr2 := "192.168.10.1:5222", "10.20.10.0:4222"

	nb.Init(broker.Addrs(addr1, addr2))

	if len(nb.Options().Addrs) != 2 {
		t.Errorf("Expected Addr count = 2, Actual Addr count = %d", len(nb.Options().Addrs))
	}

	expectedAddr := fmt.Sprintf("nats://%s,nats://%s", addr1, addr2)

	if nb.Address() != expectedAddr {
		t.Errorf("Expected = '%s', Actual = '%s'", expectedAddr, nb.Address())
	}
}
