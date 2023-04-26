package models

// Place 場所の情報
// Thumbnail サムネイル用の低画質な写真
type Place struct {
	Name      string      `json:"name"`
	Location  GeoLocation `json:"location"`
	Thumbnail *string     `json:"thumbnail"`
	Photos    []string    `json:"photos"`
}
