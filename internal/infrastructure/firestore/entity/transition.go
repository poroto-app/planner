package entity

import "poroto.app/poroto/planner/internal/domain/models"

type TransitionsEntity struct {
	FromPlaceId *string `firestore:"from,omitempty"`
	ToPlaceId   string  `firestore:"to"`
	Duration    int     `firestore:"duration"`
}

func ToTransitionsEntities(transitions []models.Transition) *[]TransitionsEntity {
	ts := make([]TransitionsEntity, len(transitions))
	for i, transition := range transitions {
		ts[i] = TransitionsEntity{
			FromPlaceId: transition.FromPlaceId,
			ToPlaceId:   transition.ToPlaceId,
			Duration:    int(transition.Duration),
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
			FromPlaceId: entity.FromPlaceId,
			ToPlaceId:   entity.ToPlaceId,
			Duration:    uint(entity.Duration),
		}
	}
	return ts
}
