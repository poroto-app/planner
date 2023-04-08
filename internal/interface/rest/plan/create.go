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

	places, err := placesApi.FindPlacesFromLocation(context.Background(), &places.FindPlacesFromLocationRequest{
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

	var plans []models.Plan
	for _, placeSearched := range places {
		plans = append(plans, models.Plan{
			Name: placeSearched.Name,
			Places: []models.Place{
				{
					Name: placeSearched.Name,
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
