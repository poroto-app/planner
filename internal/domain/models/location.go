package models

import (
	"math"

	"github.com/mmcloughlin/geohash"
)

type GeoLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// IsZero はゼロ地値かどうかを判定する
func (g GeoLocation) IsZero() bool {
	return g.Latitude == 0.0 && g.Longitude == 0.0
}

// DistanceInMeter 2点間距離(メートル)
// SEE: https://www.geodatasource.com/developers/go
func (g GeoLocation) DistanceInMeter(another GeoLocation) float64 {
	locationA := g
	locationB := another

	radianLatitudeA := toRadian(locationA.Latitude)
	radianLatitudeB := toRadian(locationB.Latitude)
	radianTheta := toRadian(locationA.Longitude - locationB.Longitude)

	distance := math.Sin(radianLatitudeA)*math.Sin(radianLatitudeB) + math.Cos(radianLatitudeA)*math.Cos(radianLatitudeB)*math.Cos(radianTheta)
	if distance > 1 {
		distance = 1
	}

	distance = math.Acos(distance)
	distance = toDegree(distance)
	distance = distance * 60 * 1.1515
	distance = toKilometers(distance) * 1000

	return distance
}

func (g GeoLocation) GeoHash() string {
	return geohash.Encode(g.Latitude, g.Longitude)
}

// GeoHashOfNeighbors は指定した精度の周辺のGeoHashを返す
// 精度（precision）の値がどのような範囲を表すかは https://en.wikipedia.org/wiki/Geohash#Digits_and_precision_in_km を参照
// 例えば、precision=4の場合は、 GeoLocation　を中心とした 北、北東、東、南東、南、南西、西、北西の各方向に
// 20kmの範囲を表すGeoHashを返す
func (g GeoLocation) GeoHashOfNeighbors(precision uint) []string {
	centerGeoHash := geohash.EncodeWithPrecision(g.Latitude, g.Longitude, precision)
	return geohash.Neighbors(centerGeoHash)
}

func (g GeoLocation) TravelTimeTo(
	destination GeoLocation,
	meterPerMinutes float64,
) uint {
	var timeInMinutes uint = 0
	distance := g.DistanceInMeter(destination)
	if distance > 0.0 && meterPerMinutes > 0.0 {
		timeInMinutes = uint(distance / meterPerMinutes)
	}
	return timeInMinutes
}

func (g GeoLocation) Equal(other GeoLocation) bool {
	return g.Latitude == other.Latitude && g.Longitude == other.Longitude
}

func toRadian(degree float64) float64 {
	return math.Pi * degree / 180
}

func toDegree(radian float64) float64 {
	return radian * 180 / math.Pi
}

func toKilometers(mile float64) float64 {
	return mile * 1.609344
}
