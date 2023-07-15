package entity

import "poroto.app/poroto/planner/internal/domain/models"

type TransitionsEntity struct {
	from     *string `firestore:"from,omitempty"`
	to       string  `firestore:"to"`
	duration int     `firestore:"duration"`
}

func ToTransitionsEntities(transitions []models.Transition) *[]TransitionsEntity {
	ts := make([]TransitionsEntity, len(transitions))
	for i, transition := range transitions {
		ts[i] = TransitionsEntity{
			from:     transition.FromPlaceId,
			to:       transition.ToPlaceId,
			duration: int(transition.Duration),
		}
	}
	return &ts
}

func FromTransitionEntities(entities *[]TransitionsEntity) []models.Transition {
	if entities == nil {
		return []models.Transition{}
	}

	ts := make([]models.Transition, len(*entities))
	for i, entity := range *entities {
		ts[i] = models.Transition{
			FromPlaceId: entity.from,
			ToPlaceId:   entity.to,
			Duration:    uint(entity.duration),
		}
	}
	return ts
}
