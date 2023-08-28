package entity

import (
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
)

// PlanEntity は保存されたプランを示す
// GeoHash はプランの最初の場所のGeoHashを示す（プランは小さい範囲で作られるため、どこをとってもあまり変わらない）
// TimeInMinutes MEMO: Firestoreではuintをサポートしていないため，intにしている
type PlanEntity struct {
	Id            string               `firestore:"id"`
	Name          string               `firestore:"name"`
	Places        []PlaceEntity        `firestore:"places"`
	GeoHash       *string              `firestore:"geohash,omitempty"`
	TimeInMinutes int                  `firestore:"time_in_minutes"`
	Transitions   *[]TransitionsEntity `firestore:"transitions,omitempty"`
	CreatedAt     time.Time            `firestore:"created_at,omitempty,serverTimestamp"`
	UpdatedAt     time.Time            `firestore:"updated_at,omitempty"`
	AuthorId      *string              `firestore:"author_id,omitempty"`
}

func ToPlanEntity(plan models.Plan) PlanEntity {
	places := make([]PlaceEntity, len(plan.Places))
	for i, place := range plan.Places {
		places[i] = ToPlaceEntity(place)
	}

	var geohash *string
	if len(plan.Places) > 0 {
		value := plan.Places[0].Location.GeoHash()
		geohash = &value
	}

	return PlanEntity{
		Id:            plan.Id,
		Name:          plan.Name,
		Places:        places,
		GeoHash:       geohash,
		TimeInMinutes: int(plan.TimeInMinutes),
		Transitions:   ToTransitionsEntities(plan.Transitions),
		AuthorId:      plan.AuthorId,
	}
}

func FromPlanEntity(entity PlanEntity) models.Plan {
	return fromPlanEntity(
		entity.Id,
		entity.Name,
		entity.Places,
		entity.TimeInMinutes,
		entity.Transitions,
		entity.AuthorId,
	)
}

func fromPlanEntity(
	id string,
	name string,
	places []PlaceEntity,
	timeInMinutes int,
	transitions *[]TransitionsEntity,
	authorId *string,
) models.Plan {
	ps := make([]models.Place, len(places))
	for i, place := range places {
		ps[i] = FromPlaceEntity(place)
	}

	return models.Plan{
		Id:            id,
		Name:          name,
		Places:        ps,
		TimeInMinutes: uint(timeInMinutes),
		Transitions:   FromTransitionEntities(transitions),
		AuthorId:      authorId,
	}
}
