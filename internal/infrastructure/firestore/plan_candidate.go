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
	collectionMetaData       = "meta_data"
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
	doc := p.doc(planCandidate.Id)
	if _, err := doc.Set(ctx, entity.ToPlanCandidateEntity(*planCandidate)); err != nil {
		return fmt.Errorf("error while saving plan candidate: %v", err)
	}

	docMetaData := p.docMetaDataV1(planCandidate.Id)
	if _, err := docMetaData.Set(ctx, entity.ToPlanCandidateMetaDataV1Entity(planCandidate.MetaData)); err != nil {
		return fmt.Errorf("error while saving plan candidate meta data: %v", err)
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

	places, err := p.placeInPlanCandidateRepository.FindByPlanCandidateId(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching places: %v", err)
	}

	planCandidate := entity.FromPlanCandidateEntity(planCandidateEntity, planCandidateMetaDataEntity, *places)

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

			// この関数は削除でのみもちいるため、メタデータの取得はできていなくても良い
			// TODO: validateが入らないようにする
			planCandidate := entity.FromPlanCandidateEntity(planCandidateEntity, entity.PlanCandidateMetaDataV1Entity{}, nil)
			planCandidates = append(planCandidates, planCandidate)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("error while finding expired plan candidates: %v", err)
	}

	return &planCandidates, nil
}

// TODO: errorだけを返すようにする
func (p *PlanCandidateFirestoreRepository) AddPlan(
	ctx context.Context,
	planCandidateId string,
	plan *models.Plan,
) (*models.PlanCandidate, error) {
	var planCandidate models.PlanCandidate
	err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc := p.doc(planCandidateId)
		_, err := tx.Get(doc)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return fmt.Errorf("plan candidate[%s] not found", planCandidateId)
			}
			return fmt.Errorf("error while getting plan candidate[%s]: %v", planCandidateId, err)
		}

		planEntity := entity.ToPlanInCandidateEntity(*plan)
		if err := tx.Update(doc, []firestore.Update{
			{
				Path:  "plans",
				Value: firestore.ArrayUnion(planEntity),
			},
		}); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error while adding plan to plan candidate: %v", err)
	}

	docMetaData := p.docMetaDataV1(planCandidateId)
	snapshotMetaData, err := docMetaData.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while finding plan candidate meta data: %v", err)
	}

	var planCandidateMetaDataEntity entity.PlanCandidateMetaDataV1Entity
	if err = snapshotMetaData.DataTo(&planCandidateMetaDataEntity); err != nil {
		return nil, fmt.Errorf("error while converting snapshot to plan candidate meta data entity: %v", err)
	}

	docPlanCandidate := p.doc(planCandidateId)
	snapshotPlanCandidate, err := docPlanCandidate.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while finding plan candidate: %v", err)
	}

	var planCandidateEntity entity.PlanCandidateEntity
	if err = snapshotPlanCandidate.DataTo(&planCandidateEntity); err != nil {
		return nil, fmt.Errorf("error while converting snapshot to plan candidate entity: %v", err)
	}

	places, err := p.placeInPlanCandidateRepository.FindByPlanCandidateId(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching places: %v", err)
	}

	planCandidate = entity.FromPlanCandidateEntity(planCandidateEntity, planCandidateMetaDataEntity, *places)

	return &planCandidate, nil
}

func (p *PlanCandidateFirestoreRepository) AddPlaceToPlan(ctx context.Context, planCandidateId string, planId string, place models.Place) error {
	err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc := p.doc(planCandidateId)
		snapshot, err := tx.Get(doc)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return fmt.Errorf("plan candidate[%s] not found", planCandidateId)
			}
			return fmt.Errorf("error while getting plan candidate[%s]: %v", planCandidateId, err)
		}

		var planCandidateEntity entity.PlanCandidateEntity
		if err = snapshot.DataTo(&planCandidateEntity); err != nil {
			return fmt.Errorf("error while converting snapshot to plan candidate entity: %v", err)
		}

		var planIndex *int
		for i, plan := range planCandidateEntity.Plans {
			if plan.Id == planId {
				idx := i
				planIndex = &idx
				break
			}
		}

		if planIndex == nil {
			return fmt.Errorf("plan[%s] not found in plan candidate[%s]", planId, planCandidateId)
		}

		// TODO: Planを別のサブコレクションに分ける
		planEntityToUpdate := planCandidateEntity.Plans[*planIndex]
		planEntityToUpdate.PlaceIdsOrdered = append(planEntityToUpdate.PlaceIdsOrdered, place.Id)
		planCandidateEntity.Plans[*planIndex] = planEntityToUpdate

		if err := tx.Update(doc, []firestore.Update{
			{
				Path:  "plans",
				Value: planCandidateEntity.Plans,
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
		doc := p.doc(planCandidateId)
		snapshot, err := tx.Get(doc)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return fmt.Errorf("plan candidate[%s] not found", planCandidateId)
			}
			return fmt.Errorf("error while getting plan candidate[%s]: %v", planCandidateId, err)
		}

		var planCandidateEntity entity.PlanCandidateEntity
		if err = snapshot.DataTo(&planCandidateEntity); err != nil {
			return fmt.Errorf("error while converting snapshot to plan candidate entity: %v", err)
		}

		docMetaData := p.docMetaDataV1(planCandidateId)
		snapshotMetaData, err := tx.Get(docMetaData)
		if err != nil {
			return fmt.Errorf("error while getting plan candidate meta data: %v", err)
		}

		var planCandidateMetaDataEntity entity.PlanCandidateMetaDataV1Entity
		if err = snapshotMetaData.DataTo(&planCandidateMetaDataEntity); err != nil {
			return fmt.Errorf("error while converting snapshot to plan candidate meta data entity: %v", err)
		}

		// TODO: 直接削除するだけで十分な設計にする（Placeを別で持つ）
		var planIndex *int
		for i, plan := range planCandidateEntity.Plans {
			if plan.Id == planId {
				idx := i
				planIndex = &idx
				break
			}
		}
		if planIndex == nil {
			return fmt.Errorf("plan[%s] not found in plan candidate[%s]", planId, planCandidateId)
		}

		planEntityToUpdate := planCandidateEntity.Plans[*planIndex]
		// placeIdsOrdered から削除
		var placeIdsOrdered []string
		for _, placeIdOrdered := range planEntityToUpdate.PlaceIdsOrdered {
			if placeIdOrdered != placeId {
				placeIdsOrdered = append(placeIdsOrdered, placeIdOrdered)
			}
		}

		planEntityToUpdate.PlaceIdsOrdered = placeIdsOrdered
		planCandidateEntity.Plans[*planIndex] = planEntityToUpdate

		if err := tx.Update(doc, []firestore.Update{
			{
				Path:  "plans",
				Value: planCandidateEntity.Plans,
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

func (p *PlanCandidateFirestoreRepository) ReplacePlace(ctx context.Context, planCandidateId, planId, placeIdToBeReplaced string, placeToReplace models.Place) error {
	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// PlanCandidateを取得
		doc := p.doc(planCandidateId)
		snapshot, err := tx.Get(doc)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return fmt.Errorf("plan candidate[%s] not found", planCandidateId)
			}
		}

		var planCandidateEntity entity.PlanCandidateEntity
		if err = snapshot.DataTo(&planCandidateEntity); err != nil {
			return fmt.Errorf("error while converting snapshot to plan candidate entity: %v", err)
		}

		// Planを取得
		var planIndex *int
		for i, plan := range planCandidateEntity.Plans {
			if plan.Id == planId {
				idx := i
				planIndex = &idx
				break
			}
		}
		if planIndex == nil {
			return fmt.Errorf("plan[%s] not found in plan candidate[%s]", planId, planCandidateId)
		}

		// 順番を更新
		var placeIndex *int
		for i, placeId := range planCandidateEntity.Plans[*planIndex].PlaceIdsOrdered {
			if placeId == placeIdToBeReplaced {
				idx := i
				placeIndex = &idx
				break
			}
		}
		if placeIndex == nil {
			return fmt.Errorf("place[%s] not found in plan[%s] in plan candidate[%s]", placeIdToBeReplaced, planId, planCandidateId)
		}

		planCandidateEntity.Plans[*planIndex].PlaceIdsOrdered[*placeIndex] = placeToReplace.Id

		if err := tx.Update(doc, []firestore.Update{
			{
				Path:  "plans",
				Value: planCandidateEntity.Plans,
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

func (p *PlanCandidateFirestoreRepository) docMetaDataV1(id string) *firestore.DocumentRef {
	return p.doc(id).Collection(collectionMetaData).Doc("v1")
}
