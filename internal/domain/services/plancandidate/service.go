package plancandidate

import "poroto.app/poroto/planner/internal/domain/repository"

type Service struct {
	planCandidateRepository     repository.PlanCandidateRepository
	placeSearchResultRepository repository.PlaceSearchResultRepository
}
