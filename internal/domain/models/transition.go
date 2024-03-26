package models

// Transition 移動情報
// FromPlaceId がnilの場合は，出発地点が現在地であることを表す
// ToPlaceId は，移動先の場所ID
// Duration は，移動時間（分）
type Transition struct {
	FromPlaceId *string `json:"from_place_id"`
	ToPlaceId   string  `json:"to_place_id"`
	Duration    uint    `json:"duration"`
}

// CreateTransition　は移動情報を更新する（プラン内の場所の順番入れ替えなどの後に用いる）
// startLocation は現在地の座標を表す
func CreateTransition(places []Place, startLocation *GeoLocation) []Transition {
	transitions := make([]Transition, 0)

	// 現在位置から作成されたプラン or 場所指定で作成されたプラン
	if startLocation != nil {
		transitions = append(transitions, Transition{
			FromPlaceId: nil,
			ToPlaceId:   places[0].Id,
			Duration:    startLocation.TravelTimeTo(places[0].Location, 80.0),
		})
	}

	for i := range places {
		if i >= len(places)-1 {
			break
		}
		transitions = append(transitions, places[i].CreateTransition(places[i+1]))
	}

	return transitions
}
