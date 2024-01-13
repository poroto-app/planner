package models

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestPlace_IsSameCategoryPlace(t *testing.T) {
	cases := []struct {
		name     string
		a        Place
		b        Place
		expected bool
	}{
		{
			name: "should return true when two places are same category",
			a: Place{
				Google: GooglePlace{
					Types: []string{CategoryRestaurant.SubCategories[0]},
				},
			},
			b: Place{
				Google: GooglePlace{
					Types: []string{CategoryRestaurant.SubCategories[1]},
				},
			},
			expected: true,
		},
		{
			name: "should return false when two places are not same category",
			a: Place{
				Google: GooglePlace{
					Types: []string{CategoryRestaurant.SubCategories[0]},
				},
			},
			b: Place{
				Google: GooglePlace{
					Types: []string{CategoryAmusements.SubCategories[0]},
				},
			},
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.a.IsSameCategoryPlace(c.b)
			if actual != c.expected {
				t.Errorf("expected: %v, actual: %v", c.expected, actual)
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
				Google: GooglePlace{
					Types: []string{},
				},
			},
			expected: 0,
		},
		{
			name: "place has one category",
			place: Place{
				Google: GooglePlace{
					Types: []string{CategoryAmusements.SubCategories[0]},
				},
			},
			expected: CategoryAmusements.EstimatedStayDuration,
		},
		{
			name: "place has two categories and return the estimated stay duration of the first one",
			place: Place{
				Google: GooglePlace{
					Types: []string{
						CategoryRestaurant.SubCategories[0],
						CategoryAmusements.SubCategories[1],
					},
				},
			},
			expected: CategoryRestaurant.EstimatedStayDuration,
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

func TestShufflePlaces(t *testing.T) {
	cases := []struct {
		name     string
		places   []Place
		expected []Place
	}{
		{
			name: "should return shuffled places",
			places: []Place{
				NewMockPlaceShinjukuStation(),
				NewMockPlaceIsetan(),
				NewMockPlaceShinjukuGyoen(),
				NewMockPlaceTakashimaya(),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			original := make([]Place, len(c.places))
			copy(original, c.places)

			actual := ShufflePlaces(c.places)
			if diff := cmp.Diff(len(c.places), len(actual)); diff != "" {
				t.Errorf("ShufflePlaces() mismatch (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(original, actual); diff == "" {
				t.Errorf("ShufflePlaces() should return shuffled places")
			}

			if diff := cmp.Diff(original, c.places); diff != "" {
				t.Errorf("ShufflePlaces() should not modify original places")
			}
		})
	}
}

func TestSortPlacesByRating(t *testing.T) {
	cases := []struct {
		name     string
		places   []Place
		expected []Place
	}{
		{
			name: "should return sorted places by rating",
			places: []Place{
				NewMockPlaceShinjukuStation(),
				NewMockPlaceIsetan(),
				NewMockPlaceShinjukuGyoen(),
				NewMockPlaceTakashimaya(),
			},
			expected: []Place{
				NewMockPlaceTakashimaya(),
				NewMockPlaceShinjukuGyoen(),
				NewMockPlaceIsetan(),
				NewMockPlaceShinjukuStation(),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := SortPlacesByRating(c.places)
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Errorf("SortPlacesByRating() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// ==============================================================
// Mocks
// ==============================================================
func NewMockPlaceShinjukuStation() Place {
	return Place{
		Id:   "sinjuku-station",
		Name: "新宿駅",
		Location: GeoLocation{
			Latitude:  35.6899573,
			Longitude: 139.7005071,
		},
		Google: GooglePlace{
			Rating: 4.5,
		},
	}
}

func NewMockPlaceIsetan() Place {
	return Place{
		Id:   "isetan",
		Name: "伊勢丹 新宿店",
		Location: GeoLocation{
			Latitude:  35.6916532,
			Longitude: 139.7046449,
		},
		Google: GooglePlace{
			Rating: 4.6,
		},
	}
}

func NewMockPlaceShinjukuGyoen() Place {
	return Place{
		Id:   "shinjuku-gyoen",
		Name: "新宿御苑",
		Location: GeoLocation{
			Latitude:  35.6867668,
			Longitude: 139.7123842,
		},
		Google: GooglePlace{
			Rating: 4.7,
		},
	}
}

func NewMockPlaceTakashimaya() Place {
	return Place{
		Id:   "takashimaya",
		Name: "新宿高島屋",
		Location: GeoLocation{
			Latitude:  35.6875312,
			Longitude: 139.7022521,
		},
		Google: GooglePlace{
			Rating: 4.8,
		},
	}
}
