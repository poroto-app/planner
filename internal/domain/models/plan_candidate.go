package models

import "time"

type PlanCandidate struct {
	Id        string
	Plans     []Plan
	ExpiresAt time.Time
}
