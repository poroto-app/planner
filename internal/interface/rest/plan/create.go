package plan

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"poroto.app/poroto/planner/internal/domain/place"
	"poroto.app/poroto/planner/internal/domain/plan"
	"poroto.app/poroto/planner/internal/infrastructure/api/google"
)

type CreatePlansRequest struct {
	Location place.GeoLocation `json:"location"`
}

type CreatePlansResponse struct {
	Plans []plan.Plan `json:"plans"`
}

func CreatePlans(c *gin.Context) {
	var request CreatePlansRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "request body is invalid",
		})
	}

	placesApi, err := google.NewPlacesApi()
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
	}

	places, err := placesApi.FindPlacesFromLocation(context.Background(), &google.FindPlacesFromLocationRequest{
		Location: google.Location{
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

	var plans []plan.Plan
	for _, placeSearched := range places {
		plans = append(plans, plan.Plan{
			Name: placeSearched.Name,
			Places: []place.Place{
				{
					Name: placeSearched.Name,
					Location: place.GeoLocation{
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
