package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"os"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/firestore/entity"
)

const (
	collectionPlaces           = "places"
	subCollectionGooglePlaces  = "google_places"
	subCollectionGoogleReviews = "google_place_reviews"
	subCollectionGooglePhotos  = "google_place_photos"
)

type PlaceRepository struct {
	client *firestore.Client
}

func (p PlaceRepository) SavePlacesFromGooglePlaces(ctx context.Context, places []models.Place) error {
	//TODO implement me
	panic("implement me")
}

func (p PlaceRepository) FindByLocation(ctx context.Context, location models.GeoLocation) ([]models.Place, error) {
	//TODO implement me
	panic("implement me")
}

func (p PlaceRepository) FindByGooglePlaceID(ctx context.Context, googlePlaceID string) (*models.Place, error) {
	query := p.collectionPlaces().Where("google_place_id", "==", googlePlaceID).Limit(1)
	iter := query.Documents(ctx)
	doc, err := iter.Next()
	if errors.Is(err, iterator.Done) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error while iterating documents: %v", err)
	}

	var placeEntity entity.PlaceEntity
	if err := doc.DataTo(&placeEntity); err != nil {
		return nil, fmt.Errorf("error while converting doc to entity: %v", err)
	}

	// TODO: Google Place の情報を取得する
	place := placeEntity.ToPlace()
	return &place, nil
}

func (p PlaceRepository) SaveGooglePlacePhotos(ctx context.Context, googlePlaceId string, photos []models.GooglePlacePhoto) error {
	//TODO implement me
	panic("implement me")
}

func (p PlaceRepository) SaveGooglePlaceDetail(ctx context.Context, googlePlaceId string, detail models.GooglePlaceDetail) error {
	//TODO implement me
	panic("implement me")
}

func NewPlaceRepository(ctx context.Context) (*PlaceRepository, error) {
	var options []option.ClientOption
	if os.Getenv("GCP_CREDENTIAL_FILE_PATH") != "" {
		options = append(options, option.WithCredentialsFile(os.Getenv("GCP_CREDENTIAL_FILE_PATH")))
	}

	client, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT_ID"), options...)
	if err != nil {
		return nil, fmt.Errorf("error while initializing firestore client: %v", err)
	}

	return &PlaceRepository{
		client: client,
	}, nil
}

func (p PlaceRepository) collectionPlaces() *firestore.CollectionRef {
	return p.client.Collection(collectionPlaces)
}

func (p PlaceRepository) docPlace(placeId string) *firestore.DocumentRef {
	return p.client.Collection(collectionPlaces).Doc(placeId)
}

func (p PlaceRepository) subCollectionGooglePlace(googlePlaceId string) *firestore.CollectionRef {
	return p.client.Collection(collectionPlaces).Doc(googlePlaceId).Collection(subCollectionGooglePlaces)
}

func (p PlaceRepository) subCollectionGooglePlaceReview(googlePlaceId string) *firestore.CollectionRef {
	return p.client.Collection(collectionPlaces).Doc(googlePlaceId).Collection(subCollectionGoogleReviews)
}

func (p PlaceRepository) subCollectionGooglePlacePhoto(googlePlaceId string) *firestore.CollectionRef {
	return p.client.Collection(collectionPlaces).Doc(googlePlaceId).Collection(subCollectionGooglePhotos)
}
