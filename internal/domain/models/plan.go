package models

type Plan struct {
	Name   string  `json:"name"`
	Places []Place `json:"places"`
}
