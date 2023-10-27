package firestore

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/firestore/entity"
)

const (
	collectionPlanCandidates = "plan_candidates"
	collectionMetaData       = "meta_data"
	subCollectionPlans       = "plans"
)

type PlanCandidateFirestoreRepository struct {
	client                         *firestore.Client
	placeInPlanCandidateRepository *PlaceInPlanCandidateRepository
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

	placeInPlanCandidateRepository, err := NewPlaceInPlanCandidateRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing placeInPlanCandidateRepository: %v", err)
	}

	return &PlanCandidateFirestoreRepository{
		client:                         client,
		placeInPlanCandidateRepository: placeInPlanCandidateRepository,
	}, nil
}

func (p *PlanCandidateFirestoreRepository) Save(ctx context.Context, planCandidate *models.PlanCandidate) error {
	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc := p.doc(planCandidate.Id)
		if err := tx.Set(doc, entity.ToPlanCandidateEntity(*planCandidate)); err != nil {
			return fmt.Errorf("error while saving plan candidate: %v", err)
		}

		docMetaData := p.docMetaDataV1(planCandidate.Id)
		if err := tx.Set(docMetaData, entity.ToPlanCandidateMetaDataV1Entity(planCandidate.MetaData)); err != nil {
			return fmt.Errorf("error while saving plan candidate meta data: %v", err)
		}

		for _, plan := range planCandidate.Plans {
			if err := tx.Set(p.subCollectionPlans(planCandidate.Id).Doc(plan.Id), entity.ToPlanInCandidateEntity(plan)); err != nil {
				return fmt.Errorf("error while saving plan in plan candidate: %v", err)
			}
		}

		return nil
	}); err != nil {
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

	// プラン候補メタデータを取得
	docMetaData := p.docMetaDataV1(planCandidateId)
	snapshotMetaData, err := docMetaData.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("error while finding plan candidate meta data: %v", err)
	}

	var planCandidateMetaDataEntity entity.PlanCandidateMetaDataV1Entity
	if err = snapshotMetaData.DataTo(&planCandidateMetaDataEntity); err != nil {
		return nil, fmt.Errorf("error while converting snapshot to plan candidate meta data entity: %v", err)
	}

	// プランを取得
	snapshotPlans, err := p.subCollectionPlans(planCandidateId).Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("error while fetching plans: %v", err)
	}

	var plans []entity.PlanInCandidateEntity
	for _, snapshotPlan := range snapshotPlans {
		var plan entity.PlanInCandidateEntity
		if err = snapshotPlan.DataTo(&plan); err != nil {
			return nil, fmt.Errorf("error while converting snapshot to plan entity: %v", err)
		}
		plans = append(plans, plan)
	}

	//　検索された場所を取得
	places, err := p.placeInPlanCandidateRepository.FindByPlanCandidateId(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching places: %v", err)
	}

	planCandidate := entity.FromPlanCandidateEntity(planCandidateEntity, planCandidateMetaDataEntity, plans, *places)

	return &planCandidate, nil
}

func (p *PlanCandidateFirestoreRepository) FindExpiredBefore(ctx context.Context, expiresAt time.Time) (*[]string, error) {
	var planCandidateIds []string

	query := p.collection().Where("expires_at", "<=", expiresAt)
	docIter := query.Documents(ctx)
	for {
		doc, err := docIter.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			}
			return nil, fmt.Errorf("error while iterating plan candidate: %v", err)
		}
		planCandidateIds = append(planCandidateIds, doc.Ref.ID)
	}

	return &planCandidateIds, nil
}

// TODO: errorだけを返すようにする
func (p *PlanCandidateFirestoreRepository) AddPlan(
	ctx context.Context,
	planCandidateId string,
	plan *models.Plan,
) (*models.PlanCandidate, error) {
	err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// プランを保存
		docPlan := p.subCollectionPlans(planCandidateId).Doc(plan.Id)
		if err := tx.Set(docPlan, entity.ToPlanInCandidateEntity(*plan)); err != nil {
			return fmt.Errorf("error while saving plan[%s] in plan candidate[%s]: %v", plan.Id, planCandidateId, err)
		}

		// 候補の最後に追加
		docPlanCandidate := p.doc(planCandidateId)
		if err := tx.Update(docPlanCandidate, []firestore.Update{
			{
				Path:  "plan_ids",
				Value: firestore.ArrayUnion(plan.Id),
			},
		}); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error while adding plan to plan candidate: %v", err)
	}

	planCandidateUpdated, err := p.Find(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while finding plan candidate: %v", err)
	}
	return planCandidateUpdated, nil
}

func (p *PlanCandidateFirestoreRepository) AddPlaceToPlan(ctx context.Context, planCandidateId string, planId string, place models.Place) error {
	err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc := p.subCollectionPlans(planCandidateId).Doc(planId)
		if err := tx.Update(doc, []firestore.Update{
			{
				Path:  "place_ids_ordered",
				Value: firestore.ArrayUnion(place.Id),
			},
			{
				Path:  "updated_at",
				Value: firestore.ServerTimestamp,
			},
		}); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error while adding place to plan candidate: %v", err)
	}

	return nil
}

func (p *PlanCandidateFirestoreRepository) RemovePlaceFromPlan(ctx context.Context, planCandidateId string, planId string, placeId string) error {
	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc := p.subCollectionPlans(planCandidateId).Doc(planId)
		if err := tx.Update(doc, []firestore.Update{
			{
				Path:  "place_ids_ordered",
				Value: firestore.ArrayRemove(placeId),
			},
			{
				Path:  "updated_at",
				Value: firestore.ServerTimestamp,
			},
		}); err != nil {
			return fmt.Errorf("error while updating plan candidate[%s]: %v", planCandidateId, err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error while removing place from plan candidate: %v", err)
	}

	return nil
}

func (p *PlanCandidateFirestoreRepository) UpdatePlacesOrder(ctx context.Context, planId string, planCandidateId string, placeIdsOrdered []string) (*models.Plan, error) {
	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc := p.subCollectionPlans(planCandidateId).Doc(planId)
		snapshot, err := tx.Get(doc)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return fmt.Errorf("plan[%s] not found in plan candidate[%s]", planId, planCandidateId)
			}

			return fmt.Errorf("error while getting plan[%s] in plan candidate[%s]: %v", planId, planCandidateId, err)
		}

		var planInCandidateEntity entity.PlanInCandidateEntity
		if err = snapshot.DataTo(&planInCandidateEntity); err != nil {
			return fmt.Errorf("error while converting snapshot to plan entity: %v", err)
		}

		if len(placeIdsOrdered) != len(planInCandidateEntity.PlaceIdsOrdered) {
			return fmt.Errorf("the length of placeIdsOrdered is invalid")
		}

		if err := tx.Update(doc, []firestore.Update{
			{
				Path:  "place_ids_ordered",
				Value: placeIdsOrdered,
			},
			{
				Path:  "updated_at",
				Value: firestore.ServerTimestamp,
			},
		}); err != nil {
			return fmt.Errorf("error while updating places order: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("error while updating places order: %v", err)
	}

	planCandidate, err := p.Find(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while finding plan candidate: %v", err)
	}

	var plan *models.Plan
	for _, p := range planCandidate.Plans {
		if p.Id == planId {
			plan = &p
			break
		}
	}

	return plan, nil
}

func (p *PlanCandidateFirestoreRepository) ReplacePlace(ctx context.Context, planCandidateId, planId, placeIdToBeReplaced string, placeToReplace models.Place) error {
	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc := p.subCollectionPlans(planCandidateId).Doc(planId)
		snapshot, err := tx.Get(doc)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return fmt.Errorf("plan[%s] not found in plan candidate[%s]", planId, planCandidateId)
			}

			return fmt.Errorf("error while getting plan[%s] in plan candidate[%s]: %v", planId, planCandidateId, err)
		}

		var planInCandidateEntity entity.PlanInCandidateEntity
		if err = snapshot.DataTo(&planInCandidateEntity); err != nil {
			return fmt.Errorf("error while converting snapshot to plan entity: %v", err)
		}

		for i, placeIdOrdered := range planInCandidateEntity.PlaceIdsOrdered {
			if placeIdOrdered == placeIdToBeReplaced {
				planInCandidateEntity.PlaceIdsOrdered[i] = placeToReplace.Id
				break
			}
		}

		if err := tx.Update(doc, []firestore.Update{
			{
				Path:  "place_ids_ordered",
				Value: planInCandidateEntity.PlaceIdsOrdered,
			},
			{
				Path:  "updated_at",
				Value: firestore.ServerTimestamp,
			},
		}); err != nil {
			return fmt.Errorf("error while updating plan candidate[%s]: %v", planCandidateId, err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error while replacing place: %v", err)
	}

	return nil
}

func (p *PlanCandidateFirestoreRepository) DeleteAll(ctx context.Context, planCandidateIds []string) error {
	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		for _, planCandidateId := range planCandidateIds {
			// プラン候補メタデータを削除
			docMetadata := p.docMetaDataV1(planCandidateId)
			if err := tx.Delete(docMetadata); err != nil {
				return fmt.Errorf("error while deleting plan candidate meta data[%s]: %v", planCandidateId, err)
			}

			// プランを削除
			docIter := tx.DocumentRefs(p.subCollectionPlans(planCandidateId))
			for {
				doc, err := docIter.Next()
				if err != nil {
					if errors.Is(err, iterator.Done) {
						break
					}
					return fmt.Errorf("error while iterating plans: %v", err)
				}

				if err := tx.Delete(doc); err != nil {
					return fmt.Errorf("error while deleting plan[%s]: %v", doc.ID, err)
				}
			}

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

func (p *PlanCandidateFirestoreRepository) subCollectionPlans(planCandidateId string) *firestore.CollectionRef {
	return p.doc(planCandidateId).Collection(subCollectionPlans)
}

func (p *PlanCandidateFirestoreRepository) docMetaDataV1(id string) *firestore.DocumentRef {
	return p.doc(id).Collection(collectionMetaData).Doc("v1")
}
