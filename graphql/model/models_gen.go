// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type AddPlaceToPlanCandidateAfterPlaceInput struct {
	PlanCandidateID string `json:"planCandidateId"`
	PlanID          string `json:"planId"`
	PlaceID         string `json:"placeId"`
	PreviousPlaceID string `json:"previousPlaceId"`
}

type AddPlaceToPlanCandidateAfterPlaceOutput struct {
	PlanCandidateID string `json:"planCandidateId"`
	Plan            *Plan  `json:"plan"`
}

type AutoReorderPlacesInPlanCandidateInput struct {
	PlanCandidateID string `json:"planCandidateId"`
	PlanID          string `json:"planId"`
}

type AutoReorderPlacesInPlanCandidateOutput struct {
	PlanCandidateID string `json:"planCandidateId"`
	Plan            *Plan  `json:"plan"`
}

type AvailablePlacesForPlan struct {
	Places []*Place `json:"places"`
}

type AvailablePlacesForPlanInput struct {
	Session string `json:"session"`
}

type CachedCreatedPlans struct {
	Plans                         []*Plan `json:"plans"`
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

type DeletePlaceFromPlanCandidateInput struct {
	PlanCandidateID string `json:"planCandidateId"`
	PlanID          string `json:"planId"`
	PlaceID         string `json:"placeId"`
}

type DeletePlaceFromPlanCandidateOutput struct {
	PlanCandidateID string `json:"planCandidateId"`
	Plan            *Plan  `json:"plan"`
}

type EditPlanTitleOfPlanCandidateInput struct {
	PlanCandidateID string `json:"planCandidateId"`
	PlanID          string `json:"planId"`
	Title           string `json:"title"`
}

type EditPlanTitleOfPlanCandidateOutput struct {
	PlanCandidateID string `json:"planCandidateId"`
	Plan            *Plan  `json:"plan"`
}

type FirebaseUserInput struct {
	FirebaseUserID    string `json:"firebaseUserId"`
	FirebaseAuthToken string `json:"firebaseAuthToken"`
}

type GeoLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type GooglePlaceReview struct {
	Rating           int     `json:"rating"`
	Text             *string `json:"text,omitempty"`
	Time             int     `json:"time"`
	AuthorName       string  `json:"authorName"`
	AuthorURL        *string `json:"authorUrl,omitempty"`
	AuthorPhotoURL   *string `json:"authorPhotoUrl,omitempty"`
	Language         *string `json:"language,omitempty"`
	OriginalLanguage *string `json:"originalLanguage,omitempty"`
}

type Image struct {
	Default string  `json:"default"`
	Small   *string `json:"small,omitempty"`
	Large   *string `json:"large,omitempty"`
}

type LikeToPlaceInPlanCandidateInput struct {
	PlanCandidateID string `json:"planCandidateId"`
	PlaceID         string `json:"placeId"`
	Like            bool   `json:"like"`
}

type LikeToPlaceInPlanCandidateOutput struct {
	Plans []*Plan `json:"plans"`
}

type LocationCategory struct {
	Name            string  `json:"name"`
	DisplayName     string  `json:"displayName"`
	Photo           *string `json:"photo,omitempty"`
	DefaultPhotoURL string  `json:"defaultPhotoUrl"`
}

type NearbyLocationCategory struct {
	ID              string   `json:"id"`
	DisplayName     string   `json:"displayName"`
	Places          []*Place `json:"places"`
	DefaultPhotoURL string   `json:"defaultPhotoUrl"`
}

type NearbyPlaceCategoriesInput struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type NearbyPlaceCategoryOutput struct {
	PlanCandidateID string                    `json:"planCandidateId"`
	Categories      []*NearbyLocationCategory `json:"categories"`
}

type Place struct {
	ID                    string               `json:"id"`
	GooglePlaceID         string               `json:"googlePlaceId"`
	Name                  string               `json:"name"`
	Location              *GeoLocation         `json:"location"`
	Images                []*Image             `json:"images"`
	EstimatedStayDuration int                  `json:"estimatedStayDuration"`
	GoogleReviews         []*GooglePlaceReview `json:"googleReviews"`
	Categories            []*PlaceCategory     `json:"categories"`
	PriceRange            *PriceRange          `json:"priceRange,omitempty"`
	LikeCount             int                  `json:"likeCount"`
}

type PlaceCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PlacesToAddForPlanCandidateInput struct {
	PlanCandidateID string `json:"planCandidateId"`
	PlanID          string `json:"planId"`
}

type PlacesToAddForPlanCandidateOutput struct {
	Places []*Place `json:"places"`
}

type PlacesToReplaceForPlanCandidateInput struct {
	PlanCandidateID string `json:"planCandidateId"`
	PlanID          string `json:"planId"`
	PlaceID         string `json:"placeId"`
}

type PlacesToReplaceForPlanCandidateOutput struct {
	Places []*Place `json:"places"`
}

type Plan struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Places        []*Place      `json:"places"`
	TimeInMinutes int           `json:"timeInMinutes"`
	Description   *string       `json:"description,omitempty"`
	Transitions   []*Transition `json:"transitions"`
	AuthorID      *string       `json:"authorId,omitempty"`
	LikedPlaceIds []string      `json:"likedPlaceIds"`
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

type PriceRange struct {
	PriceRangeMin    int `json:"priceRangeMin"`
	PriceRangeMax    int `json:"priceRangeMax"`
	GooglePriceLevel int `json:"googlePriceLevel"`
}

type ReplacePlaceOfPlanCandidateInput struct {
	PlanCandidateID  string `json:"planCandidateId"`
	PlanID           string `json:"planId"`
	PlaceIDToRemove  string `json:"placeIdToRemove"`
	PlaceIDToReplace string `json:"placeIdToReplace"`
}

type ReplacePlaceOfPlanCandidateOutput struct {
	PlanCandidateID string `json:"planCandidateId"`
	Plan            *Plan  `json:"plan"`
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
