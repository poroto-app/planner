package models

import (
	"sort"
)

type Plan struct {
	Id     string  `json:"id"`
	Name   string  `json:"name"`
	Places []Place `json:"places"`
	Author *User   `json:"author"`
}

// GetPlace 指定したIDの場所情報を取得する
func (p Plan) GetPlace(placeId string) *Place {
	for _, place := range p.Places {
		if place.Id == placeId {
			return &place
		}
	}
	return nil
}

// PlacesReorderedToMinimizeDistance は、スタート地点から移動が少なくなるように場所を並び替える
func (p Plan) PlacesReorderedToMinimizeDistance() []Place {
	if len(p.Places) == 0 {
		panic("Plan has no places")
	}

	placesReordered := make([]Place, len(p.Places))

	// 最初の場所をスタート地点とする
	startPlace := p.Places[0]
	placesReordered[0] = startPlace

	for i := 1; i < len(p.Places); i++ {
		// 一つ前の場所から近い順に並び替える
		placesSortedByDistanceToPrevLocation := make([]Place, len(p.Places))
		copy(placesSortedByDistanceToPrevLocation, p.Places)
		sort.Slice(placesSortedByDistanceToPrevLocation, func(a, b int) bool {
			prevLocation := placesReordered[i-1].Location
			return placesSortedByDistanceToPrevLocation[a].Location.DistanceInMeter(prevLocation) < placesSortedByDistanceToPrevLocation[b].Location.DistanceInMeter(prevLocation)
		})

		//　一つ前の場所から最も近い場所を選択する
		for _, place := range placesSortedByDistanceToPrevLocation {
			// すでに並び替えた場所は除外する
			var isAlreadyAdded bool
			for _, placeReordered := range placesReordered {
				if placeReordered.Id == place.Id {
					isAlreadyAdded = true
					break
				}
			}

			if isAlreadyAdded {
				continue
			}

			placesReordered[i] = place
			break
		}
	}

	return placesReordered
}

func (p Plan) Transitions(startLocation *GeoLocation) []Transition {
	return CreateTransition(p.Places, startLocation)
}

func (p Plan) TimeInMinutes(startLocation *GeoLocation) uint {
	transitions := p.Transitions(startLocation)
	var timeInMinute uint
	for _, t := range transitions {
		timeInMinute += t.Duration
	}

	for _, place := range p.Places {
		timeInMinute += place.EstimatedStayDuration()
	}
	return timeInMinute
}
