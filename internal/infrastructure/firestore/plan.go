package firestore

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/firestore/entity"
)

const (
	collectionPlans = "plans"
)

type PlanRepository struct {
	client *firestore.Client
}

func NewPlanRepository(ctx context.Context) (*PlanRepository, error) {
	client, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT_ID"), option.WithCredentialsFile("secrets/google-credential.json"))
	if err != nil {
		return nil, fmt.Errorf("error while initializing firestore client: %v", err)
	}

	return &PlanRepository{client: client}, nil
}

func (p *PlanRepository) Save(ctx context.Context, plan *models.Plan) error {
	doc := p.doc(plan.Id)
	if _, err := doc.Set(ctx, entity.ToPlanEntity(*plan)); err != nil {
		return fmt.Errorf("error while saving plan: %v", err)
	}
	return nil
}

func (p *PlanRepository) find(ctx context.Context, planId string) (*entity.PlanEntity, error) {
	doc := p.doc(planId)
	snapshot, err := doc.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("error while finding plan: %v", err)
	}

	var planEntity entity.PlanEntity
	if err = snapshot.DataTo(&planEntity); err != nil {
		return nil, fmt.Errorf("error while converting snapshot to plan entity: %v", err)
	}

	return &planEntity, nil
}

func (p *PlanRepository) Find(ctx context.Context, planId string) (*models.Plan, error) {
	planEntity, err := p.find(ctx, planId)
	if err != nil {
		return nil, err
	}

	if planEntity == nil {
		return nil, nil
	}

	plan := entity.FromPlanEntity(*planEntity)
	return &plan, nil
}

func (p *PlanRepository) collection() *firestore.CollectionRef {
	return p.client.Collection(collectionPlans)
}

func (p *PlanRepository) doc(id string) *firestore.DocumentRef {
	return p.collection().Doc(id)
}
