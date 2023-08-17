package firestore

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/firestore/entity"
)

const (
	collectionPlanCandidates = "plan_candidates"
)

type PlanCandidateFirestoreRepository struct {
	client *firestore.Client
}

func NewPlanCandidateRepository(ctx context.Context) (*PlanCandidateFirestoreRepository, error) {
	var options []option.ClientOption
	if os.Getenv("GCP_CREDENTIAL_FILE_PATH") != "" {
		options = append(options, option.WithCredentialsFile(os.Getenv("GCP_CREDENTIAL_FILE_PATH")))
	}

	client, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT_ID"), options...)
	if err != nil {
		return nil, fmt.Errorf("error while initializing firestore client: %v", err)
	}

	return &PlanCandidateFirestoreRepository{
		client: client,
	}, nil
}

func (p *PlanCandidateFirestoreRepository) Save(ctx context.Context, planCandidate *models.PlanCandidate) error {
	doc := p.doc(planCandidate.Id)
	if _, err := doc.Set(ctx, entity.ToPlanCandidateEntity(*planCandidate)); err != nil {
		return fmt.Errorf("error while saving plan candidate: %v", err)
	}
	return nil
}

func (p *PlanCandidateFirestoreRepository) Find(ctx context.Context, planCandidateId string) (*models.PlanCandidate, error) {
	doc := p.doc(planCandidateId)
	snapshot, err := doc.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("error while finding plan candidate: %v", err)
	}

	var planCandidateEntity entity.PlanCandidateEntity
	if err = snapshot.DataTo(&planCandidateEntity); err != nil {
		return nil, fmt.Errorf("error while converting snapshot to plan candidate entity: %v", err)
	}

	planCandidate := entity.FromPlanCandidateEntity(planCandidateEntity)
	return &planCandidate, nil
}

func (p *PlanCandidateFirestoreRepository) FindExpiredBefore(ctx context.Context, expiresAt time.Time) (*[]models.PlanCandidate, error) {
	var planCandidates []models.PlanCandidate
	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		query := p.collection().Where("expires_at", "<=", expiresAt)
		snapshots, err := tx.Documents(query).GetAll()
		if err != nil {
			return fmt.Errorf("error while getting all plan candidates: %v", err)
		}

		for _, snapshot := range snapshots {
			var planCandidateEntity entity.PlanCandidateEntity
			if err = snapshot.DataTo(&planCandidateEntity); err != nil {
				return fmt.Errorf("error while converting snapshot to plan candidate entity: %v", err)
			}

			planCandidate := entity.FromPlanCandidateEntity(planCandidateEntity)
			planCandidates = append(planCandidates, planCandidate)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("error while finding expired plan candidates: %v", err)
	}

	return &planCandidates, nil
}

func (p *PlanCandidateFirestoreRepository) AddPlan(
	ctx context.Context,
	planCandidateId string,
	plan *models.Plan,
) (*models.PlanCandidate, error) {
	var planCandidate models.PlanCandidate
	err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc := p.doc(planCandidateId)
		docRef, err := tx.Get(doc)
		if err != nil {
			return err
		}

		var planCandidateEntity entity.PlanCandidateEntity
		if err = docRef.DataTo(&planCandidateEntity); err != nil {
			return err
		}

		planEntity := entity.ToPlanInCandidateEntity(*plan)
		planCandidateEntity.Plans = append(planCandidateEntity.Plans, planEntity)

		if err := tx.Update(doc, []firestore.Update{
			{
				Path:  "plans",
				Value: planCandidateEntity.Plans,
			},
		}); err != nil {
			return err
		}

		planCandidate = entity.FromPlanCandidateEntity(planCandidateEntity)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error while adding plan to plan candidate: %v", err)
	}

	return &planCandidate, nil
}

func (p *PlanCandidateFirestoreRepository) UpdatePlacesOrder(ctx context.Context, planId string, planCandidateId string, placeIdsOrdered []string) (*models.Plan, error) {
	planCandidate, err := p.Find(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while finding plan candidate: %w\n", err)
	}

	if planCandidate == nil {
		return nil, fmt.Errorf("not found plan candidate[%s]\n", planCandidateId)
	}

	// planCandidateの中から指定のplanを探索
	var plan *models.Plan
	for _, p := range planCandidate.Plans {
		if p.Id == planId {
			plan = &p
			break
		}
	}
	if plan == nil {
		return nil, fmt.Errorf("not found plan[%s] in plan candidate[%s]", planId, planCandidate.Id)
	}

	// MOCK：並び替え・更新処理を実装

	return plan, nil
}

func (p *PlanCandidateFirestoreRepository) DeleteAll(ctx context.Context, planCandidateIds []string) error {
	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		for _, planCandidateId := range planCandidateIds {
			doc := p.doc(planCandidateId)
			if err := tx.Delete(doc); err != nil {
				return fmt.Errorf("error while deleting plan candidate[%s]: %v", planCandidateId, err)
			}
		}
		return nil
	}, firestore.MaxAttempts(3)); err != nil {
		return fmt.Errorf("error while deleting plan candidates: %v", err)
	}

	return nil
}

func (p *PlanCandidateFirestoreRepository) collection() *firestore.CollectionRef {
	return p.client.Collection(collectionPlanCandidates)
}

func (p *PlanCandidateFirestoreRepository) doc(id string) *firestore.DocumentRef {
	return p.collection().Doc(id)
}
