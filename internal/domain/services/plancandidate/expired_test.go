package plancandidate

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"poroto.app/poroto/planner/internal/infrastructure/mock"
)

func TestDeleteExpiredPlanCandidates(t *testing.T) {
	cases := []struct {
		name                       string
		expiresAt                  time.Time
		planCandidates             map[string]models.PlanCandidate
		placeSearchResults         map[string][]places.Place
		expectedPlanCandidates     map[string]models.PlanCandidate
		expectedPlaceSearchResults map[string][]places.Place
	}{
		{
			name:      "expired plan candidates are deleted",
			expiresAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			planCandidates: map[string]models.PlanCandidate{
				"planCandidate1": {
					Id:        "planCandidate1",
					ExpiresAt: time.Date(2019, 12, 31, 23, 59, 59, 0, time.UTC),
				},
				"planCandidate2": {
					Id:        "planCandidate2",
					ExpiresAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				"planCandidate3": {
					Id:        "planCandidate3",
					ExpiresAt: time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
				},
			},
			placeSearchResults: map[string][]places.Place{
				"planCandidate1": {{PlaceID: "place1"}},
				"planCandidate2": {{PlaceID: "place2"}},
				"planCandidate3": {{PlaceID: "place3"}},
			},
			expectedPlanCandidates: map[string]models.PlanCandidate{
				"planCandidate2": {
					Id:        "planCandidate2",
					ExpiresAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				"planCandidate3": {
					Id:        "planCandidate3",
					ExpiresAt: time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
				},
			},
			expectedPlaceSearchResults: map[string][]places.Place{
				"planCandidate2": {{PlaceID: "place2"}},
				"planCandidate3": {{PlaceID: "place3"}},
			},
		},
	}

	for _, c := range cases {
		planCandidateRepository := mock.NewPlanCandidateRepository(c.planCandidates)
		planSearchResultRepository := mock.NewPlaceSearchResultRepository(c.placeSearchResults)

		service := Service{
			planCandidateRepository:     planCandidateRepository,
			placeSearchResultRepository: planSearchResultRepository,
		}

		err := service.DeleteExpiredPlanCandidates(context.Background(), c.expiresAt)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if diff := cmp.Diff(c.expectedPlanCandidates, c.planCandidates); diff != "" {
			t.Errorf("unexpected plan candidates (-want +got):\n%s", diff)
		}

		if diff := cmp.Diff(c.expectedPlaceSearchResults, c.placeSearchResults); diff != "" {
			t.Errorf("unexpected place search results (-want +got):\n%s", diff)
		}
	}
}
