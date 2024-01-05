package factory

import (
	"fmt"
	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
)

func NewPlaceFromEntity(
	placeEntity generated.Place,
	googlePlaceEntity generated.GooglePlace,
	googlePlaceTypeSlice generated.GooglePlaceTypeSlice,
	googlePlacePhotoReferenceSlice generated.GooglePlacePhotoReferenceSlice,
	googlePlacePhotoAttributionSlice generated.GooglePlacePhotoAttributionSlice,
	googlePlacePhotoSlice generated.GooglePlacePhotoSlice,
	googlePlaceReviewSlice generated.GooglePlaceReviewSlice,
	googlePlaceOpeningPeriodSlice generated.GooglePlaceOpeningPeriodSlice,
) (*models.Place, error) {
	googlePlace, err := NewGooglePlaceFromEntity(
		googlePlaceEntity,
		googlePlaceTypeSlice,
		googlePlacePhotoReferenceSlice,
		googlePlacePhotoAttributionSlice,
		googlePlacePhotoSlice,
		googlePlaceReviewSlice,
		googlePlaceOpeningPeriodSlice,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to NewGooglePlaceFromEntity: %w", err)
	}

	if googlePlace == nil {
		return nil, err
	}

	return &models.Place{
		Id:        placeEntity.ID,
		Name:      placeEntity.Name,
		Location:  googlePlace.Location,
		Google:    *googlePlace,
		LikeCount: 0, // TODO: implement me
	}, nil
}

func NewPlaceEntityFromGooglePlaceEntity(googlePlace models.GooglePlace) generated.Place {
	return generated.Place{
		ID:   uuid.New().String(),
		Name: googlePlace.Name,
	}
}
