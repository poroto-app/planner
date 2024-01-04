package factory

import (
	"fmt"
	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
)

func NewPlaceFromEntity(
	placeEntity entities.Place,
	googlePlaceSlice entities.GooglePlaceSlice,
) (*models.Place, error) {
	googlePlaceEntity, ok := array.Find(googlePlaceSlice, func(googlePlace *entities.GooglePlace) bool {
		if googlePlace == nil {
			return false
		}
		return googlePlace.PlaceID == placeEntity.ID
	})
	if !ok {
		return nil, fmt.Errorf("failed to find google place")
	}

	googlePlace, err := NewGooglePlaceFromEntity(*googlePlaceEntity)
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
