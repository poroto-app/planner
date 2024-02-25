package mock

import (
	"context"

	"poroto.app/poroto/planner/internal/domain/models"
)

type PlaceRepository struct {
	Data map[string]models.Place
}

func NewPlaceRepository(data map[string]models.Place) *PlaceRepository {
	return &PlaceRepository{
		Data: data,
	}
}

func (p PlaceRepository) SavePlacePhotos(ctx context.Context, userId string, placeId string, photoUrl string, width int, height int) error {
	//TODO implement me
	panic("implement me")
}
