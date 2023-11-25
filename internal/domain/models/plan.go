package models

type Plan struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Places   []Place `json:"places"`
	AuthorId *string `json:"author_id"`
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

func (p Plan) Transitions(startLocation *GeoLocation) []Transition {
	return CreateTransition(p.Places, startLocation)
}

func (p Plan) TimeInMinutes(startLocation *GeoLocation) uint {
	transitions := p.Transitions(startLocation)
	var timeInMinute uint
	for _, t := range transitions {
		timeInMinute += t.Duration
	}

	for _, place := range p.Places {
		timeInMinute += place.EstimatedStayDuration()
	}
	return timeInMinute
}
