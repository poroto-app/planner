package models

import (
	"reflect"
	"testing"
)

func TestGetCategoryOfName(t *testing.T) {
	cases := []struct {
		name         string
		categoryName string
		expected     LocationCategory
	}{
		{
			name:         "amusements is name of CategoryAmusements",
			categoryName: "amusements",
			expected:     CategoryAmusements,
		},
		{
			name:         "restaurant is name of CategoryRestaurant",
			categoryName: "restaurant",
			expected:     CategoryRestaurant,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resultCategory := GetCategoryOfName(c.categoryName)
			if !reflect.DeepEqual(resultCategory, c.expected) {
				t.Errorf("expected: %v\nactual: %v", resultCategory, c.expected)
			}
		})
	}
}
