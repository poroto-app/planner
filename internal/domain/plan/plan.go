package plan

import "poroto.app/poroto/planner/internal/domain/place"

type Plan struct {
	Name   string        `json:"name"`
	Places []place.Place `json:"places"`
}
