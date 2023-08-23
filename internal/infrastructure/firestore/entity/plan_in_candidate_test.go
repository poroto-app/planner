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
			name: "順序指定ID配列に重複がある場合は場所一覧をプラン作成時の順序でドメインモデルを作成",
			entity: PlanInCandidateEntity{
				Id:   "duplication",
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
					"01",
					"01",
				},
			},
			expected: []models.Place{
				{
					Id: "01",
				},
				{
					Id: "02",
				},
			},
		},
		{
			name: "順序指定ID配列と場所一覧の示す場所が一致しない場合はプラン作成時の順序でドメインモデルを作成",
			entity: PlanInCandidateEntity{
				Id:   "inconsistent",
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
					"10",
					"20",
				},
			},
			expected: []models.Place{
				{
					Id: "01",
				},
				{
					Id: "02",
				},
			},
		},
		{
			name: "順序指定ID配列と場所一覧の数が合わない場合はプラン作成時の順序でドメインモデルを作成",
			entity: PlanInCandidateEntity{
				Id:   "over_ids",
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
					"01",
					"02",
					"03",
				},
			},
			expected: []models.Place{
				{
					Id: "01",
				},
				{
					Id: "02",
				},
			},
		},
		{
			name: "プラン作成時から場所一覧の順序が並び替えられたケース",
			entity: PlanInCandidateEntity{
				Id:   "correct",
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
				log.Printf("error occur while in converting entity to domain model: [%v]", err)
			}
			if diff := cmp.Diff(c.expected, result.Places); diff != "" {
				t.Errorf("expected %v, but got %v", c.expected, result)
			}
		})
	}
}
