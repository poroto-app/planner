package models

type Transition struct {
	FromPlaceId string `json:"from_place_id"`
	ToPlaceId   string `json:"to_place_id"`
	Duration    uint   `json:"duration"`
}
