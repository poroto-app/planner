package plan

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

type CreatePlansRequest struct {
	Location models.GeoLocation `json:"location"`
}

type CreatePlansResponse struct {
	Plans []models.Plan `json:"plans"`
}

func CreatePlans(c *gin.Context) {
	var request CreatePlansRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "request body is invalid",
		})
	}

	placesApi, err := places.NewPlacesApi()
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
	}

	placesSearched, err := placesApi.FindPlacesFromLocation(context.Background(), &places.FindPlacesFromLocationRequest{
		Location: places.Location{
			Latitude:  request.Location.Latitude,
			Longitude: request.Location.Longitude,
		},
		Radius: 1000,
	})
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Printf("error while fetching places: %v\n", err)
		return
	}

	// TODO: 移動距離ではなく、移動時間でやる
	var placesRecommend []places.Place
	placesInNear := FilterWithinDistanceRange(request.Location, 0, 500, placesSearched)
	placesInMiddle := FilterWithinDistanceRange(request.Location, 500, 1000, placesSearched)
	placesInFar := FilterWithinDistanceRange(request.Location, 1000, 2000, placesSearched)
	if len(placesInNear) > 0 {
		placesRecommend = append(placesRecommend, placesInNear[0])
	}
	if len(placesInMiddle) > 0 {
		placesRecommend = append(placesRecommend, placesInMiddle[0])
	}
	if len(placesInFar) > 0 {
		placesRecommend = append(placesRecommend, placesInFar[0])
	}

	var plans []models.Plan
	for _, placeSearched := range placesRecommend {
		placePhotos, err := placesApi.FetchPlacePhotos(context.Background(), placeSearched)
		if err != nil {
			continue
		}
		photos := []string{}
		for _, photo := range placePhotos {
			photos = append(photos, photo.ImageUrl)
		}

		plans = append(plans, models.Plan{
			Name: placeSearched.Name,
			Places: []models.Place{
				{
					Name:   placeSearched.Name,
					Photos: photos,
					Location: models.GeoLocation{
						Latitude:  placeSearched.Location.Latitude,
						Longitude: placeSearched.Location.Longitude,
					},
				},
			},
		})
	}

	c.JSON(http.StatusOK, CreatePlansResponse{
		Plans: plans,
	})
}

func FilterWithinDistanceRange(
	currentLocation models.GeoLocation,
	startInMeter float64,
	endInMeter float64,
	placesToFilter []places.Place,
) []places.Place {
	var placesWithInDistance []places.Place
	for _, place := range placesToFilter {
		distance := currentLocation.DistanceInMeter(models.GeoLocation{
			Latitude:  place.Location.Latitude,
			Longitude: place.Location.Longitude,
		})
		if startInMeter <= distance && distance < endInMeter {
			placesWithInDistance = append(placesWithInDistance, place)
		}
	}
	return placesWithInDistance
}
