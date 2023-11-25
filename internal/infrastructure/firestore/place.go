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
	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// 事前に保存する画像が存在するかを確認する
		query := p.collectionPlaces().Where("google_place_id", "==", googlePlaceId).Limit(1)
		iter := tx.Documents(query)
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			return fmt.Errorf("place not found: %s", googlePlaceId)
		}
		if err != nil {
			return fmt.Errorf("error while iterating documents: %v", err)
		}

		var placeEntity entity.PlaceEntity
		if err := doc.DataTo(&placeEntity); err != nil {
			return fmt.Errorf("error while converting doc to entity: %v", err)
		}

		// 画像を保存する
		if err := p.saveGooglePhotosTx(tx, placeEntity.Id, googlePlaceId, photos); err != nil {
			return fmt.Errorf("error while saving google place photos: %v", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error while running transaction: %v", err)
	}

	return nil
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

// saveGooglePlaceTx はGoogle Places APIから取得された複数の画像を同時に保存する
// 一枚でも保存できなかった場合はエラーを返す
func (p PlaceRepository) saveGooglePhotosTx(tx *firestore.Transaction, placeId string, googlePlaceId string, photos []models.GooglePlacePhoto) error {
	ch := make(chan *models.GooglePlacePhoto, len(photos))
	for _, photo := range photos {
		go func(tx *firestore.Transaction, ch chan<- *models.GooglePlacePhoto, googlePlaceId string, photo models.GooglePlacePhoto) {
			if err := tx.Set(p.subCollectionGooglePlacePhoto(placeId).Doc(photo.PhotoReference), photo); err != nil {
				ch <- nil
			} else {
				ch <- &photo
			}
		}(tx, ch, googlePlaceId, photo)
	}

	for i := 0; i < len(photos); i++ {
		if photo := <-ch; photo == nil {
			return fmt.Errorf("error while saving google place photo: %v", photos[i])
		}
	}

	return nil
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
