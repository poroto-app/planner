package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/firestore/entity"
)

const (
	collectionPlans = "plans"
)

type PlanRepository struct {
	client          *firestore.Client
	placeRepository *PlaceRepository
}

func NewPlanRepository(ctx context.Context) (*PlanRepository, error) {
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
		return nil, fmt.Errorf("error while initializing place repository: %v", err)
	}

	return &PlanRepository{
		client:          client,
		placeRepository: placeRepository,
	}, nil
}

func (p *PlanRepository) Save(ctx context.Context, plan *models.Plan) error {
	doc := p.doc(plan.Id)
	if _, err := doc.Set(ctx, entity.NewPlanEntityFromPlan(*plan)); err != nil {
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

	places, err := p.placeRepository.findByPlaceIds(ctx, planEntity.PlaceIds)
	if err != nil {
		return nil, fmt.Errorf("error while finding places: %v", err)
	}

	plan, err := planEntity.ToPlan(*places)
	if err != nil {
		return nil, fmt.Errorf("error while converting plan entity to plan: %v", err)
	}

	return plan, nil
}

func (p *PlanRepository) FindByAuthorId(ctx context.Context, authorId string) (*[]models.Plan, error) {
	collection := p.collection()
	query := collection.Where("author_id", "==", authorId).OrderBy("updated_at", firestore.Desc)

	snapshots, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("error while getting plans: %v", err)
	}

	// Plan を取得
	var planEntities []entity.PlanEntity
	for _, snapshot := range snapshots {
		var planEntity entity.PlanEntity
		if err = snapshot.DataTo(&planEntity); err != nil {
			return nil, fmt.Errorf("error while converting snapshot to plan entity: %v", err)
		}
	}

	// Plan に含まれる Place　を一括で取得
	places, err := p.findPlacesFromPlaceEntities(ctx, planEntities)
	if err != nil {
		return nil, fmt.Errorf("error while finding places: %v", err)
	}

	var plans []models.Plan
	for _, planEntity := range planEntities {
		plan, err := planEntity.ToPlan(*places)
		if err != nil {
			return nil, fmt.Errorf("error while converting plan entity to plan: %v", err)
		}
		plans = append(plans, *plan)
	}

	return &plans, nil
}

// SortedByCreatedAt created_atで降順に並べたPlanを取得する
// queryCursor(リストの最後の [models.Plan] のID)が指定されている場合は、そのcursor以降のPlanを取得する
func (p *PlanRepository) SortedByCreatedAt(ctx context.Context, queryCursor *string, limit int) (*[]models.Plan, error) {
	// PlanEntityを取得
	query := p.collection().OrderBy("created_at", firestore.Desc).OrderBy("id", firestore.Desc)
	if queryCursor != nil {
		plan, err := p.find(ctx, *queryCursor)
		if err != nil {
			return nil, fmt.Errorf("error while getting plans: %v", err)
		}

		if plan == nil {
			return nil, fmt.Errorf("plan of queryCursor(%s) was not found", *queryCursor)
		}

		query = query.StartAfter(plan.CreatedAt, plan.Id)
	}
	query = query.Limit(limit)

	snapshots, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("error while getting plans: %v", err)
	}

	planEntities := make([]entity.PlanEntity, len(snapshots))
	for i, snapshot := range snapshots {
		var planEntity entity.PlanEntity
		if err = snapshot.DataTo(&planEntity); err != nil {
			return nil, fmt.Errorf("error while converting snapshot to plan entity: %v", err)
		}

		planEntities[i] = planEntity
	}

	// Plan に含まれる Place　を一括で取得
	places, err := p.findPlacesFromPlaceEntities(ctx, planEntities)
	if err != nil {
		return nil, fmt.Errorf("error while finding places: %v", err)
	}

	var plans []models.Plan
	for _, planEntity := range planEntities {
		plan, err := planEntity.ToPlan(*places)
		if err != nil {
			return nil, fmt.Errorf("error while converting plan entity to plan: %v", err)
		}
		plans = append(plans, *plan)
	}

	return &plans, nil
}

// SortedByLocation location で指定した地点を含む20kmのGeoHashに含まれるプランを取得する
// TODO: 現在地から近い順に取得できるようにする
// TODO: レビュー等の指標に基づいて取得できるようにする
func (p *PlanRepository) SortedByLocation(ctx context.Context, geoLocation models.GeoLocation, queryCursor *string, limit int) (*[]models.Plan, *string, error) {

	// 20km圏内のPlanを取得する
	// SEE: https://en.wikipedia.org/wiki/Geohash#Digits_and_precision_in_km
	geohash := geoLocation.GeoHash()
	query := p.collection().Where("geohash", ">=", geohash[:4]).Where("geohash", "<=", geohash[:4]+"\uf8ff")

	query = query.OrderBy("geohash", firestore.Desc)

	if queryCursor != nil {
		query = query.StartAfter(*queryCursor)
	}

	query = query.Limit(limit)
	snapshots, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, nil, fmt.Errorf("error while getting plans: %v", err)
	}

	var nextQueryCursor *string
	planEntities := make([]entity.PlanEntity, len(snapshots))
	for i, snapshot := range snapshots {
		var planEntity entity.PlanEntity
		if err = snapshot.DataTo(&planEntity); err != nil {
			return nil, nil, fmt.Errorf("error while converting snapshot to plan entity: %v", err)
		}

		planEntities[i] = planEntity
		nextQueryCursor = planEntity.GeoHash
	}

	// Plan に含まれる Place　を一括で取得
	places, err := p.findPlacesFromPlaceEntities(ctx, planEntities)
	if err != nil {
		return nil, nil, fmt.Errorf("error while finding places: %v", err)
	}

	var plans []models.Plan
	for _, planEntity := range planEntities {
		plan, err := planEntity.ToPlan(*places)
		if err != nil {
			return nil, nil, fmt.Errorf("error while converting plan entity to plan: %v", err)
		}
		plans = append(plans, *plan)
	}

	return &plans, nextQueryCursor, nil
}

func (p *PlanRepository) findPlacesFromPlaceEntities(ctx context.Context, placeEntities []entity.PlanEntity) (*[]models.Place, error) {
	var placeIds []string
	for _, plan := range placeEntities {
		placeIds = append(placeIds, plan.PlaceIds...)
	}

	places, err := p.placeRepository.findByPlaceIds(ctx, array.StrArrayToSet(placeIds))
	if err != nil {
		return nil, fmt.Errorf("error while finding places: %v", err)
	}

	return places, nil
}

func (p *PlanRepository) collection() *firestore.CollectionRef {
	return p.client.Collection(collectionPlans)
}

func (p *PlanRepository) doc(id string) *firestore.DocumentRef {
	return p.collection().Doc(id)
}
