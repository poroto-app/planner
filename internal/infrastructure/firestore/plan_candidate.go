package firestore

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/utils"

	"poroto.app/poroto/planner/internal/domain/array"

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
	client          *firestore.Client
	PlaceRepository *PlaceRepository
	logger          *zap.Logger
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

	placeRepository, err := NewPlaceRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing placeRepository: %v", err)
	}

	logger, err := utils.NewLogger(utils.LoggerOption{
		Tag: "Firestore PlanCandidateRepository",
	})
	if err != nil {
		return nil, fmt.Errorf("error while initializing logger: %v", err)
	}

	return &PlanCandidateFirestoreRepository{
		client:          client,
		PlaceRepository: placeRepository,
		logger:          logger,
	}, nil
}

func (p *PlanCandidateFirestoreRepository) Create(ctx context.Context, planCandidateId string, expiresAt time.Time) error {
	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		planCandidateEntity := entity.PlanCandidateEntity{
			Id:        planCandidateId,
			ExpiresAt: expiresAt,
		}
		if err := tx.Set(p.doc(planCandidateId), planCandidateEntity); err != nil {
			return fmt.Errorf("error while saving plan candidate: %v", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error while saving plan candidate: %v", err)
	}

	return nil
}

func (p *PlanCandidateFirestoreRepository) Find(ctx context.Context, planCandidateId string) (*models.PlanCandidate, error) {
	chPlanCandidate := make(chan *entity.PlanCandidateEntity, 1)
	chMetaData := make(chan *entity.PlanCandidateMetaDataV1Entity, 1)
	chPlans := make(chan *[]entity.PlanInCandidateEntity, 1)
	chPlaces := make(chan *[]models.Place, 1)
	chErr := make(chan error)
	chDone := make(chan bool)

	var wg sync.WaitGroup

	// PlanCandidateEntityを取得
	wg.Add(1)
	go func(ctx context.Context, ch chan *entity.PlanCandidateEntity, planCandidateId string) {
		defer wg.Done()
		p.logger.Info(
			"start fetching plan candidate",
			zap.String("planCandidateId", planCandidateId),
		)

		snapshot, err := p.doc(planCandidateId).Get(ctx)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				chPlanCandidate <- nil
				return
			}

			chErr <- fmt.Errorf("error while finding plan candidate: %v", err)
			return
		}

		var planCandidateEntity entity.PlanCandidateEntity
		if err = snapshot.DataTo(&planCandidateEntity); err != nil {
			chErr <- fmt.Errorf("error while converting snapshot to plan candidate entity: %v", err)
		}

		chPlanCandidate <- &planCandidateEntity
		p.logger.Info(
			"successfully fetched plan candidate",
			zap.String("planCandidateId", planCandidateId),
		)
	}(ctx, chPlanCandidate, planCandidateId)

	// メタデータを取得
	wg.Add(1)
	go func(ctx context.Context, ch chan *entity.PlanCandidateMetaDataV1Entity, planCandidateId string) {
		defer wg.Done()
		p.logger.Info(
			"start fetching plan candidate meta data",
			zap.String("planCandidateId", planCandidateId),
		)

		snapshot, err := p.docMetaDataV1(planCandidateId).Get(ctx)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				ch <- nil
				return
			}

			chErr <- fmt.Errorf("error while finding plan candidate meta data: %v", err)
			return
		}

		var planCandidateMetaDataEntity entity.PlanCandidateMetaDataV1Entity
		if err = snapshot.DataTo(&planCandidateMetaDataEntity); err != nil {
			chErr <- fmt.Errorf("error while converting snapshot to plan candidate meta data entity: %v", err)
		}

		ch <- &planCandidateMetaDataEntity
		p.logger.Info(
			"successfully fetched plan candidate meta data",
			zap.String("planCandidateId", planCandidateId),
		)
	}(ctx, chMetaData, planCandidateId)

	wg.Add(1)
	go func(ctx context.Context, chPlans chan *[]entity.PlanInCandidateEntity, chPlaces chan *[]models.Place, planCandidateId string) {
		defer wg.Done()

		// プラン候補に含まれるプランを取得
		p.logger.Info(
			"start fetching plans of plan candidate",
			zap.String("planCandidateId", planCandidateId),
		)
		snapshotPlans, err := p.subCollectionPlans(planCandidateId).Documents(ctx).GetAll()
		if err != nil {
			chErr <- fmt.Errorf("error while fetching plans: %v", err)
			return
		}

		var plans []entity.PlanInCandidateEntity
		for _, snapshotPlan := range snapshotPlans {
			var plan entity.PlanInCandidateEntity
			if err = snapshotPlan.DataTo(&plan); err != nil {
				chErr <- fmt.Errorf("error while converting snapshot to plan in plan candidate: %v", err)
				return
			}
			plans = append(plans, plan)
		}

		chPlans <- &plans
		p.logger.Info(
			"successfully fetched plans of plan candidate",
			zap.String("planCandidateId", planCandidateId),
		)

		//　プランに含まれる場所を取得
		p.logger.Info(
			"start fetching places of plan candidate",
			zap.String("planCandidateId", planCandidateId),
		)
		var placeIdsInPlan []string
		for _, plan := range plans {
			placeIdsInPlan = append(placeIdsInPlan, plan.PlaceIdsOrdered...)
		}
		places, err := p.PlaceRepository.findByPlaceIds(ctx, array.StrArrayToSet(placeIdsInPlan))
		if err != nil {
			chErr <- fmt.Errorf("error while fetching places: %v", err)
			return
		}

		chPlaces <- places
		p.logger.Info(
			"successfully fetched places of plan candidate",
			zap.String("planCandidateId", planCandidateId),
		)
	}(ctx, chPlans, chPlaces, planCandidateId)

	go func() {
		wg.Wait()
		p.logger.Info(
			"end fetching plan candidate",
			zap.String("planCandidateId", planCandidateId),
		)
		close(chDone)
	}()

	var planCandidateEntity *entity.PlanCandidateEntity
	var planCandidateMetaDataEntity *entity.PlanCandidateMetaDataV1Entity
	var plans *[]entity.PlanInCandidateEntity
	var places *[]models.Place
Loop:
	for {
		select {
		case planCandidateEntity = <-chPlanCandidate:
			if planCandidateEntity == nil {
				return nil, nil
			}
		case planCandidateMetaDataEntity = <-chMetaData:
			if planCandidateMetaDataEntity == nil {
				return nil, fmt.Errorf("plan candidate meta data not found: %s", planCandidateId)
			}
		case plans = <-chPlans:
			if plans == nil {
				return nil, fmt.Errorf("plans not found: %s", planCandidateId)
			}
		case places = <-chPlaces:
			if places == nil {
				return nil, fmt.Errorf("places not found: %s", planCandidateId)
			}
		case err := <-chErr:
			return nil, err
		case <-chDone:
			break Loop
		}
	}

	if planCandidateEntity == nil {
		return nil, nil
	}

	if planCandidateMetaDataEntity == nil {
		return nil, fmt.Errorf("plan candidate meta data not found: %s", planCandidateId)
	}

	if plans == nil {
		return nil, fmt.Errorf("plans not found: %s", planCandidateId)
	}

	if places == nil {
		return nil, fmt.Errorf("places not found: %s", planCandidateId)
	}

	planCandidate := planCandidateEntity.ToPlanCandidate(*planCandidateMetaDataEntity, *plans, *places)
	return &planCandidate, nil
}

func (p *PlanCandidateFirestoreRepository) FindPlan(ctx context.Context, planCandidateId string, planId string) (*models.Plan, error) {
	doc := p.subCollectionPlans(planCandidateId).Doc(planId)

	snapshot, err := doc.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("error while finding plan: %v", err)
	}

	var planInCandidateEntity entity.PlanInCandidateEntity
	if err = snapshot.DataTo(&planInCandidateEntity); err != nil {
		return nil, fmt.Errorf("error while converting snapshot to plan entity: %v", err)
	}

	places, err := p.PlaceRepository.findByPlaceIds(ctx, planInCandidateEntity.PlaceIdsOrdered)
	if err != nil {
		return nil, fmt.Errorf("error while fetching places: %v", err)
	}

	plan, err := planInCandidateEntity.ToPlan(*places)
	if err != nil {
		return nil, fmt.Errorf("error while converting plan entity to plan: %v", err)
	}

	return plan, nil
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

func (p *PlanCandidateFirestoreRepository) AddSearchedPlacesForPlanCandidate(ctx context.Context, planCandidateId string, placeIds []string) error {
	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// 事前に要素が存在するかを確認する
		docPlanCandidate := p.client.Collection(collectionPlanCandidates).Doc(planCandidateId)
		snapshotPlanCandidate, err := tx.Get(docPlanCandidate)
		if status.Code(err) == codes.NotFound {
			return fmt.Errorf("plan candidate not found by id: %s", planCandidateId)
		}
		if err != nil {
			return fmt.Errorf("error while getting plan candidate: %v", err)
		}

		var planCandidateEntity entity.PlanCandidateEntity
		if err := snapshotPlanCandidate.DataTo(&planCandidateEntity); err != nil {
			return fmt.Errorf("error while converting snapshot to plan candidate entity: %v", err)
		}

		// 重複した場所が取得されないようにする
		placeIdsSearched := planCandidateEntity.PlaceIdsSearched
		placeIdsSearched = append(placeIdsSearched, placeIds...)
		placeIdsSearched = array.StrArrayToSet(placeIdsSearched)

		// 更新する
		if err := tx.Update(docPlanCandidate, []firestore.Update{
			{
				Path:  "place_ids_searched",
				Value: placeIdsSearched,
			},
			{
				Path:  "updated_at",
				Value: firestore.ServerTimestamp,
			},
		}); err != nil {
			return fmt.Errorf("error while updating plan candidate: %v", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error while running transaction: %v", err)
	}

	return nil
}

func (p *PlanCandidateFirestoreRepository) AddPlan(ctx context.Context, planCandidateId string, plans ...models.Plan) error {
	err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// プランを保存
		for _, plan := range plans {
			docPlan := p.subCollectionPlans(planCandidateId).Doc(plan.Id)
			if err := tx.Set(docPlan, entity.ToPlanInCandidateEntity(plan)); err != nil {
				return fmt.Errorf("error while saving plan[%s] in plan candidate[%s]: %v", plan.Id, planCandidateId, err)
			}
		}

		var planIds []interface{}
		for _, plan := range plans {
			planIds = append(planIds, plan.Id)
		}

		// 候補の最後に追加
		docPlanCandidate := p.doc(planCandidateId)
		if err := tx.Update(docPlanCandidate, []firestore.Update{
			{
				Path:  "plan_ids",
				Value: firestore.ArrayUnion(planIds...),
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
		return fmt.Errorf("error while adding plan to plan candidate: %v", err)
	}

	return nil
}

func (p *PlanCandidateFirestoreRepository) AddPlaceToPlan(ctx context.Context, planCandidateId string, planId string, previousPlaceId string, place models.Place) error {
	err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
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

		var placeIdsOrderd []string
		for _, id := range planInCandidateEntity.PlaceIdsOrdered {
			placeIdsOrderd = append(placeIdsOrderd, id)
			if id != previousPlaceId {
				continue
			}
			placeIdsOrderd = append(placeIdsOrderd, place.Id)
		}

		if err := tx.Update(doc, []firestore.Update{
			{
				Path:  "place_ids_ordered",
				Value: placeIdsOrderd,
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

func (p *PlanCandidateFirestoreRepository) UpdatePlanCandidateMetaData(ctx context.Context, planCandidateId string, meta models.PlanCandidateMetaData) error {
	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc := p.docMetaDataV1(planCandidateId)
		if err := tx.Set(doc, entity.ToPlanCandidateMetaDataV1Entity(meta)); err != nil {
			return fmt.Errorf("error while saving plan candidate meta data: %v", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error while saving plan candidate meta data: %v", err)
	}

	return nil
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
	for _, planCandidateId := range planCandidateIds {
		p.logger.Info(
			"start deleting plan candidate",
			zap.String("planCandidateId", planCandidateId),
		)
		if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
			// プランを削除
			// DocumentRefsは内部で参照を行う
			// TransactionのルールではWriteの後にReadを行うことはできないため、削除処理の最初におこなう
			p.logger.Info(
				"start deleting plans",
				zap.String("planCandidateId", planCandidateId),
			)
			docIter := tx.DocumentRefs(p.subCollectionPlans(planCandidateId))
			for {
				if docIter == nil {
					p.logger.Info(
						"docIter of plan is nil",
						zap.String("planCandidateId", planCandidateId),
					)
					break
				}

				doc, err := docIter.Next()
				if err != nil {
					if errors.Is(err, iterator.Done) {
						break
					}
					return fmt.Errorf("error while iterating plans: %v", err)
				}

				p.logger.Info(
					"start deleting plan",
					zap.String("planId", doc.ID),
					zap.String("planCandidateId", planCandidateId),
				)
				if err := tx.Delete(doc); err != nil {
					return fmt.Errorf("error while deleting plan[%s]: %v", doc.ID, err)
				}
				p.logger.Info(
					"successfully deleted plan",
					zap.String("planId", doc.ID),
					zap.String("planCandidateId", planCandidateId),
				)
			}

			// プラン候補メタデータを削除
			p.logger.Info(
				"start deleting plan candidate meta data",
				zap.String("planCandidateId", planCandidateId),
			)
			docMetadata := p.docMetaDataV1(planCandidateId)
			if err := tx.Delete(docMetadata); err != nil {
				return fmt.Errorf("error while deleting plan candidate meta data[%s]: %v", planCandidateId, err)
			}
			p.logger.Info(
				"successfully deleted plan candidate meta data",
				zap.String("planCandidateId", planCandidateId),
			)

			// プラン候補を削除
			p.logger.Info(
				"start deleting plan candidate",
				zap.String("planCandidateId", planCandidateId),
			)
			doc := p.doc(planCandidateId)
			if err := tx.Delete(doc); err != nil {
				return fmt.Errorf("error while deleting plan candidate[%s]: %v", planCandidateId, err)
			}
			p.logger.Info(
				"successfully deleted plan candidate",
				zap.String("planCandidateId", planCandidateId),
			)

			return nil
		}, firestore.MaxAttempts(3)); err != nil {
			return fmt.Errorf("error while deleting plan candidates: %v", err)
		}
	}

	return nil
}

func (p *PlanCandidateFirestoreRepository) UpdateLikeToPlaceInPlanCandidate(ctx context.Context, planCandidateId string, placeId string, like bool) error {
	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// Placeの取得
		place, err := p.PlaceRepository.findByPlaceId(ctx, placeId)
		if err != nil {
			return fmt.Errorf("error while fetching place: %v", err)
		}

		// PlanCandidateの取得
		docPlanCandidate := p.client.Collection(collectionPlanCandidates).Doc(planCandidateId)
		snapshotPlanCandidate, err := tx.Get(docPlanCandidate)
		if status.Code(err) == codes.NotFound {
			return fmt.Errorf("plan candidate not found by id: %s", planCandidateId)
		}
		if err != nil {
			return fmt.Errorf("error while getting plan candidate: %v", err)
		}

		var planCandidateEntity entity.PlanCandidateEntity
		if err := snapshotPlanCandidate.DataTo(&planCandidateEntity); err != nil {
			return fmt.Errorf("error while converting snapshot to plan candidate entity: %v", err)
		}

		// すでにLikeしている場合は、Likeを取り消し
		if place.LikeCount > 0 && array.IsContain(planCandidateEntity.LikedPlaceIds, place.Id) {
			place.LikeCount -= 1
			for i, id := range planCandidateEntity.LikedPlaceIds {
				if id == place.Id {
					planCandidateEntity.LikedPlaceIds = append(planCandidateEntity.LikedPlaceIds[:i], planCandidateEntity.LikedPlaceIds[i+1:]...)
					break
				}
			}
		} else {
			// まだLikeされていない場合は、Likeを追加
			place.LikeCount += 1
			planCandidateEntity.LikedPlaceIds = append(planCandidateEntity.LikedPlaceIds, place.Id)
		}

		// PlanCandidateを更新する
		if err := tx.Update(docPlanCandidate, []firestore.Update{
			{
				Path:  "liked_place_ids",
				Value: planCandidateEntity.LikedPlaceIds,
			},
			{
				Path:  "updated_at",
				Value: firestore.ServerTimestamp,
			},
		}); err != nil {
			return fmt.Errorf("error while updating plan candidate: %v", err)
		}

		// Placeを更新する
		if err := tx.Update(p.PlaceRepository.docPlace(placeId), []firestore.Update{
			{
				Path:  "like_count",
				Value: place.LikeCount,
			},
			{
				Path:  "updated_at",
				Value: firestore.ServerTimestamp,
			},
		}); err != nil {
			return fmt.Errorf("error while updating place: %v", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error while running transaction: %v", err)
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
