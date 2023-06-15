// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type CachedCreatedPlans struct {
	Plans                         []*Plan `json:"plans,omitempty"`
	CreatedBasedOnCurrentLocation bool    `json:"createdBasedOnCurrentLocation"`
}

type CachedCreatedPlansInput struct {
	Session string `json:"session"`
}

type ChangePlacesOrderInPlanInput struct {
	Session string `json:"session"`
	ID      string `json:"id"`
}

type ChangePlacesOrderInPlanOutput struct {
	Plan *Plan `json:"plan"`
}

type CreatePlanByLocationInput struct {
	Latitude                      float64  `json:"latitude"`
	Longitude                     float64  `json:"longitude"`
	Categories                    []string `json:"categories,omitempty"`
	FreeTime                      *int     `json:"freeTime,omitempty"`
	CreatedBasedOnCurrentLocation *bool    `json:"createdBasedOnCurrentLocation,omitempty"`
}

type CreatePlanByLocationOutput struct {
	Session string  `json:"session"`
	Plans   []*Plan `json:"plans"`
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
	Name                  string       `json:"name"`
	Location              *GeoLocation `json:"location"`
	Photos                []string     `json:"photos"`
	EstimatedStayDuration int          `json:"estimatedStayDuration"`
}

type Plan struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Places        []*Place `json:"places"`
	TimeInMinutes int      `json:"timeInMinutes"`
	Description   *string  `json:"description,omitempty"`
}

type SavePlanFromCandidateInput struct {
	Session string `json:"session"`
	PlanID  string `json:"planId"`
}

type SavePlanFromCandidateOutput struct {
	Plan *Plan `json:"plan"`
}
