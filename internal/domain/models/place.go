package models

import (
	"math"
	"math/rand"
	"poroto.app/poroto/planner/internal/domain/utils"
	"regexp"
	"sort"
)

// Place 場所の情報
type Place struct {
	Id          string       `json:"id"`
	Google      GooglePlace  `json:"google"`
	Name        string       `json:"name"`
	Location    GeoLocation  `json:"location"`
	Address     *string      `json:"address"`
	LikeCount   int          `json:"like_count"`
	PlacePhotos []PlacePhoto `json:"place_photos"`
}

func (p Place) Categories() []LocationCategory {
	return GetCategoriesFromSubCategories(p.Google.Types)
}

func (p Place) MainCategory() *LocationCategory {
	if len(p.Categories()) == 0 {
		return nil
	}
	return &p.Categories()[0]
}

// ShortenAddress 番地等の細かい情報のない住所を取得する
func (p Place) ShortenAddress() *string {
	if p.Address == nil {
		return nil
	}

	re := regexp.MustCompile(`^(.*?)[0-9０-９]`)
	match := re.FindStringSubmatch(*p.Address)
	if len(match) > 1 {
		return utils.ToPointer(match[1])
	}

	// 数字が含まれていない
	return p.Address
}

// PlacePhotosSortedByUploadedAt 新しい画像が先頭になるように並び替える
func (p Place) PlacePhotosSortedByUploadedAt() []PlacePhoto {
	placePhotos := p.PlacePhotos
	sort.Slice(placePhotos, func(i, j int) bool {
		return placePhotos[i].CreatedAt.After(placePhotos[j].CreatedAt)
	})
	return placePhotos
}

func (p Place) EstimatedStayDuration() uint {
	categoryMain := p.MainCategory()
	if categoryMain == nil {
		return 0
	}
	return categoryMain.EstimatedStayDuration
}

// EstimatedPriceRange 価格帯を推定する
func (p Place) EstimatedPriceRange() (priceRange *PriceRange) {
	// TODO: 飲食店でprice_levelが0の場合は、価格帯が不明なので、nilを返す
	return PriceRangeFromGooglePriceLevel(p.Google.PriceLevel)
}

func (p Place) CreateTransition(destination Place) Transition {
	return Transition{
		FromPlaceId: &p.Id,
		ToPlaceId:   destination.Id,
		Duration:    p.Location.TravelTimeTo(destination.Location, 80.0),
	}
}

// ShufflePlaces 場所の順番をシャッフルする
func ShufflePlaces(places []Place) []Place {
	placesCopy := make([]Place, len(places))
	copy(placesCopy, places)

	// Fisher-Yatesアルゴリズム
	// https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
	for i := len(places) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		placesCopy[i], placesCopy[j] = placesCopy[j], placesCopy[i]
	}

	return placesCopy
}

// WilsonScoreLowerBound Wilson score confidence interval for a Bernoulli parameter
func WilsonScoreLowerBound(averageRating float64, totalReviews int, confidence float64, maxRating float64) float64 {
	if totalReviews == 0 {
		return 0
	}

	// Proportion of positive reviews, estimated from the average rating
	// Assuming that the maximum rating is 5
	p := averageRating / maxRating

	// Z-score for the desired confidence level (e.g., 1.96 for 95% confidence)
	z := confidence

	// Calculate the lower bound
	denominator := 1 + z*z/float64(totalReviews)
	center := p + z*z/(2*float64(totalReviews))
	margin := z * math.Sqrt(p*(1-p)/float64(totalReviews)+z*z/(4*float64(totalReviews)*float64(totalReviews)))
	lowerBound := (center - margin) / denominator

	return lowerBound
}

// SortPlacesByRating 場所を評価の高い順に並び替える
func SortPlacesByRating(places []Place) []Place {
	placesCopy := make([]Place, len(places))
	copy(placesCopy, places)

	sort.SliceStable(placesCopy, func(i, j int) bool {
		wilsonScorePlaceI := WilsonScoreLowerBound(float64(placesCopy[i].Google.Rating), placesCopy[i].Google.UserRatingsTotal, 0.95, 5)
		wilsonScorePlaceJ := WilsonScoreLowerBound(float64(placesCopy[j].Google.Rating), placesCopy[j].Google.UserRatingsTotal, 0.95, 5)
		return wilsonScorePlaceI > wilsonScorePlaceJ
	})

	return placesCopy
}
