package models

type PlacePhoto struct {
	Id       string `json:"id"`
	PlaceId  string `json:"place_id"`
	UserId   string `json:"user_id"`
	PhotoUrl string `json:"photo_url"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

func (p PlacePhoto) ToImage() ImageSmallLarge {
	// TODO: SmallとLargeのURLを区別する
	return ImageSmallLarge{
		Small: &p.PhotoUrl,
		Large: &p.PhotoUrl,
	}
}
