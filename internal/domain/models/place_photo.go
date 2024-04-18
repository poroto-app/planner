package models

import "time"

type PlacePhoto struct {
	PlaceId   string `json:"place_id"`
	UserId    string `json:"user_id"`
	PhotoUrl  string `json:"photo_url"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p PlacePhoto) ToImage() ImageSmallLarge {
	// TODO: SmallとLargeのURLを区別する
	return ImageSmallLarge{
		Small:          &p.PhotoUrl,
		Large:          &p.PhotoUrl,
		IsGooglePhotos: false,
	}
}
