package models

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func GetPlaceTest(t *testing.T) {
	cases := []struct {
		name     string
		plan     Plan
		placeId  string
		expected *Place
	}{
		{
			name: "Valid place ID",
			plan: Plan{
				Places: []Place{
					{
						Id:   "1",
						Name: "place1",
					},
				},
			},
			placeId: "1",
			expected: &Place{
				Id:   "1",
				Name: "place1",
			},
		},
		{
			name: "Invalid place ID",
			plan: Plan{
				Places: []Place{
					{
						Id:   "1",
						Name: "place1",
					},
				},
			},
			placeId:  "2",
			expected: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.plan.GetPlace(c.placeId)
			if actual != c.expected {
				t.Errorf("expected: %v, actual: %v", c.expected, actual)
			}
		})
	}
}

func TestRecreateTransition(t *testing.T) {
	cases := []struct {
		name     string
		start    *GeoLocation
		plan     Plan
		expected []Transition
	}{
		{
			name: "現在位置から作成されたプランが順序入れ替えされた後の移動情報の再構成",
			start: &GeoLocation{
				Latitude:  35.1706431,
				Longitude: 136.8816945,
			},
			plan: Plan{
				Id:   "A",
				Name: "「名古屋駅」-「名古屋市科学館」-「名古屋市科学館」",
				Places: []Place{
					{
						Id:   "02",
						Name: "名古屋市科学館",
						Location: GeoLocation{
							Latitude:  35.165077,
							Longitude: 136.899703,
						},
					},
					{
						Id:   "01",
						Name: "名古屋市博物館",
						Location: GeoLocation{
							Latitude:  35.163926,
							Longitude: 136.901071,
						},
					},
				},
				Transitions: []Transition{
					{
						FromPlaceId: nil,
						ToPlaceId:   "01",
						Duration:    30,
					},
					{
						FromPlaceId: toStrPointer("01"),
						ToPlaceId:   "02",
						Duration:    5,
					},
				},
			},
			expected: []Transition{
				{
					FromPlaceId: nil,
					ToPlaceId:   "02",
					Duration:    21,
				},
				{
					FromPlaceId: toStrPointer("02"),
					ToPlaceId:   "01",
					Duration:    2,
				},
			},
		},
		{
			name: "指定された場所から作成されたプランが順序入れ替えされた後の移動情報の再構成",
			plan: Plan{
				Id:   "B",
				Name: "「東京タワー」-「東京スカイツリー」",
				Places: []Place{
					{
						Id:   "02",
						Name: "東京タワー",
						Location: GeoLocation{
							Latitude:  35.658581,
							Longitude: 139.745433,
						},
					},
					{
						Id:   "01",
						Name: "東京スカイツリー",
						Location: GeoLocation{
							Latitude:  35.710063,
							Longitude: 139.8107,
						},
					},
				},
				Transitions: []Transition{
					{
						FromPlaceId: toStrPointer("01"),
						ToPlaceId:   "02",
						Duration:    30,
					},
				},
				TimeInMinutes: 100,
			},
			expected: []Transition{
				{
					FromPlaceId: toStrPointer("02"),
					ToPlaceId:   "01",
					Duration:    102,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := c.plan.RecreateTransition(c.start)

			if diff := cmp.Diff(c.expected, result); diff != "" {
				t.Errorf("expected %v, but got %v", c.expected, result)
			}
		})
	}
}
