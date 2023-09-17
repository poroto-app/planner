package models

// Place 場所の情報
// Thumbnails サムネイル用の低画質な写真
// TODO: カテゴリを複数扱えるようにするために、 Category　を削除する
type Place struct {
	Id                    string             `json:"id"`
	GooglePlaceId         *string            `json:"google_place_id"`
	Name                  string             `json:"name"`
	Location              GeoLocation        `json:"location"`
	Thumbnails            []string           `json:"thumbnails"`
	Photos                []string           `json:"photos"`
	EstimatedStayDuration uint               `json:"estimated_stay_duration"`
	Category              string             `json:"category"`
	Categories            []LocationCategory `json:"categories"`
}
