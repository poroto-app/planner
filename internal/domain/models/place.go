package models

// Place 場所の情報
type Place struct {
	Id       string      `json:"id"`
	Google   GooglePlace `json:"google"`
	Name     string      `json:"name"`
	Location GeoLocation `json:"location"`
}

func NewPlaceFromGooglePlace(placeId string, googlePlace GooglePlace) Place {
	return Place{
		Id:       placeId,
		Google:   googlePlace,
		Name:     googlePlace.Name,
		Location: googlePlace.Location,
	}
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

func (p Place) IsSameCategoryPlace(other Place) bool {
	for _, categoryOfA := range p.Categories() {
		for _, categoryOfB := range other.Categories() {
			if categoryOfA.Name == categoryOfB.Name {
				return true
			}
		}
	}
	return false
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
