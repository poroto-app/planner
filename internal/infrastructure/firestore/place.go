package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"os"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
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

func (p PlaceRepository) SavePlacesFromGooglePlace(ctx context.Context, place models.Place) error {
	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		placeEntity := entity.NewPlaceEntityFromPlace(place)
		if err := tx.Set(p.docPlace(placeEntity.Id), placeEntity); err != nil {
			return fmt.Errorf("error while saving place: %v", err)
		}

		// Google Place を保存する
		googlePlaceEntity := entity.GooglePlaceEntityFromGooglePlace(place.Google)
		if err := tx.Set(p.docGooglePlace(placeEntity.Id), googlePlaceEntity); err != nil {
			return fmt.Errorf("error while saving google place: %v", err)
		}

		// Place Detail を保存する
		if place.Google.PlaceDetail != nil {
			// PhotoReferenceを保存する
			if err := p.saveGooglePhotoReferencesTx(tx, placeEntity.Id, place.Google.PlaceDetail.PhotoReferences); err != nil {
				return fmt.Errorf("error while saving google place photos: %v", err)
			}

			// Reviewを保存する
			if err := p.saveGooglePlaceReviews(tx, placeEntity.Id, place.Google.PlaceDetail.Reviews); err != nil {
				return fmt.Errorf("error while saving google place reviews: %v", err)
			}

			// 開店時間を更新する
			if err := p.updateOpeningHours(tx, placeEntity.Id, *place.Google.PlaceDetail); err != nil {
				return fmt.Errorf("error while updating opening hours: %v", err)
			}
		}

		// Place Photo を保存する
		if place.Google.Photos != nil {
			if err := p.saveGooglePhotosTx(tx, placeEntity.Id, place.Google.PlaceId, *place.Google.Photos); err != nil {
				return fmt.Errorf("error while saving google place photos: %v", err)
			}
		}

		return nil
	}, firestore.MaxAttempts(3)); err != nil {
		return fmt.Errorf("error while running transaction: %v", err)
	}

	return nil
}

func (p PlaceRepository) FindByLocation(ctx context.Context, location models.GeoLocation) ([]models.Place, error) {
	type findPlaceEntityByGeoHashResult struct {
		placeEntity *[]entity.PlaceEntity
		err         error
	}

	// 各方向に 5km 以内の場所を取得
	geohashPrecision := 5
	geohashNeighbors := location.GeoHashOfNeighbors(uint(geohashPrecision))
	// TODO: 重複がないかを確認する
	geohashNeighbors = append(geohashNeighbors, location.GeoHash()[:geohashPrecision])

	// 各方向
	ch := make(chan findPlaceEntityByGeoHashResult, len(geohashNeighbors))
	for _, geoHash := range geohashNeighbors {
		go func(ch chan<- findPlaceEntityByGeoHashResult, geoHash string) {
			query := p.collectionPlaces().Where("geohash", ">=", geoHash).Where("geohash", "<=", geoHash+"\uf8ff")
			query = query.Limit(50)
			iter := query.Documents(ctx)
			snapshots, err := iter.GetAll()
			if err != nil {
				ch <- findPlaceEntityByGeoHashResult{
					placeEntity: nil,
					err:         fmt.Errorf("error while getting place entities: %v", err),
				}
				return
			}

			var placeEntities []entity.PlaceEntity
			for _, snapshot := range snapshots {
				var placeEntity entity.PlaceEntity
				if err := snapshot.DataTo(&placeEntity); err != nil {
					ch <- findPlaceEntityByGeoHashResult{
						placeEntity: nil,
						err:         fmt.Errorf("error while converting snapshot to place entity: %v", err),
					}
					return
				}
				placeEntities = append(placeEntities, placeEntity)
			}

			ch <- findPlaceEntityByGeoHashResult{
				placeEntity: &placeEntities,
				err:         nil,
			}
		}(ch, geoHash)
	}

	var placeEntities []entity.PlaceEntity
	for i := 0; i < len(geohashNeighbors); i++ {
		result := <-ch
		if result.err != nil {
			return nil, result.err
		}
		if result.placeEntity != nil {
			placeEntities = append(placeEntities, *result.placeEntity...)
		}
	}

	// TODO: 紐づくデータを取得
	var places []models.Place
	for _, placeEntity := range placeEntities {
		places = append(places, placeEntity.ToPlace())
	}

	return places, nil
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
		placeEntity, err := p.findByGooglePlaceIdTx(tx, googlePlaceId)
		if err != nil {
			return fmt.Errorf("error while finding place by google place id: %v", err)
		}
		if placeEntity == nil {
			return fmt.Errorf("place not found by google place id: %s", googlePlaceId)
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
	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// 事前に要素が存在するかを確認する
		placeEntity, err := p.findByGooglePlaceIdTx(tx, googlePlaceId)
		if err != nil {
			return fmt.Errorf("error while finding place by google place id: %v", err)
		}
		if placeEntity == nil {
			return fmt.Errorf("place not found by google place id: %s", googlePlaceId)
		}

		// PhotoReferenceを保存する
		if err := p.saveGooglePhotoReferencesTx(tx, placeEntity.Id, detail.PhotoReferences); err != nil {
			return fmt.Errorf("error while saving google place photos: %v", err)
		}

		// Reviewを保存する
		if err := p.saveGooglePlaceReviews(tx, placeEntity.Id, detail.Reviews); err != nil {
			return fmt.Errorf("error while saving google place reviews: %v", err)
		}

		// 開店時間を更新する
		if err := p.updateOpeningHours(tx, placeEntity.Id, detail); err != nil {
			return fmt.Errorf("error while updating opening hours: %v", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error while saving google place detail: %v", err)
	}

	return nil
}

func (p PlaceRepository) findByGooglePlaceIdTx(tx *firestore.Transaction, googlePlaceId string) (*entity.PlaceEntity, error) {
	query := p.collectionPlaces().Where("google_place_id", "==", googlePlaceId).Limit(1)
	iter := tx.Documents(query)
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

	return &placeEntity, nil
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

func (p PlaceRepository) saveGooglePhotoReferencesTx(tx *firestore.Transaction, placeId string, photoReferences []models.GooglePlacePhotoReference) error {
	ch := make(chan *models.GooglePlacePhotoReference, len(photoReferences))
	for _, photoReference := range photoReferences {
		go func(tx *firestore.Transaction, ch chan<- *models.GooglePlacePhotoReference, placeId string, photoReference models.GooglePlacePhotoReference) {
			if err := tx.Set(p.subCollectionGooglePlacePhoto(placeId).Doc(photoReference.PhotoReference), photoReference); err != nil {
				ch <- nil
			} else {
				ch <- &photoReference
			}
		}(tx, ch, placeId, photoReference)
	}

	return nil
}

func (p PlaceRepository) saveGooglePlaceReviews(tx *firestore.Transaction, placeId string, reviews []models.GooglePlaceReview) error {
	ch := make(chan *models.GooglePlaceReview, len(reviews))
	for _, review := range reviews {
		go func(tx *firestore.Transaction, ch chan<- *models.GooglePlaceReview, placeId string, review models.GooglePlaceReview) {
			// 重複したレビューが保存されないように ID を MD5(Time+Text+Language) で生成する
			// AuthorName 等は頻繁に変更される可能性があるため、IDには含めない
			hashContent := fmt.Sprintf("%d-%s-%s", review.Time, utils.StrEmptyIfNil(review.Text), utils.StrEmptyIfNil(review.Language))
			id := fmt.Sprintf("%x", md5.Sum([]byte(hashContent)))
			if err := tx.Set(p.subCollectionGooglePlaceReview(placeId).Doc(id), review); err != nil {
				ch <- nil
			} else {
				ch <- &review
			}
		}(tx, ch, placeId, review)
	}

	for i := 0; i < len(reviews); i++ {
		if review := <-ch; review == nil {
			return fmt.Errorf("error while saving google place review: %v", reviews[i])
		}
	}

	return nil
}

func (p PlaceRepository) updateOpeningHours(tx *firestore.Transaction, placeId string, placeDetail models.GooglePlaceDetail) error {
	// 開店時間を更新
	if placeDetail.OpeningHours != nil {
		openingHoursEntity := entity.GooglePlaceOpeningsEntityFromGooglePlaceOpeningHours(*placeDetail.OpeningHours)
		if err := tx.Update(p.docGooglePlace(placeId), []firestore.Update{
			{
				Path:  "opening_hours",
				Value: openingHoursEntity,
			},
			{
				Path:  "updated_at",
				Value: firestore.ServerTimestamp,
			},
		}); err != nil {
			return fmt.Errorf("error while saving google place detail: %v", err)
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

func (p PlaceRepository) docGooglePlace(placeId string) *firestore.DocumentRef {
	return p.client.Collection(collectionPlaces).Doc(placeId).Collection(subCollectionGooglePlaces).Doc("v1")
}

func (p PlaceRepository) subCollectionGooglePlaceReview(placeId string) *firestore.CollectionRef {
	return p.client.Collection(collectionPlaces).Doc(placeId).Collection(subCollectionGoogleReviews)
}

func (p PlaceRepository) subCollectionGooglePlacePhoto(placeId string) *firestore.CollectionRef {
	return p.client.Collection(collectionPlaces).Doc(placeId).Collection(subCollectionGooglePhotos)
}
