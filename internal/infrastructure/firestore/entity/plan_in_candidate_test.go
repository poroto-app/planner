package entity

import (
	"log"
	"testing"

	"github.com/google/go-cmp/cmp"
	"poroto.app/poroto/planner/internal/domain/models"
)

func TestFromPlanInCandidateEntity(t *testing.T) {
	cases := []struct {
		name     string
		entity   PlanInCandidateEntity
		expected []models.Place
	}{
		{
			name: "プラン作成時から場所一覧の順序が並び替えられたケース",
			entity: PlanInCandidateEntity{
				Id:   "A",
				Name: "プラン候補A",
				Places: []PlaceEntity{
					{
						Id: "01",
					},
					{
						Id: "02",
					},
				},
				PlaceIdsOrdered: []string{
					"02",
					"01",
				},
			},
			expected: []models.Place{
				{
					Id: "02",
				},
				{
					Id: "01",
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result, err := fromPlanInCandidateEntity(
				c.entity.Id,
				c.entity.Name,
				c.entity.Places,
				c.entity.PlaceIdsOrdered,
				c.entity.TimeInMinutes,
				c.entity.Transitions,
			)
			if err != nil {
				log.Printf("Error occur while in converting Entity to Domain model: [%v]", err)
			}
			if diff := cmp.Diff(c.expected, result.Places); diff != "" {
				t.Errorf("expected %v, but got %v", c.expected, result)
			}
		})
	}
}
