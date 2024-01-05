package array

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
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

func TestMapAndFilter(t *testing.T) {
	cases := []struct {
		name      string
		slice     []int
		transform func(int) (string, bool)
		expected  []string
	}{
		{
			name:      "empty slice",
			transform: func(i int) (string, bool) { return strconv.Itoa(i), true },
			slice:     []int{},
			expected:  []string{},
		},
		{
			name:      "no match",
			transform: func(i int) (string, bool) { return strconv.Itoa(i), i < 0 },
			slice:     []int{1, 2, 3},
			expected:  []string{},
		},
		{
			name:      "one match",
			transform: func(i int) (string, bool) { return strconv.Itoa(i), i == 2 },
			slice:     []int{1, 2, 3},
			expected:  []string{"2"},
		},
		{
			name:      "multiple matches",
			transform: func(i int) (string, bool) { return strconv.Itoa(i), i > 1 },
			slice:     []int{1, 2, 3},
			expected:  []string{"2", "3"},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			actual := MapAndFilter(c.slice, c.transform)
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

func TestMapWithErr(t *testing.T) {
	cases := []struct {
		name      string
		slice     []int
		transform func(int) (*string, error)
		expected  *[]string
	}{
		{
			name:      "empty slice",
			transform: func(i int) (*string, error) { return &[]string{strconv.Itoa(i)}[0], nil },
			slice:     []int{},
			expected:  &[]string{},
		},
		{
			name:      "one element",
			transform: func(i int) (*string, error) { return &[]string{strconv.Itoa(i)}[0], nil },
			slice:     []int{1},
			expected:  &[]string{"1"},
		},
		{
			name:      "multiple elements",
			transform: func(i int) (*string, error) { return &[]string{strconv.Itoa(i)}[0], nil },
			slice:     []int{1, 2, 3},
			expected:  &[]string{"1", "2", "3"},
		},
		{
			name:      "return nil",
			transform: func(i int) (*string, error) { return nil, nil },
			slice:     []int{1, 2, 3},
			expected:  &[]string{},
		},
		{
			name: "return error",
			transform: func(i int) (*string, error) {
				if i == 2 {
					return nil, fmt.Errorf("error")
				}
				return &[]string{strconv.Itoa(i)}[0], nil
			},
			slice:    []int{1, 2, 3},
			expected: nil,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			actual, err := MapWithErr(c.slice, c.transform)
			if err != nil {
				if c.expected != nil {
					t.Errorf("expected: %v, actual: %v", c.expected, actual)
				}
				return
			}
			if c.expected == nil && actual != nil {
				t.Errorf("expected: %v, actual: %v", c.expected, actual)
			}

			if diff := cmp.Diff(*c.expected, *actual); diff != "" {
				t.Errorf("expected: %v, actual: %v", c.expected, actual)
			}
		})
	}
}
func TestFlatten(t *testing.T) {
	cases := []struct {
		name     string
		slice    [][]int
		expected []int
	}{
		{
			name:     "empty slice",
			slice:    [][]int{},
			expected: []int{},
		},
		{
			name:     "one element",
			slice:    [][]int{{1}},
			expected: []int{1},
		},
		{
			name:     "multiple elements",
			slice:    [][]int{{1}, {2, 3}},
			expected: []int{1, 2, 3},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			actual := Flatten(c.slice)
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
