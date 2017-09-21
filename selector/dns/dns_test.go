package dns

import (
	"testing"
)

func TestDNSSelector(t *testing.T) {
	data := []string{"foo", "bar", "baz"}

	s := NewSelector()

	for _, name := range data {
		next, err := s.Select(name)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 3; i++ {
			node, err := next()
			if err != nil {
				t.Fatal(err)
			}
			if node.Address != name {
				t.Fatalf("got %s expected %s", node.Address, name)
			}
		}
	}
}
