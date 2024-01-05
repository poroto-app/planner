package factory

import (
	"fmt"
	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
)

func NewPlaceFromEntity(
	placeEntity entities.Place,
	googlePlaceEntity entities.GooglePlace,
	googlePlaceTypeSlice entities.GooglePlaceTypeSlice,
	googlePlacePhotoReferenceSlice entities.GooglePlacePhotoReferenceSlice,
	googlePlacePhotoAttributionSlice entities.GooglePlacePhotoAttributionSlice,
	googlePlacePhotoSlice entities.GooglePlacePhotoSlice,
	googlePlaceReviewSlice entities.GooglePlaceReviewSlice,
	googlePlaceOpeningPeriodSlice entities.GooglePlaceOpeningPeriodSlice,
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

func NewPlaceEntityFromGooglePlaceEntity(googlePlace models.GooglePlace) entities.Place {
	return entities.Place{
		ID:   uuid.New().String(),
		Name: googlePlace.Name,
	}
}
