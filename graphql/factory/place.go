package factory

import (
	graphql "poroto.app/poroto/planner/graphql/model"
	"poroto.app/poroto/planner/internal/domain/models"
)

func PlaceFromDomainModel(place *models.Place) *graphql.Place {
	if place == nil {
		return nil
	}

	var googlePlaceReviews []*graphql.GooglePlaceReview
	if place.GooglePlaceReviews != nil {
		googlePlaceReviews = make([]*graphql.GooglePlaceReview, len(*place.GooglePlaceReviews))
		for i, review := range *place.GooglePlaceReviews {
			googlePlaceReviews[i] = GooglePlaceReviewFromDomainModel(review)
		}
	}

	return &graphql.Place{
		ID:     place.Id,
		Name:   place.Name,
		Photos: place.Photos,
		Location: &graphql.GeoLocation{
			Latitude:  place.Location.Latitude,
			Longitude: place.Location.Longitude,
		},
		EstimatedStayDuration: int(place.EstimatedStayDuration),
		GoogleReviews:         googlePlaceReviews,
	}
}
