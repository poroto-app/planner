package entity

import (
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
)

// PlanEntity は保存されたプランを示す
// GeoHash はプランの最初の場所のGeoHashを示す（プランは小さい範囲で作られるため、どこをとってもあまり変わらない）
type PlanEntity struct {
	Id        string    `firestore:"id"`
	Name      string    `firestore:"name"`
	GeoHash   *string   `firestore:"geohash,omitempty"`
	PlaceIds  []string  `firestore:"place_ids"`
	CreatedAt time.Time `firestore:"created_at,omitempty,serverTimestamp"`
	UpdatedAt time.Time `firestore:"updated_at,omitempty"`
	AuthorId  *string   `firestore:"author_id,omitempty"`
}

func NewPlanEntityFromPlan(plan models.Plan) PlanEntity {
	var geohash *string
	if len(plan.Places) > 0 {
		value := plan.Places[0].Location.GeoHash()
		geohash = &value
	}

	placeIds := make([]string, len(plan.Places))
	for i, place := range plan.Places {
		placeIds[i] = place.Id
	}

	return PlanEntity{
		Id:        plan.Id,
		Name:      plan.Name,
		GeoHash:   geohash,
		PlaceIds:  placeIds,
		AuthorId:  plan.AuthorId,
		UpdatedAt: time.Now(),
	}
}

func FromPlanEntity(entity PlanEntity) models.Plan {
	return models.Plan{
		Id:       entity.Id,
		Name:     entity.Name,
		AuthorId: entity.AuthorId,
	}
}
