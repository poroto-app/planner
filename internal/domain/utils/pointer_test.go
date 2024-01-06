package utils

import "testing"

func TestToPointer(t *testing.T) {
	cases := []struct {
		name     string
		input    int
		expected *int
	}{
		{
			name:     "success",
			input:    1,
			expected: ToPointer(1),
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			actual := ToPointer(c.input)
			if *actual != *c.expected {
				t.Errorf("expected: %v, actual: %v", c.expected, actual)
			}
		})
	}
}

func TestFromPointerOrZero(t *testing.T) {
	cases := []struct {
		name     string
		input    *int
		expected int
	}{
		{
			name:     "should return value",
			input:    ToPointer(1),
			expected: 1,
		},
		{
			name:     "should return zero",
			input:    nil,
			expected: 0,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			actual := FromPointerOrZero(c.input)
			if actual != c.expected {
				t.Errorf("expected: %v, actual: %v", c.expected, actual)
			}
		})
	}
}
