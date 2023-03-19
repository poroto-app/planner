package array

import "testing"

func TestHasIntersection(t *testing.T) {
	cases := []struct {
		name     string
		a        []string
		b        []string
		expected bool
	}{
		{
			name:     "has intersection",
			a:        []string{"a", "b", "c"},
			b:        []string{"a", "d", "e"},
			expected: true,
		}, {
			name:     "not has intersection",
			a:        []string{"a", "b", "c"},
			b:        []string{"d", "e", "f"},
			expected: false,
		}, {
			name:     "empty values",
			a:        []string{},
			b:        []string{},
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := HasIntersection(c.a, c.b)
			if result != c.expected {
				t.Errorf("expected: %v, result: %v", c.expected, result)
			}
		})
	}
}
