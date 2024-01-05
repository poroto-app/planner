package factory

import (
	"fmt"
	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
)

func NewPlaceFromEntity(placeEntity entities.Place) (*models.Place, error) {
	if placeEntity.R == nil {
		return nil, fmt.Errorf("placeEntity.R is nil")
	}

	googlePlaceEntities := placeEntity.R.GetGooglePlaces()
	if len(googlePlaceEntities) == 0 || googlePlaceEntities[0] == nil {
		return nil, fmt.Errorf("placeEntity.R.GetGooglePlaces() is empty")
	}

	googlePlace, err := NewGooglePlaceFromEntity(*googlePlaceEntities[0])
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

func NewPlaceFromGooglePlaceEntity(googlePlaceEntity entities.GooglePlace) (*models.Place, error) {
	if googlePlaceEntity.R.Place == nil {
		return nil, fmt.Errorf("googlePlaceEntity.R.Place is nil")
	}

	googlePlace, err := NewGooglePlaceFromEntity(googlePlaceEntity)
	if googlePlace == nil {
		return nil, err
	}

	return &models.Place{
		Id:        googlePlaceEntity.R.Place.ID,
		Name:      googlePlace.Name,
		Location:  googlePlace.Location,
		Google:    *googlePlace,
		LikeCount: 0, // TODO: implement me
	}, err
}

func NewPlaceEntityFromGooglePlaceEntity(googlePlace models.GooglePlace) entities.Place {
	return entities.Place{
		ID:   uuid.New().String(),
		Name: googlePlace.Name,
	}
}
