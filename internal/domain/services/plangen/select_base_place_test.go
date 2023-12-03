package plangen

import (
	"github.com/google/go-cmp/cmp"
	"poroto.app/poroto/planner/internal/domain/models"
	"testing"
)

func TestService_SelectBasePlace(t *testing.T) {
	cases := []struct {
		name     string
		input    SelectBasePlaceInput
		expected []models.Place
	}{
		{
			name: "should remove places that are far from the base location",
			input: SelectBasePlaceInput{
				BaseLocation: models.GeoLocation{
					// 新宿駅
					Latitude:  35.690817373071,
					Longitude: 139.7065625287,
				},
				Places: []models.Place{
					{
						Id:       "takashimaya",
						Name:     "新宿高島屋",
						Location: models.GeoLocation{Latitude: 35.6875312, Longitude: 139.7022521},
					},
					{
						Id:       "shinjuku-gyoen",
						Name:     "ホテル雅叙園東京",
						Location: models.GeoLocation{Latitude: 35.6305774, Longitude: 139.7142515},
					},
				},
			},
			expected: []models.Place{
				{
					Id:       "takashimaya",
					Name:     "新宿高島屋",
					Location: models.GeoLocation{Latitude: 35.6875312, Longitude: 139.7022521},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := Service{}
			actual := s.SelectBasePlace(c.input)
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Errorf("(-want +got)\n%s", diff)
			}
		})
	}
}

func TestIsAlreadyAdded(t *testing.T) {
	cases := []struct {
		name     string
		place    models.Place
		places   []models.Place
		expected bool
	}{
		{
			name:  "should return true when place is already added",
			place: models.Place{Id: "1"},
			places: []models.Place{
				{Id: "1"},
				{Id: "2"},
			},
			expected: true,
		},
		{
			name:  "should return false when place is not added",
			place: models.Place{Id: "3"},
			places: []models.Place{
				{Id: "1"},
				{Id: "2"},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := isAlreadyAdded(c.place, c.places)
			if actual != c.expected {
				t.Errorf("expected: %v, actual: %v", c.expected, actual)
			}
		})
	}
}
