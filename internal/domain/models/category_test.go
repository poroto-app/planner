package models

import (
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func TestLocationCategory_IsCategoryOf(t *testing.T) {
	cases := []struct {
		name       string
		category   LocationCategory
		categories []LocationCategory
		expected   bool
	}{
		{
			name:       "CategoryRestaurant is CategoryRestaurant",
			category:   CategoryRestaurant,
			categories: []LocationCategory{CategoryRestaurant},
			expected:   true,
		},
		{
			name:       "CategoryRestaurant is not CategoryCafe",
			category:   CategoryRestaurant,
			categories: []LocationCategory{CategoryCafe},
			expected:   false,
		},
		{
			name:       "CategoryRestaurant is CategoryRestaurant and CategoryCafe",
			category:   CategoryRestaurant,
			categories: []LocationCategory{CategoryRestaurant, CategoryCafe},
			expected:   true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			result := c.category.IsCategoryOf(c.categories...)
			if result != c.expected {
				t.Errorf("expected: %v\nactual: %v", result, c.expected)
			}
		})
	}
}

func TestGetCategoryOfName(t *testing.T) {
	cases := []struct {
		name         string
		categoryName string
		expected     *LocationCategory
	}{
		{
			name:         "The category of amusements is CategoryAmusements",
			categoryName: "amusements",
			expected:     &CategoryAmusements,
		},
		{
			name:         "The category of restaurant is CategoryRestaurant",
			categoryName: "restaurant",
			expected:     &CategoryRestaurant,
		},
		{
			name:         "The category does not exist",
			categoryName: "not_exist",
			expected:     nil,
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

func TestGetCategoriesFromSubCategories(t *testing.T) {
	cases := []struct {
		name          string
		subCategories []string
		expected      []LocationCategory
	}{
		{
			name:          "categories belong to CategoryCafe and CategoryShopping",
			subCategories: []string{"cafe", "shoe_store"},
			expected:      []LocationCategory{CategoryCafe, CategoryShopping},
		},
		{
			name:          "should not return duplicated categories",
			subCategories: []string{"cafe", "cafe"},
			expected:      []LocationCategory{CategoryCafe},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resultCategories := GetCategoriesFromSubCategories(c.subCategories)
			if diff := cmp.Diff(resultCategories, c.expected); diff != "" {
				t.Errorf("resultCategories differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestGetCategoryToFilter(t *testing.T) {
	t.Run("should have DefaultPhoto and EstimatedStayDuration", func(t *testing.T) {
		categories := GetCategoryToFilter()
		for _, category := range categories {
			if category.DefaultPhoto == "" {
				t.Errorf("category %v should have DefaultPhoto", category.Name)
			}
			if category.EstimatedStayDuration == 0 {
				t.Errorf("category %v should have EstimatedStayDuration", category.Name)
			}
		}
	})
}
