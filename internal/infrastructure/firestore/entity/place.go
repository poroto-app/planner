package entity

import "poroto.app/poroto/planner/internal/domain/models"

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
	var googlePlaceReviews []GooglePlaceReviewEntity
	if place.GooglePlaceReviews != nil {
		for _, review := range *place.GooglePlaceReviews {
			googlePlaceReviews = append(googlePlaceReviews, ToGooglePlaceReviewEntity(review))
		}
	}

	return PlaceEntity{
		Id:                    place.Id,
		GooglePlaceId:         place.GooglePlaceId,
		Name:                  place.Name,
		Location:              ToGeoLocationEntity(place.Location),
		Thumbnail:             place.Thumbnail,
		Photos:                place.Photos,
		EstimatedStayDuration: int(place.EstimatedStayDuration),
		GooglePlaceReviews:    &googlePlaceReviews,
	}
}

func FromPlaceEntity(entity PlaceEntity) models.Place {
	var googlePlaceReviews []models.GooglePlaceReview
	if entity.GooglePlaceReviews != nil {
		for _, review := range *entity.GooglePlaceReviews {
			googlePlaceReviews = append(googlePlaceReviews, FromGooglePlaceReviewEntity(review))
		}
	}

	return models.Place{
		Id:                    entity.Id,
		GooglePlaceId:         entity.GooglePlaceId,
		Name:                  entity.Name,
		Location:              FromGeoLocationEntity(entity.Location),
		Thumbnail:             entity.Thumbnail,
		Photos:                entity.Photos,
		EstimatedStayDuration: uint(entity.EstimatedStayDuration),
		GooglePlaceReviews:    &googlePlaceReviews,
	}
}
