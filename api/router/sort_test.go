package router

import (
	"sort"
	"testing"
)

func TestRouterSort(t *testing.T) {
	testData := []struct {
		Routes []Route
		Expect []int
	}{
		{
			Routes: []Route{
				{Priority: 5},
				{Priority: 2},
				{Priority: 3},
				{Priority: 0},
			},
			Expect: []int{0, 2, 3, 5},
		},
	}

	for _, d := range testData {
		r := Routes{d.Routes}
		sort.Sort(sortedRoutes{r})
		for i, j := range d.Expect {
			if r.Routes[i].Priority != j {
				t.Fatalf("Expected val %d got %d", j, r.Routes[i].Priority)
			}
		}
	}
}
