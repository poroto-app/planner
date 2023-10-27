package models

const (
	maxPriceOfLevel1 = 1000
	maxPriceOfLevel2 = 3000
	maxPriceOfLevel3 = 10000
	maxPriceOfLevel4 = 30000
)

type PriceRange struct {
	Min              int
	Max              int
	GooglePriceLevel int
}

// PriceRangeFromGooglePriceLevel Googleの価格帯レベルから価格帯を推定する
// SEE: https://developers.google.com/maps/documentation/places/web-service/details?hl=ja#Place-price_level
func PriceRangeFromGooglePriceLevel(priceLevel int) *PriceRange {
	switch priceLevel {
	case 0:
		return &PriceRange{
			Min:              0,
			Max:              0,
			GooglePriceLevel: priceLevel,
		}
	case 1:
		return &PriceRange{
			Min:              0,
			Max:              maxPriceOfLevel1,
			GooglePriceLevel: priceLevel,
		}
	case 2:
		return &PriceRange{
			Min:              0,
			Max:              maxPriceOfLevel2,
			GooglePriceLevel: priceLevel,
		}
	case 3:
		return &PriceRange{
			Min:              0,
			Max:              maxPriceOfLevel3,
			GooglePriceLevel: priceLevel,
		}
	case 4:
		return &PriceRange{
			Min:              0,
			Max:              maxPriceOfLevel4,
			GooglePriceLevel: priceLevel,
		}
	default:
		return nil
	}
}
