// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type CreatePlanByLocationInput struct {
	Latitude   float64  `json:"latitude"`
	Longitude  float64  `json:"longitude"`
	Categories []string `json:"categories,omitempty"`
}

type GeoLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type InterestCandidate struct {
	Categories []*LocationCategory `json:"categories"`
}

type LocationCategory struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Photo       string `json:"photo"`
}

type MatchInterestsInput struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Place struct {
	Name     string       `json:"name"`
	Location *GeoLocation `json:"location"`
	Photos   []string     `json:"photos,omitempty"`
}

type Plan struct {
	Name   string   `json:"name"`
	Places []*Place `json:"places"`
}
