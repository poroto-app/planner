package plan

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services"
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

	service, err := services.NewPlanService()
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	plans, err := service.CreatePlanByLocation(c.Request.Context(), request.Location)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, CreatePlansResponse{
		Plans: *plans,
	})
}
