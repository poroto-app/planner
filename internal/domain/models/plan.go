package models

type Plan struct {
	Id            string  `json:"id"`
	Name          string  `json:"name"`
	Places        []Place `json:"places"`
	TimeInMinutes uint    `json:"time_in_minutes"` // MEMO: 複数プレイスを扱うようになったら，区間ごとの移動時間も保持したい
	AuthorId      *string `json:"author_id"`
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
