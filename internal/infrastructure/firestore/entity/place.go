package entity

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
)

// PlaceEntity
// EstimatedStayDuration Firestoreではuintをサポートしていないため，intにしている
// TODO: photos, thumbnailは削除する
type PlaceEntity struct {
	Id                    string                     `firestore:"id"`
	GooglePlaceId         *string                    `firestore:"google_place_id"`
	Name                  string                     `firestore:"name"`
	Location              GeoLocationEntity          `firestore:"location"`
	Thumbnail             *string                    `firestore:"thumbnail"`
	Photos                []string                   `firestore:"photos"`
	Images                []ImageEntity              `firestore:"images"`
	EstimatedStayDuration int                        `firestore:"estimated_stay_duration"`
	GooglePlaceReviews    *[]GooglePlaceReviewEntity `firestore:"google_place_reviews,omitempty"`
	Categories            []string                   `firestore:"categories"`
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
			thumbnail = utils.StrCopyPointerValue(image.Small)
		}

		photos = append(photos, image.Default())
	}

	var images []ImageEntity
	for _, image := range place.Images {
		images = append(images, ToImageEntity(image))
	}

	var categories []string
	for _, category := range place.Categories {
		categories = append(categories, category.Name)
	}

	return PlaceEntity{
		Id:                    place.Id,
		GooglePlaceId:         place.GooglePlaceId,
		Name:                  place.Name,
		Location:              ToGeoLocationEntity(place.Location),
		Thumbnail:             thumbnail,
		Photos:                photos,
		Images:                images,
		EstimatedStayDuration: int(place.EstimatedStayDuration),
		GooglePlaceReviews:    googlePlaceReviews,
		Categories:            categories,
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
	if entity.Images == nil {
		for _, photo := range entity.Photos {
			image, err := models.NewImage(nil, utils.StrOmitEmpty(photo))
			if err != nil {
				continue
			}

			images = append(images, *image)
		}
	} else {
		for _, image := range entity.Images {
			images = append(images, FromImageEntity(image))
		}
	}

	var categories []models.LocationCategory
	for _, category := range entity.Categories {
		c := models.GetCategoryOfName(category)
		if c != nil {
			categories = append(categories, *c)
		}
	}

	return models.Place{
		Id:                    entity.Id,
		GooglePlaceId:         entity.GooglePlaceId,
		Name:                  entity.Name,
		Location:              FromGeoLocationEntity(entity.Location),
		Images:                images,
		EstimatedStayDuration: uint(entity.EstimatedStayDuration),
		GooglePlaceReviews:    googlePlaceReviews,
		Categories:            categories,
	}
}
