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

	natsBroker, ok := nb.(*nbroker)
	if !ok {
		t.Fatal("Expected broker to be of types *nbroker")
	}

	addr1f := fmt.Sprintf("nats://%s", addr1)
	addr2f := fmt.Sprintf("nats://%s", addr2)

	if natsBroker.addrs[0] != addr1f && natsBroker.addrs[1] != addr2f {
		expAddr, actAddr := fmt.Sprintf("%s,%s", addr1f, addr2f), fmt.Sprintf("%s,%s", natsBroker.addrs[0], natsBroker.addrs[1])
		t.Errorf("Expected = '%s', Actual = '%s'", expAddr, actAddr)
	}

}
