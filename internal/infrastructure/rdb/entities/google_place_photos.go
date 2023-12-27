// Code generated by SQLBoiler 4.15.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package entities

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// GooglePlacePhoto is an object representing the database table.
type GooglePlacePhoto struct {
	ID             string    `boil:"id" json:"id" toml:"id" yaml:"id"`
	PhotoReference string    `boil:"photo_reference" json:"photo_reference" toml:"photo_reference" yaml:"photo_reference"`
	Width          int       `boil:"width" json:"width" toml:"width" yaml:"width"`
	Height         int       `boil:"height" json:"height" toml:"height" yaml:"height"`
	URL            string    `boil:"url" json:"url" toml:"url" yaml:"url"`
	CreatedAt      null.Time `boil:"created_at" json:"created_at,omitempty" toml:"created_at" yaml:"created_at,omitempty"`
	UpdatedAt      null.Time `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`

	R *googlePlacePhotoR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L googlePlacePhotoL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var GooglePlacePhotoColumns = struct {
	ID             string
	PhotoReference string
	Width          string
	Height         string
	URL            string
	CreatedAt      string
	UpdatedAt      string
}{
	ID:             "id",
	PhotoReference: "photo_reference",
	Width:          "width",
	Height:         "height",
	URL:            "url",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

var GooglePlacePhotoTableColumns = struct {
	ID             string
	PhotoReference string
	Width          string
	Height         string
	URL            string
	CreatedAt      string
	UpdatedAt      string
}{
	ID:             "google_place_photos.id",
	PhotoReference: "google_place_photos.photo_reference",
	Width:          "google_place_photos.width",
	Height:         "google_place_photos.height",
	URL:            "google_place_photos.url",
	CreatedAt:      "google_place_photos.created_at",
	UpdatedAt:      "google_place_photos.updated_at",
}

// Generated where

var GooglePlacePhotoWhere = struct {
	ID             whereHelperstring
	PhotoReference whereHelperstring
	Width          whereHelperint
	Height         whereHelperint
	URL            whereHelperstring
	CreatedAt      whereHelpernull_Time
	UpdatedAt      whereHelpernull_Time
}{
	ID:             whereHelperstring{field: "`google_place_photos`.`id`"},
	PhotoReference: whereHelperstring{field: "`google_place_photos`.`photo_reference`"},
	Width:          whereHelperint{field: "`google_place_photos`.`width`"},
	Height:         whereHelperint{field: "`google_place_photos`.`height`"},
	URL:            whereHelperstring{field: "`google_place_photos`.`url`"},
	CreatedAt:      whereHelpernull_Time{field: "`google_place_photos`.`created_at`"},
	UpdatedAt:      whereHelpernull_Time{field: "`google_place_photos`.`updated_at`"},
}

// GooglePlacePhotoRels is where relationship names are stored.
var GooglePlacePhotoRels = struct {
	PhotoReferenceGooglePlacePhotoReference string
}{
	PhotoReferenceGooglePlacePhotoReference: "PhotoReferenceGooglePlacePhotoReference",
}

// googlePlacePhotoR is where relationships are stored.
type googlePlacePhotoR struct {
	PhotoReferenceGooglePlacePhotoReference *GooglePlacePhotoReference `boil:"PhotoReferenceGooglePlacePhotoReference" json:"PhotoReferenceGooglePlacePhotoReference" toml:"PhotoReferenceGooglePlacePhotoReference" yaml:"PhotoReferenceGooglePlacePhotoReference"`
}

// NewStruct creates a new relationship struct
func (*googlePlacePhotoR) NewStruct() *googlePlacePhotoR {
	return &googlePlacePhotoR{}
}

func (r *googlePlacePhotoR) GetPhotoReferenceGooglePlacePhotoReference() *GooglePlacePhotoReference {
	if r == nil {
		return nil
	}
	return r.PhotoReferenceGooglePlacePhotoReference
}

// googlePlacePhotoL is where Load methods for each relationship are stored.
type googlePlacePhotoL struct{}

var (
	googlePlacePhotoAllColumns            = []string{"id", "photo_reference", "width", "height", "url", "created_at", "updated_at"}
	googlePlacePhotoColumnsWithoutDefault = []string{"photo_reference", "width", "height", "url"}
	googlePlacePhotoColumnsWithDefault    = []string{"id", "created_at", "updated_at"}
	googlePlacePhotoPrimaryKeyColumns     = []string{"id"}
	googlePlacePhotoGeneratedColumns      = []string{}
)

type (
	// GooglePlacePhotoSlice is an alias for a slice of pointers to GooglePlacePhoto.
	// This should almost always be used instead of []GooglePlacePhoto.
	GooglePlacePhotoSlice []*GooglePlacePhoto
	// GooglePlacePhotoHook is the signature for custom GooglePlacePhoto hook methods
	GooglePlacePhotoHook func(context.Context, boil.ContextExecutor, *GooglePlacePhoto) error

	googlePlacePhotoQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	googlePlacePhotoType                 = reflect.TypeOf(&GooglePlacePhoto{})
	googlePlacePhotoMapping              = queries.MakeStructMapping(googlePlacePhotoType)
	googlePlacePhotoPrimaryKeyMapping, _ = queries.BindMapping(googlePlacePhotoType, googlePlacePhotoMapping, googlePlacePhotoPrimaryKeyColumns)
	googlePlacePhotoInsertCacheMut       sync.RWMutex
	googlePlacePhotoInsertCache          = make(map[string]insertCache)
	googlePlacePhotoUpdateCacheMut       sync.RWMutex
	googlePlacePhotoUpdateCache          = make(map[string]updateCache)
	googlePlacePhotoUpsertCacheMut       sync.RWMutex
	googlePlacePhotoUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var googlePlacePhotoAfterSelectHooks []GooglePlacePhotoHook

var googlePlacePhotoBeforeInsertHooks []GooglePlacePhotoHook
var googlePlacePhotoAfterInsertHooks []GooglePlacePhotoHook

var googlePlacePhotoBeforeUpdateHooks []GooglePlacePhotoHook
var googlePlacePhotoAfterUpdateHooks []GooglePlacePhotoHook

var googlePlacePhotoBeforeDeleteHooks []GooglePlacePhotoHook
var googlePlacePhotoAfterDeleteHooks []GooglePlacePhotoHook

var googlePlacePhotoBeforeUpsertHooks []GooglePlacePhotoHook
var googlePlacePhotoAfterUpsertHooks []GooglePlacePhotoHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *GooglePlacePhoto) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlacePhotoAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *GooglePlacePhoto) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlacePhotoBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *GooglePlacePhoto) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlacePhotoAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *GooglePlacePhoto) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlacePhotoBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *GooglePlacePhoto) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlacePhotoAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *GooglePlacePhoto) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlacePhotoBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *GooglePlacePhoto) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlacePhotoAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *GooglePlacePhoto) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlacePhotoBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *GooglePlacePhoto) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlacePhotoAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddGooglePlacePhotoHook registers your hook function for all future operations.
func AddGooglePlacePhotoHook(hookPoint boil.HookPoint, googlePlacePhotoHook GooglePlacePhotoHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		googlePlacePhotoAfterSelectHooks = append(googlePlacePhotoAfterSelectHooks, googlePlacePhotoHook)
	case boil.BeforeInsertHook:
		googlePlacePhotoBeforeInsertHooks = append(googlePlacePhotoBeforeInsertHooks, googlePlacePhotoHook)
	case boil.AfterInsertHook:
		googlePlacePhotoAfterInsertHooks = append(googlePlacePhotoAfterInsertHooks, googlePlacePhotoHook)
	case boil.BeforeUpdateHook:
		googlePlacePhotoBeforeUpdateHooks = append(googlePlacePhotoBeforeUpdateHooks, googlePlacePhotoHook)
	case boil.AfterUpdateHook:
		googlePlacePhotoAfterUpdateHooks = append(googlePlacePhotoAfterUpdateHooks, googlePlacePhotoHook)
	case boil.BeforeDeleteHook:
		googlePlacePhotoBeforeDeleteHooks = append(googlePlacePhotoBeforeDeleteHooks, googlePlacePhotoHook)
	case boil.AfterDeleteHook:
		googlePlacePhotoAfterDeleteHooks = append(googlePlacePhotoAfterDeleteHooks, googlePlacePhotoHook)
	case boil.BeforeUpsertHook:
		googlePlacePhotoBeforeUpsertHooks = append(googlePlacePhotoBeforeUpsertHooks, googlePlacePhotoHook)
	case boil.AfterUpsertHook:
		googlePlacePhotoAfterUpsertHooks = append(googlePlacePhotoAfterUpsertHooks, googlePlacePhotoHook)
	}
}

// One returns a single googlePlacePhoto record from the query.
func (q googlePlacePhotoQuery) One(ctx context.Context, exec boil.ContextExecutor) (*GooglePlacePhoto, error) {
	o := &GooglePlacePhoto{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "entities: failed to execute a one query for google_place_photos")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all GooglePlacePhoto records from the query.
func (q googlePlacePhotoQuery) All(ctx context.Context, exec boil.ContextExecutor) (GooglePlacePhotoSlice, error) {
	var o []*GooglePlacePhoto

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "entities: failed to assign all query results to GooglePlacePhoto slice")
	}

	if len(googlePlacePhotoAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all GooglePlacePhoto records in the query.
func (q googlePlacePhotoQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "entities: failed to count google_place_photos rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q googlePlacePhotoQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "entities: failed to check if google_place_photos exists")
	}

	return count > 0, nil
}

// PhotoReferenceGooglePlacePhotoReference pointed to by the foreign key.
func (o *GooglePlacePhoto) PhotoReferenceGooglePlacePhotoReference(mods ...qm.QueryMod) googlePlacePhotoReferenceQuery {
	queryMods := []qm.QueryMod{
		qm.Where("`photo_reference` = ?", o.PhotoReference),
	}

	queryMods = append(queryMods, mods...)

	return GooglePlacePhotoReferences(queryMods...)
}

// LoadPhotoReferenceGooglePlacePhotoReference allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (googlePlacePhotoL) LoadPhotoReferenceGooglePlacePhotoReference(ctx context.Context, e boil.ContextExecutor, singular bool, maybeGooglePlacePhoto interface{}, mods queries.Applicator) error {
	var slice []*GooglePlacePhoto
	var object *GooglePlacePhoto

	if singular {
		var ok bool
		object, ok = maybeGooglePlacePhoto.(*GooglePlacePhoto)
		if !ok {
			object = new(GooglePlacePhoto)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeGooglePlacePhoto)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeGooglePlacePhoto))
			}
		}
	} else {
		s, ok := maybeGooglePlacePhoto.(*[]*GooglePlacePhoto)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeGooglePlacePhoto)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeGooglePlacePhoto))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &googlePlacePhotoR{}
		}
		args = append(args, object.PhotoReference)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &googlePlacePhotoR{}
			}

			for _, a := range args {
				if a == obj.PhotoReference {
					continue Outer
				}
			}

			args = append(args, obj.PhotoReference)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`google_place_photo_references`),
		qm.WhereIn(`google_place_photo_references.photo_reference in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load GooglePlacePhotoReference")
	}

	var resultSlice []*GooglePlacePhotoReference
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice GooglePlacePhotoReference")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for google_place_photo_references")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for google_place_photo_references")
	}

	if len(googlePlacePhotoReferenceAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.PhotoReferenceGooglePlacePhotoReference = foreign
		if foreign.R == nil {
			foreign.R = &googlePlacePhotoReferenceR{}
		}
		foreign.R.PhotoReferenceGooglePlacePhotos = append(foreign.R.PhotoReferenceGooglePlacePhotos, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.PhotoReference == foreign.PhotoReference {
				local.R.PhotoReferenceGooglePlacePhotoReference = foreign
				if foreign.R == nil {
					foreign.R = &googlePlacePhotoReferenceR{}
				}
				foreign.R.PhotoReferenceGooglePlacePhotos = append(foreign.R.PhotoReferenceGooglePlacePhotos, local)
				break
			}
		}
	}

	return nil
}

// SetPhotoReferenceGooglePlacePhotoReference of the googlePlacePhoto to the related item.
// Sets o.R.PhotoReferenceGooglePlacePhotoReference to related.
// Adds o to related.R.PhotoReferenceGooglePlacePhotos.
func (o *GooglePlacePhoto) SetPhotoReferenceGooglePlacePhotoReference(ctx context.Context, exec boil.ContextExecutor, insert bool, related *GooglePlacePhotoReference) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE `google_place_photos` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, []string{"photo_reference"}),
		strmangle.WhereClause("`", "`", 0, googlePlacePhotoPrimaryKeyColumns),
	)
	values := []interface{}{related.PhotoReference, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.PhotoReference = related.PhotoReference
	if o.R == nil {
		o.R = &googlePlacePhotoR{
			PhotoReferenceGooglePlacePhotoReference: related,
		}
	} else {
		o.R.PhotoReferenceGooglePlacePhotoReference = related
	}

	if related.R == nil {
		related.R = &googlePlacePhotoReferenceR{
			PhotoReferenceGooglePlacePhotos: GooglePlacePhotoSlice{o},
		}
	} else {
		related.R.PhotoReferenceGooglePlacePhotos = append(related.R.PhotoReferenceGooglePlacePhotos, o)
	}

	return nil
}

// GooglePlacePhotos retrieves all the records using an executor.
func GooglePlacePhotos(mods ...qm.QueryMod) googlePlacePhotoQuery {
	mods = append(mods, qm.From("`google_place_photos`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`google_place_photos`.*"})
	}

	return googlePlacePhotoQuery{q}
}

// FindGooglePlacePhoto retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindGooglePlacePhoto(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*GooglePlacePhoto, error) {
	googlePlacePhotoObj := &GooglePlacePhoto{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `google_place_photos` where `id`=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, googlePlacePhotoObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "entities: unable to select from google_place_photos")
	}

	if err = googlePlacePhotoObj.doAfterSelectHooks(ctx, exec); err != nil {
		return googlePlacePhotoObj, err
	}

	return googlePlacePhotoObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *GooglePlacePhoto) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("entities: no google_place_photos provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if queries.MustTime(o.CreatedAt).IsZero() {
			queries.SetScanner(&o.CreatedAt, currTime)
		}
		if queries.MustTime(o.UpdatedAt).IsZero() {
			queries.SetScanner(&o.UpdatedAt, currTime)
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(googlePlacePhotoColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	googlePlacePhotoInsertCacheMut.RLock()
	cache, cached := googlePlacePhotoInsertCache[key]
	googlePlacePhotoInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			googlePlacePhotoAllColumns,
			googlePlacePhotoColumnsWithDefault,
			googlePlacePhotoColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(googlePlacePhotoType, googlePlacePhotoMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(googlePlacePhotoType, googlePlacePhotoMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `google_place_photos` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `google_place_photos` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `google_place_photos` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, googlePlacePhotoPrimaryKeyColumns))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	_, err = exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "entities: unable to insert into google_place_photos")
	}

	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.ID,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "entities: unable to populate default values for google_place_photos")
	}

CacheNoHooks:
	if !cached {
		googlePlacePhotoInsertCacheMut.Lock()
		googlePlacePhotoInsertCache[key] = cache
		googlePlacePhotoInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the GooglePlacePhoto.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *GooglePlacePhoto) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		queries.SetScanner(&o.UpdatedAt, currTime)
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	googlePlacePhotoUpdateCacheMut.RLock()
	cache, cached := googlePlacePhotoUpdateCache[key]
	googlePlacePhotoUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			googlePlacePhotoAllColumns,
			googlePlacePhotoPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("entities: unable to update google_place_photos, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `google_place_photos` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, googlePlacePhotoPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(googlePlacePhotoType, googlePlacePhotoMapping, append(wl, googlePlacePhotoPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "entities: unable to update google_place_photos row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "entities: failed to get rows affected by update for google_place_photos")
	}

	if !cached {
		googlePlacePhotoUpdateCacheMut.Lock()
		googlePlacePhotoUpdateCache[key] = cache
		googlePlacePhotoUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q googlePlacePhotoQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "entities: unable to update all for google_place_photos")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "entities: unable to retrieve rows affected for google_place_photos")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o GooglePlacePhotoSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("entities: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), googlePlacePhotoPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `google_place_photos` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, googlePlacePhotoPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "entities: unable to update all in googlePlacePhoto slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "entities: unable to retrieve rows affected all in update all googlePlacePhoto")
	}
	return rowsAff, nil
}

var mySQLGooglePlacePhotoUniqueColumns = []string{
	"id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *GooglePlacePhoto) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("entities: no google_place_photos provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if queries.MustTime(o.CreatedAt).IsZero() {
			queries.SetScanner(&o.CreatedAt, currTime)
		}
		queries.SetScanner(&o.UpdatedAt, currTime)
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(googlePlacePhotoColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLGooglePlacePhotoUniqueColumns, o)

	if len(nzUniques) == 0 {
		return errors.New("cannot upsert with a table that cannot conflict on a unique column")
	}

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzUniques {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	googlePlacePhotoUpsertCacheMut.RLock()
	cache, cached := googlePlacePhotoUpsertCache[key]
	googlePlacePhotoUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			googlePlacePhotoAllColumns,
			googlePlacePhotoColumnsWithDefault,
			googlePlacePhotoColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			googlePlacePhotoAllColumns,
			googlePlacePhotoPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("entities: unable to upsert google_place_photos, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`google_place_photos`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `google_place_photos` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(googlePlacePhotoType, googlePlacePhotoMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(googlePlacePhotoType, googlePlacePhotoMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	_, err = exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "entities: unable to upsert for google_place_photos")
	}

	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(googlePlacePhotoType, googlePlacePhotoMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "entities: unable to retrieve unique values for google_place_photos")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "entities: unable to populate default values for google_place_photos")
	}

CacheNoHooks:
	if !cached {
		googlePlacePhotoUpsertCacheMut.Lock()
		googlePlacePhotoUpsertCache[key] = cache
		googlePlacePhotoUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single GooglePlacePhoto record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *GooglePlacePhoto) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("entities: no GooglePlacePhoto provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), googlePlacePhotoPrimaryKeyMapping)
	sql := "DELETE FROM `google_place_photos` WHERE `id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "entities: unable to delete from google_place_photos")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "entities: failed to get rows affected by delete for google_place_photos")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q googlePlacePhotoQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("entities: no googlePlacePhotoQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "entities: unable to delete all from google_place_photos")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "entities: failed to get rows affected by deleteall for google_place_photos")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o GooglePlacePhotoSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(googlePlacePhotoBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), googlePlacePhotoPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `google_place_photos` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, googlePlacePhotoPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "entities: unable to delete all from googlePlacePhoto slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "entities: failed to get rows affected by deleteall for google_place_photos")
	}

	if len(googlePlacePhotoAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *GooglePlacePhoto) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindGooglePlacePhoto(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *GooglePlacePhotoSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := GooglePlacePhotoSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), googlePlacePhotoPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `google_place_photos`.* FROM `google_place_photos` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, googlePlacePhotoPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "entities: unable to reload all in GooglePlacePhotoSlice")
	}

	*o = slice

	return nil
}

// GooglePlacePhotoExists checks if the GooglePlacePhoto row exists.
func GooglePlacePhotoExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `google_place_photos` where `id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "entities: unable to check if google_place_photos exists")
	}

	return exists, nil
}

// Exists checks if the GooglePlacePhoto row exists.
func (o *GooglePlacePhoto) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return GooglePlacePhotoExists(ctx, exec, o.ID)
}
