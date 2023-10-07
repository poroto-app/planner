package models

import (
	"github.com/google/go-cmp/cmp"
	"poroto.app/poroto/planner/internal/domain/utils"
	"testing"
)

func TestCreateTransition(t *testing.T) {
	cases := []struct {
		name     string
		start    *GeoLocation
		places   []Place
		expected []Transition
	}{
		{
			name: "現在位置から作成されたプランの移動情報の再構成",
			start: &GeoLocation{
				Latitude:  35.1706431,
				Longitude: 136.8816945,
			},
			places: []Place{
				{
					Id:   "01",
					Name: "名古屋市科学館",
					Location: GeoLocation{
						Latitude:  35.165077,
						Longitude: 136.899703,
					},
				},
				{
					Id:   "02",
					Name: "名古屋市博物館",
					Location: GeoLocation{
						Latitude:  35.163926,
						Longitude: 136.901071,
					},
				},
			},
			expected: []Transition{
				{
					FromPlaceId: nil,
					ToPlaceId:   "01",
					Duration:    21,
				},
				{
					FromPlaceId: utils.StrPointer("01"),
					ToPlaceId:   "02",
					Duration:    2,
				},
			},
		},
		{
			name: "指定された場所から作成されたプランの移動情報の再構成",
			places: []Place{
				{
					Id:   "01",
					Name: "東京タワー",
					Location: GeoLocation{
						Latitude:  35.658581,
						Longitude: 139.745433,
					},
				},
				{
					Id:   "02",
					Name: "東京スカイツリー",
					Location: GeoLocation{
						Latitude:  35.710063,
						Longitude: 139.8107,
					},
				},
			},
			expected: []Transition{
				{
					FromPlaceId: utils.StrPointer("01"),
					ToPlaceId:   "02",
					Duration:    102,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := CreateTransition(c.places, c.start)

			if diff := cmp.Diff(c.expected, result); diff != "" {
				t.Errorf("expected %v, but got %v", c.expected, result)
			}
		})
	}
}
