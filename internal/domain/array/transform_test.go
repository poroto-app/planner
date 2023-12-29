package array

import (
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	cases := []struct {
		name      string
		slice     []int
		transform func(int) string
		expected  []string
	}{
		{
			name:      "empty slice",
			transform: func(i int) string { return strconv.Itoa(i) },
			slice:     []int{},
			expected:  []string{},
		},
		{
			name:      "one element",
			transform: func(i int) string { return strconv.Itoa(i) },
			slice:     []int{1},
			expected:  []string{"1"},
		},
		{
			name:      "multiple elements",
			transform: func(i int) string { return strconv.Itoa(i) },
			slice:     []int{1, 2, 3},
			expected:  []string{"1", "2", "3"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := Map(c.slice, c.transform)
			if len(actual) != len(c.expected) {
				t.Errorf("expected: %v, actual: %v", c.expected, actual)
			}
			for i, v := range actual {
				if v != c.expected[i] {
					t.Errorf("expected: %v, actual: %v", c.expected, actual)
				}
			}
		})
	}
}
