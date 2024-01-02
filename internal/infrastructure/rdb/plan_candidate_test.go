package rdb

import (
	"context"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
	"testing"
	"time"
)

func TestPlanCandidateRepository_Create(t *testing.T) {
	cases := []struct {
		name            string
		planCandidateId string
		expiresAt       time.Time
	}{
		{
			name:            "success",
			planCandidateId: uuid.New().String(),
			expiresAt:       time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	planCandidateRepository, err := NewPlanCandidateRepository(testDB)
	if err != nil {
		t.Fatalf("failed to create plan candidate repository: %v", err)
	}

	testContext := context.Background()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := cleanup(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to cleanup: %v", err)
				}
			})

			if err := planCandidateRepository.Create(testContext, c.planCandidateId, c.expiresAt); err != nil {
				t.Fatalf("failed to create plan candidate: %v", err)
			}

			exists, err := entities.PlanCandidateSetExists(testContext, testDB, c.planCandidateId)
			if err != nil {
				t.Fatalf("failed to check plan candidate existence: %v", err)
			}

			if !exists {
				t.Fatalf("plan candidate does not exist")
			}

		})
	}
}

func TestPlanCandidateRepository_AddSearchedPlacesForPlanCandidate(t *testing.T) {
	cases := []struct {
		name            string
		planCandidateId string
		placeIds        []string
	}{
		{
			name:            "success",
			planCandidateId: uuid.New().String(),
			placeIds:        []string{uuid.New().String(), uuid.New().String()},
		},
	}

	planCandidateRepository, err := NewPlanCandidateRepository(testDB)
	if err != nil {
		t.Fatalf("failed to create plan candidate repository: %v", err)
	}

	testContext := context.Background()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := cleanup(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to cleanup: %v", err)
				}
			})

			// 事前にPlanCandidateSetを作成しておく
			if err := planCandidateRepository.Create(testContext, c.planCandidateId, time.Now().Add(time.Hour)); err != nil {
				t.Fatalf("failed to create plan candidate: %v", err)
			}

			// 事前にPlaceを作成しておく
			for _, placeId := range c.placeIds {
				placeEntity := entities.Place{ID: placeId}
				if err := placeEntity.Insert(testContext, testDB, boil.Infer()); err != nil {
					t.Fatalf("failed to insert place: %v", err)
				}
			}

			if err := planCandidateRepository.AddSearchedPlacesForPlanCandidate(testContext, c.planCandidateId, c.placeIds); err != nil {
				t.Fatalf("failed to add searched places for plan candidate: %v", err)
			}

			numPlanCandidateSetSearchedPlaces, err := entities.
				PlanCandidateSetSearchedPlaces(entities.PlanCandidateSetSearchedPlaceWhere.PlanCandidateSetID.EQ(c.planCandidateId)).
				Count(testContext, testDB)
			if err != nil {
				t.Fatalf("failed to get plan candidate places: %v", err)
			}

			if int(numPlanCandidateSetSearchedPlaces) != len(c.placeIds) {
				t.Fatalf("number of plan candidate places is not expected: %v", numPlanCandidateSetSearchedPlaces)
			}
		})
	}
}
