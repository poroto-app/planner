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
	Session                       *string  `json:"session,omitempty"`
	Latitude                      float64  `json:"latitude"`
	Longitude                     float64  `json:"longitude"`
	GooglePlaceID                 *string  `json:"googlePlaceId,omitempty"`
	CategoriesPreferred           []string `json:"categoriesPreferred,omitempty"`
	CategoriesDisliked            []string `json:"categoriesDisliked,omitempty"`
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

type FirebaseUserInput struct {
	FirebaseUserID    string `json:"firebaseUserId"`
	FirebaseAuthToken string `json:"firebaseAuthToken"`
}

type GeoLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type InterestCandidate struct {
	Session    string              `json:"session"`
	Categories []*LocationCategory `json:"categories"`
}

type LocationCategory struct {
	Name            string `json:"name"`
	DisplayName     string `json:"displayName"`
	Photo           string `json:"photo"`
	DefaultPhotoURL string `json:"defaultPhotoUrl"`
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
	Author        *User         `json:"author,omitempty"`
}

type PlansByLocationInput struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Limit     *int    `json:"limit,omitempty"`
	PageKey   *string `json:"pageKey,omitempty"`
}

type PlansByLocationOutput struct {
	Plans   []*Plan `json:"plans"`
	PageKey *string `json:"pageKey,omitempty"`
}

type PlansByUserInput struct {
	UserID string `json:"userId"`
}

type PlansByUserOutput struct {
	Plans []*Plan `json:"plans"`
}

type SavePlanFromCandidateInput struct {
	Session   string  `json:"session"`
	PlanID    string  `json:"planId"`
	AuthToken *string `json:"authToken,omitempty"`
}

type SavePlanFromCandidateOutput struct {
	Plan *Plan `json:"plan"`
}

type Transition struct {
	From     *Place `json:"from,omitempty"`
	To       *Place `json:"to"`
	Duration int    `json:"duration"`
}

type User struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	PhotoURL *string `json:"photoUrl,omitempty"`
}
