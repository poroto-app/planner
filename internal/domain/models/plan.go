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
// startLocation は現在地の座標を表す
func (p Plan) RecreateTransition(startLocation *GeoLocation) []Transition {
	transitions := make([]Transition, 0)

	// 現在位置から作成されたプラン or 場所指定で作成されたプラン
	if startLocation != nil {
		transitions = append(transitions, Transition{
			FromPlaceId: nil,
			ToPlaceId:   p.Places[0].Id,
			Duration:    startLocation.TravelTimeTo(p.Places[0].Location, 80.0),
		})

		for i, place := range p.Places {
			if i >= len(p.Places)-1 {
				break
			}

			transitions = append(transitions, Transition{
				FromPlaceId: &p.Places[i].Id,
				ToPlaceId:   p.Places[i+1].Id,
				Duration:    place.Location.TravelTimeTo(p.Places[i+1].Location, 80.0),
			})
		}
	} else {
		for i, place := range p.Places {
			if i >= len(p.Places)-1 {
				break
			}
			transitions = append(transitions, Transition{
				FromPlaceId: &p.Places[i].Id,
				ToPlaceId:   p.Places[i+1].Id,
				Duration:    place.Location.TravelTimeTo(p.Places[i+1].Location, 80.0),
			})
		}
	}
	return transitions
}
