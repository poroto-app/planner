package array

import "testing"

func TestFind(t *testing.T) {
	cases := []struct {
		name      string
		slice     []int
		condition func(int) bool
		expected  int
	}{
		{
			name:  "Find element in slice",
			slice: []int{1, 2, 3, 4, 5},
			condition: func(i int) bool {
				return i == 3
			},
			expected: 3,
		},
		{
			name:  "Find element not in slice",
			slice: []int{1, 2, 3, 4, 5},
			condition: func(i int) bool {
				return i == 6
			},
		},
		{
			name:  "Find element in empty slice",
			slice: []int{},
			condition: func(i int) bool {
				return i == 6
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			actual, ok := Find(c.slice, c.condition)
			if !ok {
				if c.expected != 0 {
					t.Errorf("expected %d, but actual %d", c.expected, actual)
				}
			} else {
				if actual != c.expected {
					t.Errorf("expected %d, but actual %d", c.expected, actual)
				}
			}
		})
	}
}
