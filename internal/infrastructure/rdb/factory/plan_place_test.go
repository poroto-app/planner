package factory

import (
	"github.com/google/go-cmp/cmp"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
	"testing"
)

func TestNewPlanPlaceSliceFromDomainModel(t *testing.T) {
	cases := []struct {
		name       string
		planPlaces []models.Place
		planId     string
		expected   []generated.PlanPlace
	}{
		{
			name: "should return a valid slice",
			planPlaces: []models.Place{
				{
					Id: "ec7c607d-454a-4644-929a-c3b1e078842d",
				},
				{
					Id: "339809cf-d515-4a64-bbcd-c6a899051273",
				},
			},
			expected: []generated.PlanPlace{
				{
					PlaceID:   "ec7c607d-454a-4644-929a-c3b1e078842d",
					SortOrder: 0,
				},
				{
					PlaceID:   "339809cf-d515-4a64-bbcd-c6a899051273",
					SortOrder: 1,
				},
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			actual := NewPlanPlaceSliceFromDomainMode(c.planPlaces, c.planId)
			if len(actual) != len(c.expected) {
				t.Errorf("wrong plan place slice length, want: %d, got: %d", len(c.expected), len(actual))
			}

			for i, expected := range c.expected {
				if actual[i].PlaceID != expected.PlaceID {
					t.Errorf("wrong place id, want: %s, got: %s", expected.PlaceID, actual[i].PlaceID)
				}
				if actual[i].SortOrder != expected.SortOrder {
					t.Errorf("wrong sort order, want: %d, got: %d", expected.SortOrder, actual[i].SortOrder)
				}
			}
		})
	}
}

func TestNewPlacesFromEntities(t *testing.T) {
	cases := []struct {
		name           string
		planPlaceSlice generated.PlanPlaceSlice
		places         []models.Place
		planId         string
		expected       []models.Place
	}{
		{
			name:   "should return a valid slice",
			planId: "4a81310e-03c9-4862-8a13-d3c75475bc6e",
			planPlaceSlice: generated.PlanPlaceSlice{
				{
					PlaceID:   "ec7c607d-454a-4644-929a-c3b1e078842d",
					PlanID:    "4a81310e-03c9-4862-8a13-d3c75475bc6e",
					SortOrder: 1,
				},
				{
					PlaceID:   "339809cf-d515-4a64-bbcd-c6a899051273",
					PlanID:    "4a81310e-03c9-4862-8a13-d3c75475bc6e",
					SortOrder: 0,
				},
				{
					// 異なるプラン
					PlaceID:   "d0a3e6b7-11c6-4050-8d3f-68220349a8a7",
					PlanID:    "d0a3e6b7-11c6-4050-8d3f-68220349a8a7",
					SortOrder: 0,
				},
			},
			places: []models.Place{
				{Id: "339809cf-d515-4a64-bbcd-c6a899051273"},
				{Id: "ec7c607d-454a-4644-929a-c3b1e078842d"},
			},
			expected: []models.Place{
				{Id: "339809cf-d515-4a64-bbcd-c6a899051273"},
				{Id: "ec7c607d-454a-4644-929a-c3b1e078842d"},
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewPlanPlacesFromEntities(c.planPlaceSlice, c.places, c.planId)
			if err != nil {
				t.Errorf("error converting entities to domain models: %v", err)
			}

			if diff := cmp.Diff(c.expected, *actual); diff != "" {
				t.Errorf("wrong plan places (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNewPlacesFromEntities_ShouldFail(t *testing.T) {
	cases := []struct {
		name           string
		planPlaceSlice generated.PlanPlaceSlice
		places         []models.Place
		planId         string
	}{
		{
			name:   "should fail if place not found",
			planId: "4a81310e-03c9-4862-8a13-d3c75475bc6e",
			planPlaceSlice: generated.PlanPlaceSlice{
				{
					PlaceID:   "ec7c607d-454a-4644-929a-c3b1e078842d",
					PlanID:    "4a81310e-03c9-4862-8a13-d3c75475bc6e",
					SortOrder: 1,
				},
				{
					PlaceID:   "339809cf-d515-4a64-bbcd-c6a899051273",
					PlanID:    "4a81310e-03c9-4862-8a13-d3c75475bc6e",
					SortOrder: 0,
				},
			},
			places: []models.Place{
				{Id: "339809cf-d515-4a64-bbcd-c6a899051273"},
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewPlanPlacesFromEntities(c.planPlaceSlice, c.places, c.planId)
			if err == nil {
				t.Errorf("error should be returned")
			}
		})
	}
}
