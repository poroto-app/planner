package repository

import (
	"context"
	"poroto.app/poroto/planner/internal/domain/models"
)

type PlaceRepository interface {
	SavePlacesFromGooglePlaces(ctx context.Context, places []models.Place) error

	FindByLocation(ctx context.Context, location models.GeoLocation) ([]models.Place, error)

	FindByGooglePlaceID(ctx context.Context, googlePlaceID string) (*models.Place, error)

	SaveGooglePlacePhotos(ctx context.Context, googlePlaceId string, photos []models.GooglePlacePhoto) error

	SaveGooglePlaceDetail(ctx context.Context, googlePlaceId string, detail models.GooglePlaceDetail) error
}
