package entity

import (
	"poroto.app/poroto/planner/internal/domain/models"
)

// PlaceEntity
type PlaceEntity struct {
	Id                 string                     `firestore:"id"`
	GooglePlaceId      *string                    `firestore:"google_place_id"`
	Name               string                     `firestore:"name"`
	Location           GeoLocationEntity          `firestore:"location"`
	Images             []ImageEntity              `firestore:"images"`
	GooglePlaceReviews *[]GooglePlaceReviewEntity `firestore:"google_place_reviews,omitempty"`
	Categories         []string                   `firestore:"categories"`
	PriceLevel         int                        `firestore:"price_level"`
}

func ToPlaceEntity(place models.Place) PlaceEntity {
	return PlaceEntity{
		Id:       place.Id,
		Name:     place.Name,
		Location: ToGeoLocationEntity(place.Location),
	}
}

func FromPlaceEntity(entity PlaceEntity) models.Place {
	var googlePlaceReviews *[]models.GooglePlaceReview
	if entity.GooglePlaceReviews != nil {
		googlePlaceReviews = new([]models.GooglePlaceReview)
		for _, review := range *entity.GooglePlaceReviews {
			*googlePlaceReviews = append(*googlePlaceReviews, review.ToGooglePlaceReview())
		}
	}

	var images []models.Image
	for _, image := range entity.Images {
		images = append(images, FromImageEntity(image))
	}

	var categories []models.LocationCategory
	for _, category := range entity.Categories {
		c := models.GetCategoryOfName(category)
		if c != nil {
			categories = append(categories, *c)
		}
	}

	return models.Place{
		Id:       entity.Id,
		Name:     entity.Name,
		Location: FromGeoLocationEntity(entity.Location),
	}
}
