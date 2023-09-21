package models

type GooglePlaceReview struct {
	Rating                uint    `json:"rating"`
	Text                  *string `json:"text"`
	Time                  int     `json:"time"`
	AuthorName            string  `json:"author_name"`
	AuthorUrl             *string `json:"author_url"`
	AuthorProfileImageUrl *string `json:"profile_photo_url"`
	Language              *string `json:"language"`
	OriginalLanguage      *string `json:"original_language"`
}
