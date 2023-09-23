package entity

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
)

// PlaceEntity
// EstimatedStayDuration Firestoreではuintをサポートしていないため，intにしている
type PlaceEntity struct {
	Id                    string                     `firestore:"id"`
	GooglePlaceId         *string                    `firestore:"google_place_id"`
	Name                  string                     `firestore:"name"`
	Location              GeoLocationEntity          `firestore:"location"`
	Thumbnail             *string                    `firestore:"thumbnail"`
	Photos                []string                   `firestore:"photos"`
	EstimatedStayDuration int                        `firestore:"estimated_stay_duration"`
	GooglePlaceReviews    *[]GooglePlaceReviewEntity `firestore:"google_place_reviews,omitempty"`
}

func ToPlaceEntity(place models.Place) PlaceEntity {
	var googlePlaceReviews *[]GooglePlaceReviewEntity
	if place.GooglePlaceReviews != nil {
		googlePlaceReviews = new([]GooglePlaceReviewEntity)
		for _, review := range *place.GooglePlaceReviews {
			*googlePlaceReviews = append(*googlePlaceReviews, ToGooglePlaceReviewEntity(review))
		}
	}

	// TODO: DELETE ME
	var thumbnail *string
	var photos []string
	for _, image := range place.Images {
		if thumbnail == nil && image.Small != nil {
			thumbnail = utils.StrPointer(*image.Small)
		}

		photos = append(photos, image.Default())
	}

	return PlaceEntity{
		Id:                    place.Id,
		GooglePlaceId:         place.GooglePlaceId,
		Name:                  place.Name,
		Location:              ToGeoLocationEntity(place.Location),
		Thumbnail:             thumbnail,
		Photos:                photos,
		EstimatedStayDuration: int(place.EstimatedStayDuration),
		GooglePlaceReviews:    googlePlaceReviews,
	}
}

func FromPlaceEntity(entity PlaceEntity) models.Place {
	var googlePlaceReviews *[]models.GooglePlaceReview
	if entity.GooglePlaceReviews != nil {
		googlePlaceReviews = new([]models.GooglePlaceReview)
		for _, review := range *entity.GooglePlaceReviews {
			*googlePlaceReviews = append(*googlePlaceReviews, FromGooglePlaceReviewEntity(review))
		}
	}

	// TODO: DELETE ME
	var images []models.Image
	for _, photo := range entity.Photos {
		image, err := models.NewImage(nil, utils.StrOmitEmpty(photo))
		if err != nil {
			continue
		}

		images = append(images, *image)
	}

	return models.Place{
		Id:                    entity.Id,
		GooglePlaceId:         entity.GooglePlaceId,
		Name:                  entity.Name,
		Location:              FromGeoLocationEntity(entity.Location),
		Images:                images,
		EstimatedStayDuration: uint(entity.EstimatedStayDuration),
		GooglePlaceReviews:    googlePlaceReviews,
	}
}
