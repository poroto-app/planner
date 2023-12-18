package firestore

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"
	"sync"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/firestore/entity"
)

const (
	collectionPlaces           = "places"
	subCollectionGooglePlaces  = "google_places"
	subCollectionGoogleReviews = "google_place_reviews"
	subCollectionGooglePhotos  = "google_place_photos"
)

type PlaceRepository struct {
	client *firestore.Client
	logger *zap.Logger
}

func NewPlaceRepository(ctx context.Context) (*PlaceRepository, error) {
	var options []option.ClientOption
	if os.Getenv("GCP_CREDENTIAL_FILE_PATH") != "" {
		options = append(options, option.WithCredentialsFile(os.Getenv("GCP_CREDENTIAL_FILE_PATH")))
	}

	client, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT_ID"), options...)
	if err != nil {
		return nil, fmt.Errorf("error while initializing firestore client: %v", err)
	}

	logger, err := utils.NewLogger(utils.LoggerOption{
		Tag: "Firestore PlaceRepository",
	})
	if err != nil {
		return nil, fmt.Errorf("error while initializing logger: %v", err)
	}

	return &PlaceRepository{
		client: client,
		logger: logger,
	}, nil
}

func (p PlaceRepository) SavePlacesFromGooglePlace(ctx context.Context, googlePlace models.GooglePlace) (*models.Place, error) {
	var placeEntity entity.PlaceEntity
	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// すでに保存されている場合はそれを取得する
		savedPlaceEntity, err := p.findByGooglePlaceIdTx(tx, googlePlace.PlaceId)
		if err != nil {
			return fmt.Errorf("error while finding place by google place id: %v", err)
		}

		if savedPlaceEntity != nil {
			placeEntity = *savedPlaceEntity

			// TODO: サービスの部分で取得処理を書く（保存しただけなのに、Service内で取得していない値が入るのは気持ち悪い）
			// すでに保存されている Google Place の情報を取得する
			gp, err := p.fetchGooglePlace(ctx, placeEntity.Id)
			if err != nil {
				return fmt.Errorf("error while fetching google place: %v", err)
			}

			if gp != nil {
				p.logger.Debug(
					"Skip saving place because it is already saved",
					zap.String("placeId", placeEntity.Id),
					zap.String("googlePlaceId", googlePlace.PlaceId),
				)
				googlePlace = *gp
				return nil
			}
		}

		// 保存されていない場合は新規に保存する
		placeDoc := p.collectionPlaces().NewDoc()
		placeId := placeDoc.ID
		placeEntity = entity.NewPlaceEntityFromGooglePlace(placeId, googlePlace)
		if err := tx.Set(placeDoc, placeEntity); err != nil {
			return fmt.Errorf("error while saving place: %v", err)
		}

		// Google Place を保存する
		googlePlaceEntity := entity.GooglePlaceEntityFromGooglePlace(googlePlace)
		if err := tx.Set(p.docGooglePlace(placeId), googlePlaceEntity); err != nil {
			return fmt.Errorf("error while saving google place: %v", err)
		}

		// Place Detail を保存する
		if googlePlace.PlaceDetail != nil {
			// PhotoReferenceを保存する
			if err := p.saveGooglePhotoReferencesTx(tx, placeEntity.Id, googlePlace.PlaceDetail.PhotoReferences); err != nil {
				return fmt.Errorf("error while saving google place photos: %v", err)
			}

			// Reviewを保存する
			if err := p.saveGooglePlaceReviews(ctx, tx, placeEntity.Id, googlePlace.PlaceDetail.Reviews); err != nil {
				return fmt.Errorf("error while saving google place reviews: %v", err)
			}

			// 開店時間を更新する
			if err := p.updateOpeningHours(tx, placeEntity.Id, *googlePlace.PlaceDetail); err != nil {
				return fmt.Errorf("error while updating opening hours: %v", err)
			}
		} else if len(googlePlace.PhotoReferences) > 0 {
			// Nearby Searchで画像を取得している場合は保存する
			if err := p.saveGooglePhotoReferencesTx(tx, placeEntity.Id, googlePlace.PhotoReferences); err != nil {
				return fmt.Errorf("error while saving google place photos: %v", err)
			}
		}

		// Place Photo を保存する
		if googlePlace.Photos != nil {
			if err := p.saveGooglePhotosTx(ctx, tx, placeEntity.Id, googlePlace.PlaceId, *googlePlace.Photos); err != nil {
				return fmt.Errorf("error while saving google place photos: %v", err)
			}
		}

		return nil
	}, firestore.MaxAttempts(3)); err != nil {
		return nil, fmt.Errorf("error while running transaction: %v", err)
	}

	place := placeEntity.ToPlace(googlePlace)
	return &place, nil
}

func (p PlaceRepository) FindByLocation(ctx context.Context, location models.GeoLocation) ([]models.Place, error) {
	type findPlaceEntityByGeoHashResult struct {
		placeEntity *[]entity.PlaceEntity
		err         error
	}

	type fetchGooglePlaceResult struct {
		placeEntity entity.PlaceEntity
		googlePlace *models.GooglePlace
		err         error
	}

	// 各方向に 5km 以内の場所を取得
	geohashPrecision := 5
	geohashNeighbors := location.GeoHashOfNeighbors(uint(geohashPrecision))
	// TODO: 重複がないかを確認する
	geohashNeighbors = append(geohashNeighbors, location.GeoHash()[:geohashPrecision])

	ch := make(chan findPlaceEntityByGeoHashResult, len(geohashNeighbors))
	for _, geoHash := range geohashNeighbors {
		go func(ch chan<- findPlaceEntityByGeoHashResult, geoHash string) {
			query := p.collectionPlaces().Where("geohash", ">=", geoHash).Where("geohash", "<=", geoHash+"\uf8ff")
			query = query.Limit(50)
			iter := query.Documents(ctx)
			snapshots, err := iter.GetAll()
			if err != nil {
				ch <- findPlaceEntityByGeoHashResult{
					placeEntity: nil,
					err:         fmt.Errorf("error while getting place entities: %v", err),
				}
				return
			}

			var placeEntities []entity.PlaceEntity
			for _, snapshot := range snapshots {
				var placeEntity entity.PlaceEntity
				if err := snapshot.DataTo(&placeEntity); err != nil {
					ch <- findPlaceEntityByGeoHashResult{
						placeEntity: nil,
						err:         fmt.Errorf("error while converting snapshot to place entity: %v", err),
					}
					return
				}
				placeEntities = append(placeEntities, placeEntity)
			}

			ch <- findPlaceEntityByGeoHashResult{
				placeEntity: &placeEntities,
				err:         nil,
			}
		}(ch, geoHash)
	}

	var placeEntities []entity.PlaceEntity
	for i := 0; i < len(geohashNeighbors); i++ {
		result := <-ch
		if result.err != nil {
			return nil, result.err
		}
		if result.placeEntity != nil {
			placeEntities = append(placeEntities, *result.placeEntity...)
		}
	}

	// Places APIの検索結果を取得
	chGooglePlaces := make(chan fetchGooglePlaceResult, len(placeEntities))
	for _, placeEntity := range placeEntities {
		go func(ch chan<- fetchGooglePlaceResult, placeEntity entity.PlaceEntity) {
			googlePlace, err := p.fetchGooglePlace(ctx, placeEntity.Id)
			if err != nil {
				ch <- fetchGooglePlaceResult{
					placeEntity: placeEntity,
					googlePlace: nil,
					err:         fmt.Errorf("error while fetching google place: %v", err),
				}
				return
			}

			ch <- fetchGooglePlaceResult{
				placeEntity: placeEntity,
				googlePlace: googlePlace,
				err:         nil,
			}
		}(chGooglePlaces, placeEntity)
	}

	var places []models.Place
	for i := 0; i < len(placeEntities); i++ {
		result := <-chGooglePlaces
		if result.err != nil {
			return nil, result.err
		}
		if result.googlePlace != nil {
			place := result.placeEntity.ToPlace(*result.googlePlace)
			places = append(places, place)
		}
	}

	return places, nil
}

func (p PlaceRepository) FindByGooglePlaceID(ctx context.Context, googlePlaceID string) (*models.Place, error) {
	query := p.collectionPlaces().Where("google_place_id", "==", googlePlaceID).Limit(1)
	iter := query.Documents(ctx)
	doc, err := iter.Next()
	if errors.Is(err, iterator.Done) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error while iterating documents: %v", err)
	}

	var placeEntity entity.PlaceEntity
	if err := doc.DataTo(&placeEntity); err != nil {
		return nil, fmt.Errorf("error while converting doc to entity: %v", err)
	}

	googlePlace, err := p.fetchGooglePlace(ctx, placeEntity.Id)
	if err != nil {
		return nil, fmt.Errorf("error while fetching google place: %v", err)
	}

	place := placeEntity.ToPlace(*googlePlace)
	return &place, nil
}

func (p PlaceRepository) FindByPlanCandidateId(ctx context.Context, planCandidateId string) ([]models.Place, error) {
	// PlanCandidateを取得する
	snapshotPlanCandidate, err := p.client.Collection(collectionPlanCandidates).Doc(planCandidateId).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return nil, fmt.Errorf("plan candidate not found by id: %s", planCandidateId)
	}
	if err != nil {
		return nil, fmt.Errorf("error while getting plan candidate: %v", err)
	}

	var planCandidateEntity entity.PlanCandidateEntity
	if err := snapshotPlanCandidate.DataTo(&planCandidateEntity); err != nil {
		return nil, fmt.Errorf("error while converting snapshot to plan candidate entity: %v", err)
	}

	// 重複した場所が取得されないようにする
	placeIdsSearchedForPlanCandidate := array.StrArrayToSet(planCandidateEntity.PlaceIdsSearched)

	places, err := p.findByPlaceIds(ctx, placeIdsSearchedForPlanCandidate)
	if err != nil {
		return nil, fmt.Errorf("error while fetching places: %v", err)
	}

	return *places, nil
}

func (p PlaceRepository) SaveGooglePlacePhotos(ctx context.Context, googlePlaceId string, photos []models.GooglePlacePhoto) error {
	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// 事前に保存する画像が存在するかを確認する
		placeEntity, err := p.findByGooglePlaceIdTx(tx, googlePlaceId)
		if err != nil {
			return fmt.Errorf("error while finding place by google place id: %v", err)
		}
		if placeEntity == nil {
			return fmt.Errorf("place not found by google place id: %s", googlePlaceId)
		}

		// 画像を保存する
		if err := p.saveGooglePhotosTx(ctx, tx, placeEntity.Id, googlePlaceId, photos); err != nil {
			return fmt.Errorf("error while saving google place photos: %v", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error while running transaction: %v", err)
	}

	return nil
}

func (p PlaceRepository) SaveGooglePlaceDetail(ctx context.Context, googlePlaceId string, detail models.GooglePlaceDetail) error {
	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// 事前に要素が存在するかを確認する
		placeEntity, err := p.findByGooglePlaceIdTx(tx, googlePlaceId)
		if err != nil {
			return fmt.Errorf("error while finding place by google place id: %v", err)
		}
		if placeEntity == nil {
			return fmt.Errorf("place not found by google place id: %s", googlePlaceId)
		}

		// PhotoReferenceを保存する
		if err := p.saveGooglePhotoReferencesTx(tx, placeEntity.Id, detail.PhotoReferences); err != nil {
			return fmt.Errorf("error while saving google place photos: %v", err)
		}

		// Reviewを保存する
		if err := p.saveGooglePlaceReviews(ctx, tx, placeEntity.Id, detail.Reviews); err != nil {
			return fmt.Errorf("error while saving google place reviews: %v", err)
		}

		// 開店時間を更新する
		if err := p.updateOpeningHours(tx, placeEntity.Id, detail); err != nil {
			return fmt.Errorf("error while updating opening hours: %v", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error while saving google place detail: %v", err)
	}

	return nil
}

// findByPlaceIds placeIds で指定された複数の場所を取得する
// 　一つでも保存されていないものがあればエラーを返す
func (p PlaceRepository) findByPlaceIds(ctx context.Context, placeIds []string) (*[]models.Place, error) {
	chPlace := make(chan *models.Place, len(placeIds))
	chErr := make(chan error)
	defer close(chPlace)
	defer close(chErr)

	for _, placeId := range placeIds {
		go func(ch chan<- *models.Place, chErr chan<- error, placeId string) {
			place, err := p.findByPlaceId(ctx, placeId)
			if utils.HandleWrappedErrWithCh(ctx, chErr, err, fmt.Errorf("error while fetching place: %v", err)) {
				return
			}

			// 保存されていない場合はエラーを返す
			if place == nil {
				utils.HandleWrappedErrWithCh(ctx, chErr, fmt.Errorf("place not found by place id: %s", placeId), nil)
				return
			}

			utils.SendOrAbort(ctx, ch, place)
		}(chPlace, chErr, placeId)
	}

	var places []models.Place
	for i := 0; i < len(placeIds); i++ {
		select {
		case place := <-chPlace:
			places = append(places, *place)
		case err := <-chErr:
			return nil, err
		}
	}

	return &places, nil
}

func (p PlaceRepository) findByPlaceId(ctx context.Context, placeId string) (*models.Place, error) {
	chPlace := make(chan *entity.PlaceEntity, 1)
	chGooglePlace := make(chan *models.GooglePlace, 1)
	chErr := make(chan error)
	defer close(chPlace)
	defer close(chGooglePlace)
	defer close(chErr)

	asyncProcesses := []func(ctx context.Context){
		func(ctx context.Context) {
			// Placeを取得する
			var placeEntity entity.PlaceEntity
			snapshotPlace, err := p.docPlace(placeId).Get(ctx)
			if utils.HandleWrappedErrWithCh(ctx, chErr, err, fmt.Errorf("error while finding place by place id: %v", err)) {
				return
			}

			if err := snapshotPlace.DataTo(&placeEntity); utils.HandleWrappedErrWithCh(ctx, chErr, err, fmt.Errorf("error while converting snapshot to place entity: %v", err)) {
				return
			}

			utils.SendOrAbort(ctx, chPlace, &placeEntity)
		},
		func(ctx context.Context) {
			// Google Placeを取得する
			googlePlace, err := p.fetchGooglePlace(ctx, placeId)
			if utils.HandleWrappedErrWithCh(ctx, chErr, err, fmt.Errorf("error while fetching google place: %v", err)) {
				return
			}

			utils.SendOrAbort(ctx, chGooglePlace, googlePlace)
		},
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	chDone := make(chan bool)
	var wg sync.WaitGroup
	for _, asyncProcess := range asyncProcesses {
		wg.Add(1)
		go func(asyncProcess func(ctx context.Context)) {
			defer wg.Done()
			asyncProcess(ctx)
		}(asyncProcess)
	}

	go func() {
		wg.Wait()
		close(chDone)
	}()

	var placeEntity *entity.PlaceEntity
	var googlePlace *models.GooglePlace
Loop:
	for {
		select {
		case placeEntity = <-chPlace:
			if placeEntity == nil {
				return nil, fmt.Errorf("place is not found by place id: %s", placeId)
			}
		case googlePlace = <-chGooglePlace:
			if googlePlace == nil {
				return nil, fmt.Errorf("google place is not found by place id: %s", placeId)
			}
		case err := <-chErr:
			return nil, err
		case <-chDone:
			break Loop
		}
	}

	place := placeEntity.ToPlace(*googlePlace)
	return &place, nil
}

// fetchGooglePlace は placeId に紐づく Google Places API の検索結果を取得する
func (p PlaceRepository) fetchGooglePlace(ctx context.Context, placeId string) (*models.GooglePlace, error) {
	chGooglePlace := make(chan *entity.GooglePlaceEntity, 1)
	chReviews := make(chan *[]entity.GooglePlaceReviewEntity, 1)
	chPhotos := make(chan *[]entity.GooglePlacePhotoEntity, 1)
	chErr := make(chan error)
	defer close(chGooglePlace)
	defer close(chReviews)
	defer close(chPhotos)
	defer close(chErr)

	asyncProcesses := []func(){
		func() {
			// Google Placeを取得する
			snapshotGooglePlace, err := p.docGooglePlace(placeId).Get(ctx)
			if utils.HandleWrappedErrWithCh(ctx, chErr, err, fmt.Errorf("error while finding google place by place id: %v", err)) {
				return
			}

			var googlePlaceEntity entity.GooglePlaceEntity
			if err := snapshotGooglePlace.DataTo(&googlePlaceEntity); utils.HandleWrappedErrWithCh(ctx, chErr, err, fmt.Errorf("error while converting snapshot to google place entity: %v", err)) {
				return
			}

			utils.SendOrAbort(ctx, chGooglePlace, &googlePlaceEntity)
		},
		func() {
			// Reviewを取得する
			snapshotsReviews, err := p.subCollectionGooglePlaceReview(placeId).Documents(ctx).GetAll()
			if utils.HandleWrappedErrWithCh(ctx, chErr, err, fmt.Errorf("error while getting google place reviews: %v", err)) {
				return
			}

			var reviews []entity.GooglePlaceReviewEntity
			for _, snapshotReview := range snapshotsReviews {
				var review entity.GooglePlaceReviewEntity
				if err := snapshotReview.DataTo(&review); utils.HandleWrappedErrWithCh(ctx, chErr, err, fmt.Errorf("error while converting snapshot to google place review entity: %v", err)) {
					return
				}
				reviews = append(reviews, review)
			}

			utils.SendOrAbort(ctx, chReviews, &reviews)
		},
		func() {
			// Photoを取得する
			snapshotsPhotos, err := p.subCollectionGooglePlacePhoto(placeId).Documents(ctx).GetAll()
			if utils.HandleWrappedErrWithCh(ctx, chErr, err, fmt.Errorf("error while getting google place photos: %v", err)) {
				return
			}

			var photos []entity.GooglePlacePhotoEntity
			for _, snapshotPhoto := range snapshotsPhotos {
				var photo entity.GooglePlacePhotoEntity
				if err := snapshotPhoto.DataTo(&photo); utils.HandleWrappedErrWithCh(ctx, chErr, err, fmt.Errorf("error while converting snapshot to google place photo entity: %v", err)) {
					return
				}
				photos = append(photos, photo)
			}

			utils.SendOrAbort(ctx, chPhotos, &photos)
		},
	}

	chDone := make(chan bool)
	var wg sync.WaitGroup
	for _, asyncProcess := range asyncProcesses {
		wg.Add(1)
		go func(asyncProcess func()) {
			defer wg.Done()
			asyncProcess()
		}(asyncProcess)
	}

	go func() {
		wg.Wait()
		close(chDone)
	}()

	var googlePlaceEntity *entity.GooglePlaceEntity
	var reviewEntities *[]entity.GooglePlaceReviewEntity
	var photoEntities *[]entity.GooglePlacePhotoEntity
Loop:
	for {
		select {
		case googlePlaceEntity = <-chGooglePlace:
			if googlePlaceEntity == nil {
				return nil, fmt.Errorf("google place is nil")
			}
		case reviewEntities = <-chReviews:
			if reviewEntities == nil {
				return nil, fmt.Errorf("google place reviews is nil")
			}
		case photoEntities = <-chPhotos:
			if photoEntities == nil {
				return nil, fmt.Errorf("google place photos is nil")
			}
		case err := <-chErr:
			return nil, err
		case <-chDone:
			break Loop
		}
	}

	googlePlace := googlePlaceEntity.ToGooglePlace(photoEntities, reviewEntities)

	return &googlePlace, nil
}

func (p PlaceRepository) findByGooglePlaceIdTx(tx *firestore.Transaction, googlePlaceId string) (*entity.PlaceEntity, error) {
	query := p.collectionPlaces().Where("google_place_id", "==", googlePlaceId).Limit(1)
	iter := tx.Documents(query)
	doc, err := iter.Next()
	if errors.Is(err, iterator.Done) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error while iterating documents: %v", err)
	}

	var placeEntity entity.PlaceEntity
	if err := doc.DataTo(&placeEntity); err != nil {
		return nil, fmt.Errorf("error while converting doc to entity: %v", err)
	}

	return &placeEntity, nil
}

// saveGooglePlaceTx はGoogle Places APIから取得された複数の画像を同時に保存する
// 一枚でも保存できなかった場合はエラーを返す
func (p PlaceRepository) saveGooglePhotosTx(ctx context.Context, tx *firestore.Transaction, placeId string, googlePlaceId string, photos []models.GooglePlacePhoto) error {
	ch := make(chan *models.GooglePlacePhoto, len(photos))
	chErr := make(chan error)
	defer close(ch)
	defer close(chErr)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, photo := range photos {
		go func(ctx context.Context, tx *firestore.Transaction, ch chan<- *models.GooglePlacePhoto, googlePlaceId string, photo models.GooglePlacePhoto) {
			photoEntity := entity.GooglePlacePhotoEntityFromGooglePlacePhoto(photo)
			if err := tx.Set(p.subCollectionGooglePlacePhoto(placeId).Doc(photo.PhotoReference), photoEntity); utils.HandleWrappedErrWithCh(ctx, chErr, err, fmt.Errorf("error while saving google place photo: %v", err)) {
				return
			}

			utils.SendOrAbort(ctx, ch, &photo)
		}(ctx, tx, ch, googlePlaceId, photo)
	}

	for i := 0; i < len(photos); i++ {
		if photo := <-ch; photo == nil {
			return fmt.Errorf("error while saving google place photo: %v", photos[i])
		}
	}

	return nil
}

func (p PlaceRepository) saveGooglePhotoReferencesTx(tx *firestore.Transaction, placeId string, photoReferences []models.GooglePlacePhotoReference) error {
	ch := make(chan *models.GooglePlacePhotoReference, len(photoReferences))
	chErr := make(chan error)
	defer close(ch)
	defer close(chErr)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, photoReference := range photoReferences {
		go func(ctx context.Context, tx *firestore.Transaction, ch chan<- *models.GooglePlacePhotoReference, placeId string, photoReference models.GooglePlacePhotoReference) {
			p.logger.Info(
				"start saving google place photo references",
				zap.String("placeId", placeId),
				zap.String("photoReference", photoReference.PhotoReference),
				zap.String("width", fmt.Sprintf("%d", photoReference.Width)),
				zap.String("height", fmt.Sprintf("%d", photoReference.Height)),
			)
			photoEntity := entity.GooglePlacePhotoEntityFromGooglePhotoReference(photoReference)
			if err := tx.Set(p.subCollectionGooglePlacePhoto(placeId).Doc(photoReference.PhotoReference), photoEntity); utils.HandleWrappedErrWithCh(ctx, chErr, err, fmt.Errorf("error while saving google place photo reference: %v", err)) {
				return
			}

			if !utils.SendOrAbort(ctx, ch, &photoReference) {
				return
			}

			p.logger.Info(
				"successfully saved google place photo references",
				zap.String("placeId", placeId),
				zap.String("photoReference", photoReference.PhotoReference),
			)
		}(ctx, tx, ch, placeId, photoReference)
	}

	for range photoReferences {
		select {
		case photoReference := <-ch:
			if photoReference == nil {
				return fmt.Errorf("error while saving google place photo reference: %v", photoReference)
			}
		case err := <-chErr:
			return err
		}
	}

	return nil
}

func (p PlaceRepository) saveGooglePlaceReviews(ctx context.Context, tx *firestore.Transaction, placeId string, reviews []models.GooglePlaceReview) error {
	ch := make(chan *models.GooglePlaceReview, len(reviews))
	chErr := make(chan error)
	defer close(ch)
	defer close(chErr)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, review := range reviews {
		go func(ctx context.Context, tx *firestore.Transaction, ch chan<- *models.GooglePlaceReview, placeId string, review models.GooglePlaceReview) {
			// 重複したレビューが保存されないように ID を MD5(Time+Text+Language) で生成する
			// AuthorName 等は頻繁に変更される可能性があるため、IDには含めない
			hashContent := fmt.Sprintf("%d-%s-%s", review.Time, utils.StrEmptyIfNil(review.Text), utils.StrEmptyIfNil(review.Language))
			id := fmt.Sprintf("%x", md5.Sum([]byte(hashContent)))

			reviewEntity := entity.GooglePlaceReviewEntityFromGooglePlaceReview(review)
			if err := tx.Set(p.subCollectionGooglePlaceReview(placeId).Doc(id), reviewEntity); utils.HandleWrappedErrWithCh(ctx, chErr, err, fmt.Errorf("error while saving google place review: %v", err)) {
				return
			}

			utils.SendOrAbort(ctx, ch, &review)
		}(ctx, tx, ch, placeId, review)
	}

	for range reviews {
		select {
		case review := <-ch:
			if review == nil {
				return fmt.Errorf("error while saving google place review: %v", review)
			}
		case err := <-chErr:
			return err
		}
	}

	return nil
}

func (p PlaceRepository) updateOpeningHours(tx *firestore.Transaction, placeId string, placeDetail models.GooglePlaceDetail) error {
	// 開店時間を更新
	if placeDetail.OpeningHours != nil {
		openingHoursEntity := entity.GooglePlaceOpeningsEntityFromGooglePlaceOpeningHours(*placeDetail.OpeningHours)
		if err := tx.Update(p.docGooglePlace(placeId), []firestore.Update{
			{
				Path:  "opening_hours",
				Value: openingHoursEntity,
			},
			{
				Path:  "updated_at",
				Value: firestore.ServerTimestamp,
			},
		}); err != nil {
			return fmt.Errorf("error while saving google place detail: %v", err)
		}
	}

	return nil
}

func (p PlaceRepository) collectionPlaces() *firestore.CollectionRef {
	return p.client.Collection(collectionPlaces)
}

func (p PlaceRepository) docPlace(placeId string) *firestore.DocumentRef {
	return p.client.Collection(collectionPlaces).Doc(placeId)
}

func (p PlaceRepository) docGooglePlace(placeId string) *firestore.DocumentRef {
	return p.client.Collection(collectionPlaces).Doc(placeId).Collection(subCollectionGooglePlaces).Doc("v1")
}

func (p PlaceRepository) subCollectionGooglePlaceReview(placeId string) *firestore.CollectionRef {
	return p.client.Collection(collectionPlaces).Doc(placeId).Collection(subCollectionGoogleReviews)
}

func (p PlaceRepository) subCollectionGooglePlacePhoto(placeId string) *firestore.CollectionRef {
	return p.client.Collection(collectionPlaces).Doc(placeId).Collection(subCollectionGooglePhotos)
}
