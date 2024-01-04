package array

import "testing"

func TestFilter(t *testing.T) {
	cases := []struct {
		name      string
		slice     []int
		condition func(int) bool
		expected  []int
	}{
		{
			name:      "empty slice",
			condition: func(i int) bool { return i > 0 },
			slice:     []int{},
			expected:  []int{},
		},
		{
			name:      "no match",
			condition: func(i int) bool { return i < 0 },
			slice:     []int{1, 2, 3},
			expected:  []int{},
		},
		{
			name:      "one match",
			condition: func(i int) bool { return i == 2 },
			slice:     []int{1, 2, 3},
			expected:  []int{2},
		},
		{
			name:      "multiple matches",
			condition: func(i int) bool { return i > 1 },
			slice:     []int{1, 2, 3},
			expected:  []int{2, 3},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := Filter(c.slice, c.condition)
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

func TestDistinctBy(t *testing.T) {
	cases := []struct {
		name     string
		slice    []int
		selector func(int) int
		expected []int
	}{
		{
			name:     "empty slice",
			selector: func(i int) int { return i },
			slice:    []int{},
			expected: []int{},
		},
		{
			name:     "identical elements",
			selector: func(i int) int { return i },
			slice:    []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		{
			name:     "duplicated elements",
			selector: func(i int) int { return i },
			slice:    []int{1, 2, 1},
			expected: []int{1, 2},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := DistinctBy(c.slice, c.selector)
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
