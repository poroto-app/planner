// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

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

type Image struct {
	Default string  `json:"default"`
	Small   *string `json:"small,omitempty"`
	Large   *string `json:"large,omitempty"`
}

type InterestCandidate struct {
	Session    string              `json:"session"`
	Categories []*LocationCategory `json:"categories"`
}

type LocationCategory struct {
	Name            string  `json:"name"`
	DisplayName     string  `json:"displayName"`
	Photo           *string `json:"photo,omitempty"`
	DefaultPhotoURL string  `json:"defaultPhotoUrl"`
}

type MatchInterestsInput struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Place struct {
	ID                    string       `json:"id"`
	GooglePlaceID         *string      `json:"googlePlaceId,omitempty"`
	Name                  string       `json:"name"`
	Location              *GeoLocation `json:"location"`
	Photos                []string     `json:"photos"`
	Images                []*Image     `json:"images"`
	EstimatedStayDuration int          `json:"estimatedStayDuration"`
}

type Plan struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Places        []*Place      `json:"places"`
	TimeInMinutes int           `json:"timeInMinutes"`
	Description   *string       `json:"description,omitempty"`
	Transitions   []*Transition `json:"transitions"`
	AuthorID      *string       `json:"authorId,omitempty"`
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
	Plans  []*Plan `json:"plans"`
	Author *User   `json:"author"`
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

type ImageSize string

const (
	ImageSizeSmall ImageSize = "SMALL"
	ImageSizeLarge ImageSize = "LARGE"
)

var AllImageSize = []ImageSize{
	ImageSizeSmall,
	ImageSizeLarge,
}

func (e ImageSize) IsValid() bool {
	switch e {
	case ImageSizeSmall, ImageSizeLarge:
		return true
	}
	return false
}

func (e ImageSize) String() string {
	return string(e)
}

func (e *ImageSize) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ImageSize(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ImageSize", str)
	}
	return nil
}

func (e ImageSize) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
