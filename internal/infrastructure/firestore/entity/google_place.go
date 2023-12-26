package entity

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
)

// GooglePlaceEntity NearbySearchで取得できるデータをキャッシュ
type GooglePlaceEntity struct {
	PlaceID          string                         `firestore:"place_id"`
	Name             string                         `firestore:"name"`
	Types            []string                       `firestore:"types"`
	Location         GeoLocationEntity              `firestore:"location"`
	OpenNow          bool                           `firestore:"open_now"`
	Rating           float32                        `firestore:"rating"`
	UserRatingsTotal int                            `firestore:"user_ratings_total"`
	PriceLevel       int                            `firestore:"price_level"`
	FormattedAddress string                         `firestore:"formatted_address"`
	Vicinity         string                         `firestore:"vicinity"`
	OpeningHours     *GooglePlaceOpeningHoursEntity `firestore:"opening_hours"`
}

func GooglePlaceEntityFromGooglePlace(place models.GooglePlace) GooglePlaceEntity {
	var openingHours *GooglePlaceOpeningHoursEntity
	if place.PlaceDetail != nil {
		if place.PlaceDetail.OpeningHours != nil {
			o := GooglePlaceOpeningsEntityFromGooglePlaceOpeningHours(*place.PlaceDetail.OpeningHours)
			openingHours = &o
		}
	}

	return GooglePlaceEntity{
		PlaceID:          place.PlaceId,
		Name:             place.Name,
		Types:            place.Types,
		Rating:           place.Rating,
		UserRatingsTotal: place.UserRatingsTotal,
		PriceLevel:       place.PriceLevel,
		OpeningHours:     openingHours,
		FormattedAddress: utils.StrEmptyIfNil(place.FormattedAddress),
		Vicinity:         utils.StrEmptyIfNil(place.Vicinity),
		Location: GeoLocationEntity{
			Latitude:  place.Location.Latitude,
			Longitude: place.Location.Longitude,
		},
	}
}

func (g GooglePlaceEntity) ToGooglePlace(photoEntities *[]GooglePlacePhotoEntity, reviewEntities *[]GooglePlaceReviewEntity) models.GooglePlace {
	location := models.GeoLocation{
		Latitude:  g.Location.Latitude,
		Longitude: g.Location.Longitude,
	}

	return models.GooglePlace{
		PlaceId:          g.PlaceID,
		Name:             g.Name,
		Types:            g.Types,
		Location:         location,
		Rating:           g.Rating,
		UserRatingsTotal: g.UserRatingsTotal,
		PriceLevel:       g.PriceLevel,
		FormattedAddress: utils.StrOmitEmpty(g.FormattedAddress),
		Vicinity:         utils.StrOmitEmpty(g.Vicinity),
		Photos:           g.toGooglePlacePhotos(photoEntities),
		PlaceDetail:      g.toGooglePlaceDetail(photoEntities, reviewEntities),
	}
}

func (g GooglePlaceEntity) toGooglePlaceDetail(photoEntities *[]GooglePlacePhotoEntity, reviewEntities *[]GooglePlaceReviewEntity) *models.GooglePlaceDetail {
	isOpeningHoursEmpty := g.OpeningHours == nil
	isPhotoEmpty := photoEntities == nil || len(*photoEntities) == 0
	isReviewEmpty := reviewEntities == nil || len(*reviewEntities) == 0
	if isOpeningHoursEmpty && isPhotoEmpty && isReviewEmpty {
		return nil
	}

	placeDetail := &models.GooglePlaceDetail{}
	if g.OpeningHours != nil {
		// Opening Hoursを取得
		o := g.OpeningHours.ToGooglePlaceOpeningHours()
		placeDetail.OpeningHours = &o
	}

	if photoEntities != nil {
		// Photo Referenceを取得
		var photoReferences []models.GooglePlacePhotoReference
		for _, photo := range *photoEntities {
			photoReferences = append(photoReferences, photo.ToGooglePlacePhotoReference())
		}
		placeDetail.PhotoReferences = photoReferences
	}

	if reviewEntities != nil {
		// Reviewを取得
		var reviews []models.GooglePlaceReview
		for _, reviewEntity := range *reviewEntities {
			reviews = append(reviews, reviewEntity.ToGooglePlaceReview())
		}
		placeDetail.Reviews = reviews
	}

	return placeDetail
}

func (g GooglePlaceEntity) toGooglePlacePhotos(photoEntities *[]GooglePlacePhotoEntity) *[]models.GooglePlacePhoto {
	if photoEntities == nil {
		return nil
	}

	var photos []models.GooglePlacePhoto
	for _, photoEntity := range *photoEntities {
		p := photoEntity.ToGooglePlacePhoto()
		if p != nil {
			photos = append(photos, *p)
		}
	}

	if len(photos) == 0 {
		return nil
	}

	return &photos
}
