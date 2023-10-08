package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.34

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"poroto.app/poroto/planner/graphql/factory"
	"poroto.app/poroto/planner/graphql/model"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/plan"
	"poroto.app/poroto/planner/internal/domain/services/plancandidate"
	"poroto.app/poroto/planner/internal/domain/services/plangen"
)

// CreatePlanByLocation is the resolver for the createPlanByLocation field.
func (r *mutationResolver) CreatePlanByLocation(ctx context.Context, input model.CreatePlanByLocationInput) (*model.CreatePlanByLocationOutput, error) {
	planGenService, err := plangen.NewService(ctx)
	if err != nil {
		log.Println("error while initializing plan generator service: ", err)
		return nil, fmt.Errorf("internal server error")
	}

	planCandidateService, err := plancandidate.NewService(ctx)
	if err != nil {
		log.Println("error while initializing plan candidate service: ", err)
		return nil, fmt.Errorf("internal server error")
	}

	// TODO: 必須パラメータにする
	createBasedOnCurrentLocation := false
	if input.CreatedBasedOnCurrentLocation != nil {
		createBasedOnCurrentLocation = *input.CreatedBasedOnCurrentLocation
	}

	locationStart := models.GeoLocation{
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
	}

	session := uuid.New().String()
	if input.Session != nil {
		session = *input.Session
	}

	plans, err := planGenService.CreatePlanByLocation(
		ctx,
		session,
		locationStart,
		input.GooglePlaceID,
		&input.CategoriesPreferred,
		&input.CategoriesDisliked,
		input.FreeTime,
		createBasedOnCurrentLocation,
	)
	if err != nil {
		log.Println(err)
	}

	var categoriesPreferred, categoriesDisliked *[]models.LocationCategory
	if input.CategoriesPreferred != nil {
		var categories []models.LocationCategory
		for _, categoryName := range input.CategoriesPreferred {
			category := models.GetCategoryOfName(categoryName)
			if category != nil {
				categories = append(categories, *category)
			}
		}
		categoriesPreferred = &categories
	}

	if input.CategoriesDisliked != nil {
		var categories []models.LocationCategory
		for _, categoryName := range input.CategoriesDisliked {
			category := models.GetCategoryOfName(categoryName)
			if category != nil {
				categories = append(categories, *category)
			}
		}
		categoriesDisliked = &categories
	}

	// TODO: ServiceではPlanではなくPlanCandidateを生成するようにし、保存まで行う
	if err := planCandidateService.SavePlanCandidate(ctx, session, *plans, models.PlanCandidateMetaData{
		CategoriesPreferred:           categoriesPreferred,
		CategoriesRejected:            categoriesDisliked,
		FreeTime:                      input.FreeTime,
		CreatedBasedOnCurrentLocation: createBasedOnCurrentLocation,
		LocationStart:                 &locationStart,
	}); err != nil {
		log.Println("error while caching plan candidate: ", err)
	}

	return &model.CreatePlanByLocationOutput{
		Session: session,
		Plans:   factory.PlansFromDomainModel(plans, &locationStart),
	}, nil
}

// CreatePlanByPlace is the resolver for the createPlanByPlace field.
func (r *mutationResolver) CreatePlanByPlace(ctx context.Context, input model.CreatePlanByPlaceInput) (*model.CreatePlanByPlaceOutput, error) {
	planGenService, err := plangen.NewService(ctx)
	if err != nil {
		log.Println("error while initializing plan generator service: ", err)
		return nil, fmt.Errorf("internal server error")
	}

	planCreated, err := planGenService.CreatePlanFromPlace(
		ctx,
		input.Session,
		input.PlaceID,
	)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("internal server error")
	}

	graphqlPlan, err := factory.PlanFromDomainModel(*planCreated, nil)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("internal server error")
	}

	return &model.CreatePlanByPlaceOutput{
		Plan: graphqlPlan,
	}, nil
}

// ChangePlacesOrderInPlanCandidate is the resolver for the changePlacesOrderInPlanCandidate field.
func (r *mutationResolver) ChangePlacesOrderInPlanCandidate(ctx context.Context, input model.ChangePlacesOrderInPlanCandidateInput) (*model.ChangePlacesOrderInPlanCandidateOutput, error) {
	service, err := plancandidate.NewService(ctx)
	if err != nil {
		log.Println(fmt.Errorf("error while initizalizing PlanService: %v", err))
		return nil, fmt.Errorf("internal server error")
	}

	var currentLocation *model.GeoLocation
	if input.CurrentLatitude != nil && input.CurrentLongitude != nil {
		currentLocation = &model.GeoLocation{
			Latitude:  *input.CurrentLatitude,
			Longitude: *input.CurrentLongitude,
		}
	}

	planUpdated, err := service.ChangePlacesOrderPlanCandidate(ctx, input.PlanID, input.Session, input.PlaceIds, currentLocation)
	if err != nil {
		return nil, fmt.Errorf("could not change places order")
	}

	planCandidate, err := service.FindPlanCandidate(ctx, input.Session)
	if err != nil {
		log.Println(fmt.Errorf("error while finding plan candidate: %v", err))
		return nil, fmt.Errorf("internal server error")
	}

	graphqlPlan, err := factory.PlanFromDomainModel(*planUpdated, planCandidate.MetaData.LocationStart)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("internal server error")
	}

	return &model.ChangePlacesOrderInPlanCandidateOutput{
		Plan: graphqlPlan,
	}, nil
}

// SavePlanFromCandidate is the resolver for the savePlanFromCandidate field.
func (r *mutationResolver) SavePlanFromCandidate(ctx context.Context, input model.SavePlanFromCandidateInput) (*model.SavePlanFromCandidateOutput, error) {
	service, err := plan.NewService(ctx)
	if err != nil {
		log.Println(fmt.Errorf("error while initizalizing PlanService: %v", err))
		return nil, fmt.Errorf("internal server error")
	}

	planSaved, err := service.SavePlanFromPlanCandidate(ctx, input.Session, input.PlanID, input.AuthToken)
	if err != nil {
		log.Println(fmt.Errorf("error while initizalizing PlanService: %v", err))
		return nil, fmt.Errorf("could not save plan")
	}

	graphqlPlan, err := factory.PlanFromDomainModel(*planSaved, nil)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("internal server error")
	}

	return &model.SavePlanFromCandidateOutput{
		Plan: graphqlPlan,
	}, nil
}

// AddPlaceToPlanCandidate is the resolver for the addPlaceToPlanCandidate field.
func (r *mutationResolver) AddPlaceToPlanCandidate(ctx context.Context, input model.AddPlaceToPlanCandidateInput) (*model.AddPlaceToPlanCandidateOutput, error) {
	s, err := plancandidate.NewService(ctx)
	if err != nil {
		log.Println(fmt.Errorf("error while initizalizing PlanService: %v", err))
		return nil, fmt.Errorf("internal server error")
	}

	planCandidate, err := s.AddPlace(ctx, input.PlanCandidateID, input.PlanID, input.PlaceID)
	if err != nil {
		log.Println(fmt.Errorf("error while adding place to plan candidate: %v", err))
		return nil, fmt.Errorf("could not add place to plan candidate")
	}

	graphqlPlanInPlanCandidate, err := factory.PlanFromDomainModel(*planCandidate)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("internal server error")
	}

	return &model.AddPlaceToPlanCandidateOutput{
		Plan: graphqlPlanInPlanCandidate,
	}, nil
}

// DeletePlaceFromPlanCandidate is the resolver for the deletePlaceFromPlanCandidate field.
func (r *mutationResolver) DeletePlaceFromPlanCandidate(ctx context.Context, input model.DeletePlaceFromPlanCandidateInput) (*model.DeletePlaceFromPlanCandidateOutput, error) {
	s, err := plancandidate.NewService(ctx)
	if err != nil {
		log.Println(fmt.Errorf("error while initizalizing PlanService: %v", err))
		return nil, fmt.Errorf("internal server error")
	}

	planUpdated, err := s.RemovePlaceFromPlan(ctx, input.PlanCandidateID, input.PlanID, input.PlaceID)
	if err != nil {
		log.Println(fmt.Errorf("error while deleting place from plan candidate: %v", err))
		return nil, fmt.Errorf("could not delete place from plan candidate")
	}

	graphqlPlanInPlanCandidate, err := factory.PlanFromDomainModel(*planUpdated)
	if err != nil {
		log.Printf("error while converting plan to graphql model: %v", err)
		return nil, fmt.Errorf("internal server error")
	}

	return &model.DeletePlaceFromPlanCandidateOutput{
		PlanCandidateID: input.PlanCandidateID,
		Plan:            graphqlPlanInPlanCandidate,
	}, nil
}

// ReplacePlaceOfPlanCandidate is the resolver for the replacePlaceOfPlanCandidate field.
func (r *mutationResolver) ReplacePlaceOfPlanCandidate(ctx context.Context, input model.ReplacePlaceOfPlanCandidateInput) (*model.ReplacePlaceOfPlanCandidateOutput, error) {
	panic(fmt.Errorf("not implemented: ReplacePlaceOfPlanCandidate - replacePlaceOfPlanCandidate"))
}

// EditPlanTitleOfPlanCandidate is the resolver for the editPlanTitleOfPlanCandidate field.
func (r *mutationResolver) EditPlanTitleOfPlanCandidate(ctx context.Context, input model.EditPlanTitleOfPlanCandidateInput) (*model.EditPlanTitleOfPlanCandidateOutput, error) {
	panic(fmt.Errorf("not implemented: EditPlanTitleOfPlanCandidate - editPlanTitleOfPlanCandidate"))
}
