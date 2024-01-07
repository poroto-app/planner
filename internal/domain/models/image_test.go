package models

import (
	"poroto.app/poroto/planner/internal/domain/utils"
	"testing"
)

func TestImageDefault(t *testing.T) {
	cases := []struct {
		name     string
		image    ImageSmallLarge
		expected string
	}{
		{
			name: "small image",
			image: ImageSmallLarge{
				Small: utils.StrPointer("small"),
			},
			expected: "small",
		},
		{
			name: "large image",
			image: ImageSmallLarge{
				Large: utils.StrPointer("large"),
			},
			expected: "large",
		},
		{
			name: "both small and large image",
			image: ImageSmallLarge{
				Small: utils.StrPointer("small"),
				Large: utils.StrPointer("large"),
			},
			expected: "large",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.image.Default()
			if actual != c.expected {
				t.Errorf("expected: %s, actual: %s", c.expected, actual)
			}
		})
	}
}
