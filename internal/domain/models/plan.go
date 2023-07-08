package models

import "fmt"

type Plan struct {
	Id            string       `json:"id"`
	Name          string       `json:"name"`
	Places        []Place      `json:"places"`
	Transitions   []Transition `json:"transitions"`
	TimeInMinutes uint         `json:"time_in_minutes"` // MEMO: 複数プレイスを扱うようになったら，区間ごとの移動時間も保持したい
}

// GetPlace 指定したIDの場所情報を取得する
func (p Plan) GetPlace(placeId string) *Place {
	for _, place := range p.Places {
		if place.Id == placeId {
			return &place
		}
	}
	return nil
}

// GetTransition 指定した場所から次の場所への移動情報を取得する
func (p Plan) GetTransition(startPlaceId string) (nextPlace *Place, duration uint, err error) {
	placeStart := p.GetPlace(startPlaceId)
	if placeStart == nil {
		err = fmt.Errorf("could not find place %s in plan", startPlaceId)
		return
	}

	var transition *Transition
	for _, t := range p.Transitions {
		if t.FromPlaceId == startPlaceId {
			transition = &t
			break
		}
	}
	if transition == nil {
		err = fmt.Errorf("could not find transition from %s in plan", startPlaceId)
		return
	}

	nextPlace = p.GetPlace(transition.ToPlaceId)
	if nextPlace == nil {
		err = fmt.Errorf("could not find place %s in plan", transition.ToPlaceId)
		return
	}

	return nextPlace, transition.Duration, nil
}
