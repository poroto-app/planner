package models

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestGetPlace(t *testing.T) {
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
						Id: "1",
					},
				},
			},
			placeId: "1",
			expected: &Place{
				Id: "1",
			},
		},
		{
			name: "Invalid place ID",
			plan: Plan{
				Places: []Place{
					{
						Id: "1",
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
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}

func TestPlan_PlacesReorderedToMinimizeDistance(t *testing.T) {
	cases := []struct {
		name     string
		plan     Plan
		expected []Place
	}{
		{
			name: "should start with first place",
			plan: Plan{
				Places: []Place{
					NewMockPlaceShinjukuStation(),
				},
			},
			expected: []Place{
				NewMockPlaceShinjukuStation(),
			},
		},
		{
			name: "should return places reordered to minimize distance from first place",
			plan: Plan{
				Places: []Place{
					// 新宿駅
					NewMockPlaceShinjukuStation(),
					// 伊勢丹
					NewMockPlaceIsetan(),
					// 新宿御苑
					NewMockPlaceShinjukuGyoen(),
					// 高島屋
					NewMockPlaceTakashimaya(),
				},
			},
			expected: []Place{
				// 新宿駅
				NewMockPlaceShinjukuStation(),
				// 高島屋
				NewMockPlaceTakashimaya(),
				// 伊勢丹
				NewMockPlaceIsetan(),
				// 新宿御苑
				NewMockPlaceShinjukuGyoen(),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.plan.PlacesReorderedToMinimizeDistance()
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}
