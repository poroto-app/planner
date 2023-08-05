// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AvailablePlacesForPlan struct {
	Places []*Place `json:"places"`
}

type AvailablePlacesForPlanInput struct {
	Session string `json:"session"`
}

type CachedCreatedPlans struct {
	Plans                         []*Plan `json:"plans,omitempty"`
	CreatedBasedOnCurrentLocation bool    `json:"createdBasedOnCurrentLocation"`
}

type CachedCreatedPlansInput struct {
	Session string `json:"session"`
}

type ChangePlacesOrderInPlanCandidateInput struct {
	Session          string   `json:"session"`
	PlanID           string   `json:"planId"`
	PlaceIds         []string `json:"placeIds"`
	CurrentLatitude  *float64 `json:"currentLatitude,omitempty"`
	CurrentLongitude *float64 `json:"currentLongitude,omitempty"`
}

type ChangePlacesOrderInPlanCandidateOutput struct {
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

type CreatePlanByPlaceInput struct {
	Session string `json:"session"`
	PlaceID string `json:"placeId"`
}

type CreatePlanByPlaceOutput struct {
	Session string `json:"session"`
	Plan    *Plan  `json:"plan"`
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
	ID                    string       `json:"id"`
	Name                  string       `json:"name"`
	Location              *GeoLocation `json:"location"`
	Photos                []string     `json:"photos"`
	EstimatedStayDuration int          `json:"estimatedStayDuration"`
}

type Plan struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Places        []*Place      `json:"places"`
	TimeInMinutes int           `json:"timeInMinutes"`
	Description   *string       `json:"description,omitempty"`
	Transitions   []*Transition `json:"transitions"`
}

type SavePlanFromCandidateInput struct {
	Session string `json:"session"`
	PlanID  string `json:"planId"`
}

type SavePlanFromCandidateOutput struct {
	Plan *Plan `json:"plan"`
}

type Transition struct {
	From     *Place `json:"from,omitempty"`
	To       *Place `json:"to"`
	Duration int    `json:"duration"`
}
