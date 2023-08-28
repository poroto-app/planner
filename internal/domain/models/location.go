package models

import (
	"math"

	"github.com/mmcloughlin/geohash"
)

type GeoLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
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
