package rdb

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
)

func savePlaces(ctx context.Context, db *sql.DB, places []models.Place) error {
	places = array.DistinctBy(places, func(place models.Place) string { return place.Id })
	for _, place := range places {
		placeEntity := generated.Place{ID: place.Id}
		if err := placeEntity.Insert(ctx, db, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert place: %v", err)
		}

		if place.Google.PlaceId == "" {
			continue
		}

		googlePlaceEntity := generated.GooglePlace{GooglePlaceID: place.Google.PlaceId, PlaceID: place.Id}
		if err := placeEntity.AddGooglePlaces(ctx, db, true, &googlePlaceEntity); err != nil {
			return fmt.Errorf("failed to insert google place: %v", err)
		}
	}

	return nil
}

func savePlanCandidate(ctx context.Context, db *sql.DB, planCandidateSet models.PlanCandidate) error {
	// PlanCandidateSetを作成
	planCandidateSetEntity := generated.PlanCandidateSet{
		ID:        planCandidateSet.Id,
		ExpiresAt: planCandidateSet.ExpiresAt,
	}
	if err := planCandidateSetEntity.Insert(ctx, db, boil.Infer()); err != nil {
		return fmt.Errorf("failed to insert plan candidate set: %v", err)
	}

	// PlanCandidateSetMetaDataを作成
	if !planCandidateSet.MetaData.IsZero() {
		planCandidateSetMetaDataEntity := generated.PlanCandidateSetMetaDatum{
			ID:                           uuid.New().String(),
			PlanCandidateSetID:           planCandidateSet.Id,
			IsCreatedFromCurrentLocation: planCandidateSet.MetaData.CreatedBasedOnCurrentLocation,
			LatitudeStart:                planCandidateSet.MetaData.LocationStart.Latitude,
			LongitudeStart:               planCandidateSet.MetaData.LocationStart.Longitude,
			PlanDurationMinutes:          null.IntFromPtr(planCandidateSet.MetaData.FreeTime),
		}
		if err := planCandidateSetMetaDataEntity.Insert(ctx, db, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert plan candidate set meta data: %v", err)
		}

		if planCandidateSet.MetaData.CategoriesPreferred != nil {
			for _, category := range *planCandidateSet.MetaData.CategoriesPreferred {
				planCandidateSetCategoryEntity := generated.PlanCandidateSetMetaDataCategory{
					ID:                 uuid.New().String(),
					PlanCandidateSetID: planCandidateSet.Id,
					Category:           category.Name,
					IsSelected:         true,
				}
				if err := planCandidateSetCategoryEntity.Insert(ctx, db, boil.Infer()); err != nil {
					return fmt.Errorf("failed to insert plan candidate set category: %v", err)
				}
			}
		}

		if planCandidateSet.MetaData.CategoriesRejected != nil {
			for _, category := range *planCandidateSet.MetaData.CategoriesRejected {
				planCandidateSetCategoryEntity := generated.PlanCandidateSetMetaDataCategory{
					ID:                 uuid.New().String(),
					PlanCandidateSetID: planCandidateSet.Id,
					Category:           category.Name,
					IsSelected:         false,
				}
				if err := planCandidateSetCategoryEntity.Insert(ctx, db, boil.Infer()); err != nil {
					return fmt.Errorf("failed to insert plan candidate set category: %v", err)
				}
			}
		}
	}

	// PlanCandidateを作成
	planCandidateRepository, err := NewPlanCandidateRepository(db)
	if err != nil {
		return fmt.Errorf("failed to create plan candidate repository: %v", err)
	}

	if err := planCandidateRepository.AddPlan(ctx, planCandidateSet.Id, planCandidateSet.Plans...); err != nil {
		return fmt.Errorf("failed to add plan: %v", err)
	}

	return nil
}
