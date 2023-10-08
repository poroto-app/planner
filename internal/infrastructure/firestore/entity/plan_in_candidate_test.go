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
		expected *models.Plan
	}{
		{
			name: "順序指定ID配列に重複がある場合は空のプランを返す",
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
			expected: nil,
		},
		{
			name: "順序指定ID配列と場所一覧の示す場所が一致しない場合は空のプランを返す",
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
			expected: nil,
		},
		{
			name: "順序指定ID配列と場所一覧の数が合わない場合は空のプランを返す",
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
			expected: nil,
		},
		{
			name: "正常な場合はプラン作成時から場所一覧の順序が並び替えられたプランを返す",
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
			expected: &models.Plan{
				Id:   "correct",
				Name: "プラン候補A",
				Places: []models.Place{
					{
						Id: "02",
					},
					{
						Id: "01",
					},
				},
				TimeInMinutes: 0,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result, err := FromPlanInCandidateEntity(
				c.entity.Id,
				c.entity.Name,
				c.entity.Places,
				c.entity.PlaceIdsOrdered,
				c.entity.TimeInMinutes,
			)
			if err != nil {
				log.Printf("error occur while in converting entity to domain model: [%v]", err)
			}
			if diff := cmp.Diff(c.expected, result); diff != "" {
				t.Errorf("expected %v, but got %v", c.expected, result)
			}
		})
	}
}

func TestValidatePlanInCandidateEntity(t *testing.T) {
	cases := []struct {
		name            string
		entity          []PlaceEntity
		placeIdsOrdered []string
		expected        bool
	}{
		{
			name: "順序指定ID配列に重複がある場合は false",
			entity: []PlaceEntity{
				{
					Id: "01",
				},
				{
					Id: "02",
				},
			},
			placeIdsOrdered: []string{
				"01",
				"01",
			},
			expected: false,
		},
		{
			name: "順序指定ID配列と場所一覧の示す場所が一致しない場合は false",
			entity: []PlaceEntity{
				{
					Id: "01",
				},
				{
					Id: "02",
				},
			},
			placeIdsOrdered: []string{
				"10",
				"20",
			},
			expected: false,
		},
		{
			name: "順序指定ID配列と場所一覧の示す場所が一致しない場合は false",
			entity: []PlaceEntity{
				{
					Id: "01",
				},
				{
					Id: "02",
				},
			},
			placeIdsOrdered: []string{
				"10",
				"20",
			},
			expected: false,
		},
		{
			name: "順序指定ID配列と場所一覧の数が合わない場合は false",
			entity: []PlaceEntity{
				{
					Id: "01",
				},
				{
					Id: "02",
				},
			},
			placeIdsOrdered: []string{
				"01",
				"02",
				"03",
			},
			expected: false,
		},
		{
			name: "正常な場合は true",
			entity: []PlaceEntity{
				{
					Id: "01",
				},
				{
					Id: "02",
				},
			},
			placeIdsOrdered: []string{
				"02",
				"01",
			},
			expected: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := validatePlanInCandidateEntity(
				c.entity,
				c.placeIdsOrdered,
			)
			if diff := cmp.Diff(c.expected, result); diff != "" {
				t.Errorf("expected %v, but got %v", c.expected, result)
			}
		})
	}
}
