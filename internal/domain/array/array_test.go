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
func TestIsContain(t *testing.T) {
	cases := []struct {
		name     string
		array    []string
		value    string
		expected bool
	}{
		{
			name:     "is contain",
			array:    []string{"a", "b", "c"},
			value:    "a",
			expected: true,
		}, {
			name:     "is not contain",
			array:    []string{"a", "b", "c"},
			value:    "d",
			expected: false,
		}, {
			name:     "empty array",
			array:    []string{},
			value:    "a",
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := IsContain(c.array, c.value)
			if result != c.expected {
				t.Errorf("expected: %v, result: %v", c.expected, result)
			}
		})
	}
}
