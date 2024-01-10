package rdb

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/factory"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
)

func savePlaces(ctx context.Context, db *sql.DB, places []models.Place) error {
	places = array.DistinctBy(places, func(place models.Place) string { return place.Id })
	for _, place := range places {
		placeEntity := generated.Place{ID: place.Id, Name: place.Name}
		if err := placeEntity.Insert(ctx, db, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert place: %v", err)
		}

		if place.Google.PlaceId == "" {
			continue
		}

		googlePlaceEntity := factory.NewGooglePlaceEntityFromGooglePlace(place.Google, place.Id)
		if _, err := queries.Raw(
			fmt.Sprintf(
				"INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, POINT(?, ?) )",
				generated.TableNames.GooglePlaces,
				generated.GooglePlaceColumns.GooglePlaceID,
				generated.GooglePlaceColumns.PlaceID,
				generated.GooglePlaceColumns.Name,
				generated.GooglePlaceColumns.FormattedAddress,
				generated.GooglePlaceColumns.Vicinity,
				generated.GooglePlaceColumns.PriceLevel,
				generated.GooglePlaceColumns.Rating,
				generated.GooglePlaceColumns.UserRatingsTotal,
				generated.GooglePlaceColumns.Latitude,
				generated.GooglePlaceColumns.Longitude,
				generated.GooglePlaceColumns.Location,
			),
			googlePlaceEntity.GooglePlaceID,
			googlePlaceEntity.PlaceID,
			googlePlaceEntity.Name,
			googlePlaceEntity.FormattedAddress,
			googlePlaceEntity.Vicinity,
			googlePlaceEntity.PriceLevel,
			googlePlaceEntity.Rating,
			googlePlaceEntity.UserRatingsTotal,
			googlePlaceEntity.Latitude,
			googlePlaceEntity.Longitude,
			googlePlaceEntity.Longitude,
			googlePlaceEntity.Latitude,
		).Exec(db); err != nil {
			return fmt.Errorf("failed to insert google place: %w", err)
		}

		photoReferenceSlice := factory.NewGooglePlacePhotoReferenceSliceFromGooglePlacePhotoReferences(place.Google.PhotoReferences, place.Google.PlaceId)
		if _, err := photoReferenceSlice.InsertAll(ctx, db, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert google place photo references: %v", err)
		}

		for _, photoReference := range place.Google.PhotoReferences {
			photoAttribution := factory.NewGooglePlacePhotoAttributionSliceFromPhotoReference(photoReference, place.Google.PlaceId)
			if _, err := photoAttribution.InsertAll(ctx, db, boil.Infer()); err != nil {
				return fmt.Errorf("failed to insert google place photo attributions: %v", err)
			}
		}

		googlePlaceTypeSlice := factory.NewGooglePlaceTypeSliceFromGooglePlace(place.Google)
		if _, err := googlePlaceTypeSlice.InsertAll(ctx, db, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert google place types: %v", err)
		}

		placeRepository, err := NewPlaceRepository(db)
		if err != nil {
			return fmt.Errorf("failed to create place repository: %v", err)
		}

		if place.Google.PlaceDetail != nil {
			if err := placeRepository.SaveGooglePlaceDetail(ctx, place.Google.PlaceId, *place.Google.PlaceDetail); err != nil {
				return fmt.Errorf("failed to save google place detail: %v", err)
			}
		}

		if place.Google.Photos != nil {
			if err := placeRepository.SaveGooglePlacePhotos(ctx, place.Google.PlaceId, *place.Google.Photos); err != nil {
				return fmt.Errorf("failed to save google place photos: %v", err)
			}
		}
	}

	return nil
}

func saveUsers(ctx context.Context, db *sql.DB, users []models.User) error {
	for _, user := range users {
		userEntity := generated.User{
			ID: user.Id,
		}
		if err := userEntity.Insert(ctx, db, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert user: %v", err)
		}
	}
	return nil
}

func savePlans(ctx context.Context, db *sql.DB, plans []models.Plan) error {
	for _, plan := range plans {
		planEntity := factory.NewPlanEntityFromDomainModel(plan)
		var startLocation models.GeoLocation
		if len(plan.Places) > 0 {
			startLocation = plan.Places[0].Location
		}
		if _, err := queries.Raw(
			fmt.Sprintf(
				"INSERT INTO %s (%s, %s, %s, %s) VALUES (?, ?, ?, POINT(?, ?))",
				generated.TableNames.Plans,
				generated.PlanColumns.ID,
				generated.PlanColumns.UserID,
				generated.PlanColumns.Name,
				generated.PlanColumns.Location,
			),
			planEntity.ID,
			planEntity.UserID,
			planEntity.Name,
			startLocation.Longitude,
			startLocation.Latitude,
		).ExecContext(ctx, db); err != nil {
			return fmt.Errorf("failed to insert plan: %w", err)
		}

		planPlaceSlice := factory.NewPlanPlaceSliceFromDomainMode(plan.Places, plan.Id)
		if _, err := planPlaceSlice.InsertAll(ctx, db, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert plan places: %v", err)
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
