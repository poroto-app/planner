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
			name:         "The category of amusements is CategoryAmusements",
			categoryName: "amusements",
			expected:     CategoryAmusements,
		},
		{
			name:         "The category of restaurant is CategoryRestaurant",
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

func TestCategoryOfSubCategory(t *testing.T) {
	cases := []struct {
		name        string
		subCategory string
		expected    LocationCategory
	}{

		{
			name:        "cafe belongs to CategoryCafe",
			subCategory: "cafe",
			expected:    CategoryCafe,
		},
		{
			name:        "shoe_store belongs to CategoryShopping",
			subCategory: "shoe_store",
			expected:    CategoryShopping,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resultCategory := CategoryOfSubCategory(c.subCategory)
			if !reflect.DeepEqual(*resultCategory, c.expected) {
				t.Errorf("expected: %v\nactual: %v", *resultCategory, c.expected)
			}
		})
	}
}
