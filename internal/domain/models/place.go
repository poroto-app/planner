package models

type Place struct {
	Name     string      `json:"name"`
	Location GeoLocation `json:"location"`
	Photos   []string    `json:"photos"`
}
