package factory

import (
	"poroto.app/poroto/planner/internal/domain/models"
	graphql "poroto.app/poroto/planner/internal/interface/graphql/model"
)

func PlaceFromDomainModel(place *models.Place) *graphql.Place {
	if place == nil {
		return nil
	}

	var images []*graphql.Image
	if place.Google.Photos != nil {
		for _, photo := range *place.Google.Photos {
			image := photo.ToImage()
			images = append(images, &graphql.Image{
				Default: image.Default(),
				Small:   image.Small,
				Large:   image.Large,
			})
		}
	}
	if place.PlacePhotos != nil {
		for _, photo := range place.PlacePhotos {
			image := photo.ToImage()
			images = append(images, &graphql.Image{
				Default: image.Default(),
				Small:   image.Small,
				Large:   image.Large,
			})
		}
	}
	if len(images) == 0 {
		// not nil な値にする
		images = make([]*graphql.Image, 0)
	}

	var googlePlaceReviews []*graphql.GooglePlaceReview
	if place.Google.PlaceDetail != nil && place.Google.PlaceDetail.Reviews != nil {
		googlePlaceReviews = make([]*graphql.GooglePlaceReview, len(place.Google.PlaceDetail.Reviews))
		for i, review := range place.Google.PlaceDetail.Reviews {
			googlePlaceReviews[i] = GooglePlaceReviewFromDomainModel(review)
		}
	} else {
		// not nil な値にする
		googlePlaceReviews = make([]*graphql.GooglePlaceReview, 0)
	}

	var placeCategories []*graphql.PlaceCategory
	for _, category := range place.Categories() {
		placeCategories = append(placeCategories, &graphql.PlaceCategory{
			ID:   category.Name,
			Name: category.DisplayName,
		})
	}

	return &graphql.Place{
		ID:            place.Id,
		GooglePlaceID: place.Google.PlaceId,
		Name:          place.Name,
		Images:        images,
		Location: &graphql.GeoLocation{
			Latitude:  place.Location.Latitude,
			Longitude: place.Location.Longitude,
		},
		EstimatedStayDuration: int(place.EstimatedStayDuration()),
		GoogleReviews:         googlePlaceReviews,
		Categories:            placeCategories,
		PriceRange:            PriceRangeFromDomainModel(place.EstimatedPriceRange()),
		LikeCount:             int(place.LikeCount),
	}
}
