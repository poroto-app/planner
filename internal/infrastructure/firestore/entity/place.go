package entity

import (
	"poroto.app/poroto/planner/internal/domain/models"
)

// PlaceEntity
// EstimatedStayDuration Firestoreではuintをサポートしていないため，intにしている
type PlaceEntity struct {
	Id                    string                     `firestore:"id"`
	GooglePlaceId         *string                    `firestore:"google_place_id"`
	Name                  string                     `firestore:"name"`
	Location              GeoLocationEntity          `firestore:"location"`
	Images                []ImageEntity              `firestore:"images"`
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

	var images []ImageEntity
	for _, image := range place.Images {
		images = append(images, ToImageEntity(image))
	}

	return PlaceEntity{
		Id:                    place.Id,
		GooglePlaceId:         place.GooglePlaceId,
		Name:                  place.Name,
		Location:              ToGeoLocationEntity(place.Location),
		Images:                images,
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

	var images []models.Image
	for _, image := range entity.Images {
		images = append(images, FromImageEntity(image))
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
