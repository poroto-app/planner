package plangenerator

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"poroto.app/poroto/planner/internal/domain/models"
)

func TestAddTransition(t *testing.T) {
	cases := []struct {
		name                         string
		placesInPlan                 []models.Place
		transitions                  []models.Transition
		duration                     uint
		createBasedOnCurrentLocation bool
		expected                     []models.Transition
	}{
		{
			name:                         "placesInPlanが空の場合は何もしない",
			placesInPlan:                 []models.Place{},
			transitions:                  []models.Transition{},
			duration:                     10,
			createBasedOnCurrentLocation: true,
			expected:                     []models.Transition{},
		},
		{

			name:                         "placesInPlanが空の場合は何もしない",
			placesInPlan:                 []models.Place{},
			transitions:                  []models.Transition{},
			duration:                     10,
			createBasedOnCurrentLocation: false,
			expected:                     []models.Transition{},
		},
		{
			name: "現在地から作成したプランの場合は、出発値の情報がnilになる",
			placesInPlan: []models.Place{
				{
					Id: "1",
				},
			},
			transitions:                  []models.Transition{},
			duration:                     10,
			createBasedOnCurrentLocation: true,
			expected: []models.Transition{
				{
					FromPlaceId: nil,
					ToPlaceId:   "1",
					Duration:    10,
				},
			},
		},
		{
			name: "現在地から作成したプランの場合は、2箇所目以降の場所では出発値の情報が前の場所のIDになる",
			placesInPlan: []models.Place{
				{
					Id: "1",
				},
				{
					Id: "2",
				},
			},
			transitions: []models.Transition{
				{
					FromPlaceId: nil,
					ToPlaceId:   "1",
					Duration:    10,
				},
			},
			duration:                     10,
			createBasedOnCurrentLocation: true,
			expected: []models.Transition{
				{
					FromPlaceId: nil,
					ToPlaceId:   "1",
					Duration:    10,
				},
				{
					FromPlaceId: toStringPointer("1"),
					ToPlaceId:   "2",
					Duration:    10,
				},
			},
		},
		{
			name: "場所を指定して作成した場合、プラン内の場所が１箇所だけの場合は移動が存在しない",
			placesInPlan: []models.Place{
				{
					Id: "1",
				},
			},
			transitions:                  []models.Transition{},
			duration:                     10,
			createBasedOnCurrentLocation: false,
			expected:                     []models.Transition{},
		},
		{
			name: "場所を指定して作成した場合、指定した場所を出発地点とする",
			placesInPlan: []models.Place{
				{
					Id: "1",
				},
				{
					Id: "2",
				},
			},
			transitions:                  []models.Transition{},
			duration:                     10,
			createBasedOnCurrentLocation: false,
			expected: []models.Transition{
				{
					FromPlaceId: toStringPointer("1"),
					ToPlaceId:   "2",
					Duration:    10,
				},
			},
		},
	}

	for _, c := range cases {
		s := Service{}
		t.Run(c.name, func(t *testing.T) {
			result := s.AddTransition(
				c.placesInPlan,
				c.transitions,
				c.duration,
				c.createBasedOnCurrentLocation,
			)
			if diff := cmp.Diff(c.expected, result); diff != "" {
				t.Errorf("expected %v, but got %v", c.expected, result)
			}
		})
	}
}

func toStringPointer(value string) *string {
	return &value
}
