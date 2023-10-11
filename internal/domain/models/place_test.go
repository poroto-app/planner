package models

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestPlace_MainCategory(t *testing.T) {
	cases := []struct {
		name     string
		place    Place
		expected *LocationCategory
	}{
		{
			name: "place has no category",
			place: Place{
				Categories: []LocationCategory{},
			},
			expected: nil,
		},
		{
			name: "place has one category",
			place: Place{
				Categories: []LocationCategory{
					CategoryAmusements,
				},
			},
			expected: &CategoryAmusements,
		},
		{
			name: "place has two categories and the first one is main category",
			place: Place{
				Categories: []LocationCategory{
					CategoryAmusements,
					CategoryBookStore,
				},
			},
			expected: &CategoryAmusements,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.place.MainCategory()
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Errorf("MainCategory() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPlace_EstimatedStayDuration(t *testing.T) {
	cases := []struct {
		name     string
		place    Place
		expected uint
	}{
		{
			name: "place has no category",
			place: Place{
				Categories: []LocationCategory{},
			},
			expected: 0,
		},
		{
			name: "place has one category",
			place: Place{
				Categories: []LocationCategory{
					{EstimatedStayDuration: 10},
				},
			},
			expected: 10,
		},
		{
			name: "place has two categories and return the estimated stay duration of the first one",
			place: Place{
				Categories: []LocationCategory{
					{EstimatedStayDuration: 10},
					{EstimatedStayDuration: 20},
				},
			},
			expected: 10,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.place.EstimatedStayDuration()
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Errorf("EstimatedStayDuration() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
