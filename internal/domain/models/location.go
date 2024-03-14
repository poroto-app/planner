package models

import (
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	"math"
)

type GeoLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// IsZero はゼロ値かどうかを判定する
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

// CalculateMBR 特定の位置からの距離を元に、緯度の差分を計算する
func (g GeoLocation) CalculateMBR(distance float64) (minLocation GeoLocation, maxLocation GeoLocation) {
	// 地球の半径（メートル単位）
	const earthRadius = 6371e3

	latLng := s2.LatLngFromDegrees(g.Latitude, g.Longitude)
	point := s2.PointFromLatLng(latLng)

	angle := s1.Angle(distance / earthRadius)
	cap := s2.CapFromCenterAngle(point, angle)
	rect := cap.RectBound()

	minLocation = GeoLocation{
		Latitude:  rect.Lo().Lat.Degrees(),
		Longitude: rect.Lo().Lng.Degrees(),
	}

	maxLocation = GeoLocation{
		Latitude:  rect.Hi().Lat.Degrees(),
		Longitude: rect.Hi().Lng.Degrees(),
	}

	return minLocation, maxLocation
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
