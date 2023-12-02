package factory

import (
	graphql "poroto.app/poroto/planner/graphql/model"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
)

func PlaceFromDomainModel(place *models.Place) *graphql.Place {
	if place == nil {
		return nil
	}

	// TODO: not nil な値にする
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
	} else {
		images = make([]*graphql.Image, 0)
	}

	// TODO: not nil な値にする
	var googlePlaceReviews []*graphql.GooglePlaceReview
	if place.Google.PlaceDetail != nil && place.Google.PlaceDetail.Reviews != nil {
		googlePlaceReviews = make([]*graphql.GooglePlaceReview, len(place.Google.PlaceDetail.Reviews))
		for i, review := range place.Google.PlaceDetail.Reviews {
			googlePlaceReviews[i] = GooglePlaceReviewFromDomainModel(review)
		}
	} else {
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
		GooglePlaceID: utils.StrPointer(place.Google.PlaceId),
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
	}
}
