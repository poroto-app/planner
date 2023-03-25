package plan

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"poroto.app/poroto/planner/internal/domain/place"
	"poroto.app/poroto/planner/internal/domain/plan"
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

	// TODO: 実際にGoogle Places APIを利用して取得する
	c.JSON(http.StatusOK, CreatePlansResponse{
		Plans: []plan.Plan{
			{
				Name: "カフェでほっと一息",
				Places: []place.Place{
					{
						Name: "スターバックス コーヒー 町田パリオ店",
						Location: place.GeoLocation{
							Latitude:  35.543261261835,
							Longitude: 139.44401761365,
						},
					},
				},
			},
		},
	})
}
