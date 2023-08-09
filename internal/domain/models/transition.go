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
