package models

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

// 移動情報を更新する（プラン内の場所の順番入れ替えなどの後に用いる）
func (p Plan) RecreateTransition(startLocation *GeoLocation) []Transition {
	var firstFromPlaceId *string
	var firstFromPlaceLocation *GeoLocation
	var firstToPlaceId int

	// 現在位置から作成されたプラン or 場所指定で作成されたプラン
	if startLocation != nil {
		firstFromPlaceId = nil
		firstFromPlaceLocation = startLocation
		firstToPlaceId = 0
	} else {
		firstFromPlaceId = &p.Places[0].Id
		firstFromPlaceLocation = &p.Places[0].Location
		firstToPlaceId = 1
	}
	transitions := make([]Transition, len(p.Places)-firstToPlaceId)

	transitions[0] = Transition{
		FromPlaceId: firstFromPlaceId,
		ToPlaceId:   p.Places[firstToPlaceId].Id,
		Duration:    firstFromPlaceLocation.TravelTimeTo(p.Places[firstToPlaceId].Location, 80.0),
	}

	for i := 0; i < len(p.Places)-firstToPlaceId-1; i++ {
		transitions[i+1] = Transition{
			FromPlaceId: &p.Places[firstToPlaceId+i].Id,
			ToPlaceId:   p.Places[firstToPlaceId+i+1].Id,
			Duration:    p.Places[firstToPlaceId+i].Location.TravelTimeTo(p.Places[firstToPlaceId+i+1].Location, 80.0),
		}
	}

	return transitions
}
