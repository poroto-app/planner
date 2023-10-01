package factory

import (
	graphql "poroto.app/poroto/planner/graphql/model"
	"poroto.app/poroto/planner/internal/domain/models"
)

func PlaceFromDomainModel(place *models.Place) *graphql.Place {
	if place == nil {
		return nil
	}

	var images []*graphql.Image
	for _, image := range place.Images {
		images = append(images, &graphql.Image{
			Default: image.Default(),
			Small:   image.Small,
			Large:   image.Large,
		})
	}

	// TODO: DELETE ME
	var photos []string
	for _, image := range place.Images {
		photos = append(photos, image.Default())
	}

	var googlePlaceReviews []*graphql.GooglePlaceReview
	if place.GooglePlaceReviews != nil {
		googlePlaceReviews = make([]*graphql.GooglePlaceReview, len(*place.GooglePlaceReviews))
		for i, review := range *place.GooglePlaceReviews {
			googlePlaceReviews[i] = GooglePlaceReviewFromDomainModel(review)
		}
	}

	var placeCategories []*graphql.PlaceCategory
	for _, category := range place.Categories {
		placeCategories = append(placeCategories, &graphql.PlaceCategory{
			ID:   category.Name,
			Name: category.DisplayName,
		})
	}

	return &graphql.Place{
		ID:            place.Id,
		GooglePlaceID: place.GooglePlaceId,
		Name:          place.Name,
		Photos:        photos,
		Images:        images,
		Location: &graphql.GeoLocation{
			Latitude:  place.Location.Latitude,
			Longitude: place.Location.Longitude,
		},
		EstimatedStayDuration: int(place.EstimatedStayDuration),
		GoogleReviews:         googlePlaceReviews,
		Categories:            placeCategories,
	}
}
