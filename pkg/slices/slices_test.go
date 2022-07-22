package slices_test

import (
	"fmt"
	"testing"

	"github.com/kristofferostlund/gpx-to-gmaps/pkg/slices"
)

func TestPickSpaced(t *testing.T) {
	for _, test := range []struct {
		items         []int
		maxCount      int
		expectedCount int
	}{
		{
			items:         fill(100),
			maxCount:      10,
			expectedCount: 10,
		},
		{
			items:         fill(123),
			maxCount:      10,
			expectedCount: 10,
		},
	} {
		t.Run(fmt.Sprintf("%d items, max count %d returns %d", len(test.items), test.maxCount, test.expectedCount), func(t *testing.T) {
			actual := slices.PickSpaced(test.items, test.maxCount)
			if got, want := len(actual), test.expectedCount; got != want {
				t.Errorf("got count %d, want %d (%v)", got, want, actual)
			}
		})
	}
}

func fill(num int) []int {
	out := make([]int, 0, num)
	for i := 0; i < num; i++ {
		out = append(out, i)
	}
	return out
}
