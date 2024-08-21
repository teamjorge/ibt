package utilities

import (
	"sort"
	"testing"
)

func TestDistinct(t *testing.T) {
	t.Run("test distinct normal", func(t *testing.T) {
		items := []int{5, 1, 5, 5, 3, 1, 4}

		distinctItems := GetDistinct(items)
		if len(distinctItems) != 4 {
			t.Errorf("expected items to have len %d. received %d", 4, len(distinctItems))
		}

		sort.Ints(distinctItems)

		if distinctItems[0] != 1 && distinctItems[1] != 3 && distinctItems[2] != 4 && distinctItems[3] != 5 {
			t.Errorf("expected items to be %v. received: %v", []int{1, 3, 4, 5}, distinctItems)
		}
	})

	t.Run("test distinct empty", func(t *testing.T) {
		var items []int = nil

		distinctItems := GetDistinct(items)
		if len(distinctItems) != 0 {
			t.Errorf("expected items to have len %d. received %d", 0, len(distinctItems))
		}
	})
}
