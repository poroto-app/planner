package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.34

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/plan"
	"poroto.app/poroto/planner/internal/domain/services/plancandidate"
	"poroto.app/poroto/planner/internal/domain/services/plangen"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/interface/graphql/factory"
	"poroto.app/poroto/planner/internal/interface/graphql/model"
)

// CreatePlanByLocation is the resolver for the createPlanByLocation field.
func (r *mutationResolver) CreatePlanByLocation(ctx context.Context, input model.CreatePlanByLocationInput) (*model.CreatePlanByLocationOutput, error) {
	planGenService, err := plangen.NewService(r.DB)
	if err != nil {
		log.Println("error while initializing plan generator service: ", err)
		return nil, fmt.Errorf("internal server error")
	}

	planCandidateService, err := plancandidate.NewService(r.DB)
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

	// プラン候補の作成
	var planCandidateId string
	if input.Session != nil {
		planCandidateId = *input.Session
	} else {
		planCandidateId = uuid.New().String()
		if err := planCandidateService.CreatePlanCandidate(ctx, planCandidateId); err != nil {
			log.Printf("error while creating plan candidate: %v", err)
			return nil, fmt.Errorf("internal server error")
		}
	}

	// プランの作成
	plans, err := planGenService.CreatePlanByLocation(
		ctx,
		plangen.CreatePlanByLocationInput{
			PlanCandidateId:              planCandidateId,
			LocationStart:                locationStart,
			CategoryNamesPreferred:       &input.CategoriesPreferred,
			CategoryNamesDisliked:        &input.CategoriesDisliked,
			FreeTime:                     input.FreeTime,
			CreateBasedOnCurrentLocation: createBasedOnCurrentLocation,
		},
	)
	if err != nil {
		log.Printf("error while creating plan by location: %v", err)
		return nil, fmt.Errorf("internal server error")
	}

	// 作成されたプランの保存
	if err := planCandidateService.SavePlans(ctx, plancandidate.SavePlansInput{
		PlanCandidateId:              planCandidateId,
		Plans:                        *plans,
		LocationStart:                &locationStart,
		CategoryNamesPreferred:       &input.CategoriesPreferred,
		CategoryNamesRejected:        &input.CategoriesDisliked,
		FreeTime:                     input.FreeTime,
		CreateBasedOnCurrentLocation: createBasedOnCurrentLocation,
	}); err != nil {
		log.Printf("error while saving plans of plan candidate(%s): %v", planCandidateId, err)
	}

	return &model.CreatePlanByLocationOutput{
		Session: planCandidateId,
		Plans:   factory.PlansFromDomainModel(plans, &locationStart),
	}, nil
}

// CreatePlanByPlace is the resolver for the createPlanByPlace field.
func (r *mutationResolver) CreatePlanByPlace(ctx context.Context, input model.CreatePlanByPlaceInput) (*model.CreatePlanByPlaceOutput, error) {
	planGenService, err := plangen.NewService(r.DB)
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

// CreatePlanByGooglePlaceID is the resolver for the createPlanByGooglePlaceId field.
func (r *mutationResolver) CreatePlanByGooglePlaceID(ctx context.Context, input model.CreatePlanByGooglePlaceIDInput) (*model.CreatePlanByGooglePlaceIDOutput, error) {
	logger, err := utils.NewLogger(utils.LoggerOption{Tag: "GraphQL"})
	if err != nil {
		log.Println("error while initializing logger: ", err)
		return nil, fmt.Errorf("internal server error")
	}

	planGenService, err := plangen.NewService(r.DB)
	if err != nil {
		logger.Error("error while initializing plan generator service", zap.Error(err))
		return nil, fmt.Errorf("internal server error")
	}

	planCandidateService, err := plancandidate.NewService(r.DB)
	if err != nil {
		logger.Error("error while initializing plan candidate service", zap.Error(err))
		return nil, fmt.Errorf("internal server error")
	}

	logger.Info(
		"CreatePlanByGooglePlaceID",
		zap.String("googlePlaceId", input.GooglePlaceID),
	)

	// プラン候補の作成
	var planCandidateId string
	if input.PlanCandidateID != nil {
		planCandidateId = *input.PlanCandidateID
	} else {
		planCandidateId = uuid.New().String()
		if err := planCandidateService.CreatePlanCandidate(ctx, planCandidateId); err != nil {
			logger.Error("error while creating plan candidate", zap.Error(err))
			return nil, fmt.Errorf("internal server error")
		}
	}

	// プランの作成
	createPlanByGooglePlaceidResult, err := planGenService.CreatePlanByGooglePlaceId(ctx, plangen.CreatePlanByGooglePlaceIdInput{
		PlanCandidateId:        planCandidateId,
		GooglePlaceId:          input.GooglePlaceID,
		CategoryNamesPreferred: &input.CategoriesPreferred,
		CategoryNamesDisliked:  &input.CategoriesDisliked,
		FreeTime:               input.FreeTime,
	})
	if err != nil {
		logger.Error("error while creating plan by google place id", zap.Error(err))
		return nil, fmt.Errorf("internal server error")
	}

	// 作成されたプランの保存
	if err := planCandidateService.SavePlans(ctx, plancandidate.SavePlansInput{
		PlanCandidateId:              planCandidateId,
		Plans:                        createPlanByGooglePlaceidResult.Plans,
		CategoryNamesPreferred:       &input.CategoriesPreferred,
		CategoryNamesRejected:        &input.CategoriesDisliked,
		FreeTime:                     input.FreeTime,
		CreateBasedOnCurrentLocation: false,
	}); err != nil {
		logger.Error("error while saving plans of plan candidate", zap.Error(err))
		return nil, fmt.Errorf("internal server error")
	}

	graphqlPlans := factory.PlansFromDomainModel(&createPlanByGooglePlaceidResult.Plans, &createPlanByGooglePlaceidResult.StartPlace.Location)
	return &model.CreatePlanByGooglePlaceIDOutput{
		PlanCandidate: &model.PlanCandidate{
			ID:                            planCandidateId,
			Plans:                         graphqlPlans,
			LikedPlaceIds:                 nil,
			CreatedBasedOnCurrentLocation: false,
		},
	}, nil
}

// ChangePlacesOrderInPlanCandidate is the resolver for the changePlacesOrderInPlanCandidate field.
func (r *mutationResolver) ChangePlacesOrderInPlanCandidate(ctx context.Context, input model.ChangePlacesOrderInPlanCandidateInput) (*model.ChangePlacesOrderInPlanCandidateOutput, error) {
	service, err := plancandidate.NewService(r.DB)
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
	logger, err := utils.NewLogger(utils.LoggerOption{Tag: "GraphQL"})
	if err != nil {
		log.Println("error while initializing logger: ", err)
		return nil, fmt.Errorf("internal server error")
	}

	service, err := plan.NewService(ctx, r.DB)
	if err != nil {
		logger.Error("error while initializing PlanService", zap.Error(err))
		return nil, fmt.Errorf("internal server error")
	}

	logger.Info(
		"SavePlanFromCandidate",
		zap.String("planCandidateId", input.Session),
		zap.String("planId", input.PlanID),
	)

	planSaved, err := service.SavePlanFromPlanCandidate(ctx, input.Session, input.PlanID, input.AuthToken)
	if err != nil {
		logger.Error("error while saving plan from plan candidate", zap.Error(err))
		return nil, fmt.Errorf("could not save plan")
	}

	graphqlPlan, err := factory.PlanFromDomainModel(*planSaved, nil)
	if err != nil {
		logger.Error("error while converting plan to graphql model", zap.Error(err))
		return nil, fmt.Errorf("internal server error")
	}

	return &model.SavePlanFromCandidateOutput{
		Plan: graphqlPlan,
	}, nil
}

// AddPlaceToPlanCandidateAfterPlace is the resolver for the addPlaceToPlanCandidateAfterPlace field.
func (r *mutationResolver) AddPlaceToPlanCandidateAfterPlace(ctx context.Context, input *model.AddPlaceToPlanCandidateAfterPlaceInput) (*model.AddPlaceToPlanCandidateAfterPlaceOutput, error) {
	s, err := plancandidate.NewService(r.DB)
	if err != nil {
		log.Println(fmt.Errorf("error while initizalizing PlanService: %v", err))
		return nil, fmt.Errorf("internal server error")
	}

	planInPlanCandidate, err := s.AddPlaceAfterPlace(ctx, input.PlanCandidateID, input.PlanID, input.PreviousPlaceID, input.PlaceID)
	if err != nil {
		log.Println(fmt.Errorf("error while adding place to plan candidate: %v", err))
		return nil, fmt.Errorf("could not add place to plan candidate")
	}

	planCandidate, err := s.FindPlanCandidate(ctx, input.PlanCandidateID)
	if err != nil {
		log.Println(fmt.Errorf("error while finding plan candidate: %v", err))
		return nil, fmt.Errorf("internal server error")
	}

	graphqlPlanInPlanCandidate, err := factory.PlanFromDomainModel(*planInPlanCandidate, planCandidate.MetaData.LocationStart)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("internal server error")
	}

	return &model.AddPlaceToPlanCandidateAfterPlaceOutput{
		Plan: graphqlPlanInPlanCandidate,
	}, nil
}

// DeletePlaceFromPlanCandidate is the resolver for the deletePlaceFromPlanCandidate field.
func (r *mutationResolver) DeletePlaceFromPlanCandidate(ctx context.Context, input model.DeletePlaceFromPlanCandidateInput) (*model.DeletePlaceFromPlanCandidateOutput, error) {
	s, err := plancandidate.NewService(r.DB)
	if err != nil {
		log.Println(fmt.Errorf("error while initizalizing PlanService: %v", err))
		return nil, fmt.Errorf("internal server error")
	}

	planUpdated, err := s.RemovePlaceFromPlan(ctx, input.PlanCandidateID, input.PlanID, input.PlaceID)
	if err != nil {
		log.Println(fmt.Errorf("error while deleting place from plan candidate: %v", err))
		return nil, fmt.Errorf("could not delete place from plan candidate")
	}

	planCandidate, err := s.FindPlanCandidate(ctx, input.PlanCandidateID)
	if err != nil {
		log.Println(fmt.Errorf("error while finding plan candidate: %v", err))
		return nil, fmt.Errorf("internal server error")
	}

	graphqlPlanInPlanCandidate, err := factory.PlanFromDomainModel(*planUpdated, planCandidate.MetaData.LocationStart)
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
	s, err := plancandidate.NewService(r.DB)
	if err != nil {
		log.Println(fmt.Errorf("error while initizalizing PlanService: %v", err))
		return nil, fmt.Errorf("internal server error")
	}

	plan, err := s.ReplacePlace(ctx, input.PlanCandidateID, input.PlanID, input.PlaceIDToRemove, input.PlaceIDToReplace)
	if err != nil {
		log.Println(fmt.Errorf("error while replacing place of plan candidate: %v", err))
		return nil, fmt.Errorf("could not replace place of plan candidate")
	}

	planCandidate, err := s.FindPlanCandidate(ctx, input.PlanCandidateID)
	if err != nil {
		log.Println(fmt.Errorf("error while finding plan candidate: %v", err))
		return nil, fmt.Errorf("internal server error")
	}

	graphqlPlanInPlanCandidate, err := factory.PlanFromDomainModel(*plan, planCandidate.MetaData.LocationStart)
	if err != nil {
		log.Printf("error while converting plan to graphql model: %v", err)
		return nil, fmt.Errorf("internal server error")
	}

	return &model.ReplacePlaceOfPlanCandidateOutput{
		PlanCandidateID: input.PlanCandidateID,
		Plan:            graphqlPlanInPlanCandidate,
	}, nil
}

// EditPlanTitleOfPlanCandidate is the resolver for the editPlanTitleOfPlanCandidate field.
func (r *mutationResolver) EditPlanTitleOfPlanCandidate(ctx context.Context, input model.EditPlanTitleOfPlanCandidateInput) (*model.EditPlanTitleOfPlanCandidateOutput, error) {
	panic(fmt.Errorf("not implemented: EditPlanTitleOfPlanCandidate - editPlanTitleOfPlanCandidate"))
}

// AutoReorderPlacesInPlanCandidate is the resolver for the autoReorderPlacesInPlanCandidate field.
func (r *mutationResolver) AutoReorderPlacesInPlanCandidate(ctx context.Context, input model.AutoReorderPlacesInPlanCandidateInput) (*model.AutoReorderPlacesInPlanCandidateOutput, error) {
	planCandidateService, err := plancandidate.NewService(r.DB)
	if err != nil {
		log.Println(fmt.Errorf("error while initizalizing PlanService: %v", err))
		return nil, fmt.Errorf("internal server error")
	}

	planUpdated, err := planCandidateService.AutoReorderPlaces(ctx, plancandidate.AutoReorderPlacesInput{
		PlanCandidateId: input.PlanCandidateID,
		PlanId:          input.PlanID,
	})
	if err != nil {
		log.Println(fmt.Errorf("error while auto reordering places in plan candidate: %v", err))
		return nil, fmt.Errorf("could not auto reorder places in plan candidate")
	}

	graphqlPlanInPlanCandidate, err := factory.PlanFromDomainModel(*planUpdated, nil)
	if err != nil {
		log.Printf("error while converting plan to graphql model: %v", err)
		return nil, fmt.Errorf("internal server error")
	}

	return &model.AutoReorderPlacesInPlanCandidateOutput{
		PlanCandidateID: input.PlanCandidateID,
		Plan:            graphqlPlanInPlanCandidate,
	}, nil
}

// LikeToPlaceInPlanCandidate is the resolver for the likeToPlaceInPlanCandidate field.
func (r *mutationResolver) LikeToPlaceInPlanCandidate(ctx context.Context, input model.LikeToPlaceInPlanCandidateInput) (*model.LikeToPlaceInPlanCandidateOutput, error) {
	planCandidateService, err := plancandidate.NewService(r.DB)
	if err != nil {
		log.Println(fmt.Errorf("error while initizalizing PlanService: %v", err))
		return nil, fmt.Errorf("internal server error")
	}

	planCandidateUpdated, err := planCandidateService.LikeToPlaceInPlanCandidate(ctx, input.PlanCandidateID, input.PlaceID, input.Like)
	if err != nil {
		log.Println(fmt.Errorf("error while liking to place in plan candidate: %v", err))
		return nil, fmt.Errorf("could not like to place in plan candidate")
	}

	graphqlPlanCandidate := factory.PlanCandidateFromDomainModel(planCandidateUpdated)

	if err != nil {
		log.Printf("error while converting plan candidate to graphql model: %v", err)
		return nil, fmt.Errorf("internal server error")
	}

	return &model.LikeToPlaceInPlanCandidateOutput{
		PlanCandidate: graphqlPlanCandidate,
	}, nil
}
