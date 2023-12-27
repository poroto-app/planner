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

// GooglePlacePhotoAttribution is an object representing the database table.
type GooglePlacePhotoAttribution struct {
	ID              string    `boil:"id" json:"id" toml:"id" yaml:"id"`
	PhotoReference  string    `boil:"photo_reference" json:"photo_reference" toml:"photo_reference" yaml:"photo_reference"`
	HTMLAttribution string    `boil:"html_attribution" json:"html_attribution" toml:"html_attribution" yaml:"html_attribution"`
	CreatedAt       null.Time `boil:"created_at" json:"created_at,omitempty" toml:"created_at" yaml:"created_at,omitempty"`
	UpdatedAt       null.Time `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`

	R *googlePlacePhotoAttributionR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L googlePlacePhotoAttributionL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var GooglePlacePhotoAttributionColumns = struct {
	ID              string
	PhotoReference  string
	HTMLAttribution string
	CreatedAt       string
	UpdatedAt       string
}{
	ID:              "id",
	PhotoReference:  "photo_reference",
	HTMLAttribution: "html_attribution",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

var GooglePlacePhotoAttributionTableColumns = struct {
	ID              string
	PhotoReference  string
	HTMLAttribution string
	CreatedAt       string
	UpdatedAt       string
}{
	ID:              "google_place_photo_attributions.id",
	PhotoReference:  "google_place_photo_attributions.photo_reference",
	HTMLAttribution: "google_place_photo_attributions.html_attribution",
	CreatedAt:       "google_place_photo_attributions.created_at",
	UpdatedAt:       "google_place_photo_attributions.updated_at",
}

// Generated where

var GooglePlacePhotoAttributionWhere = struct {
	ID              whereHelperstring
	PhotoReference  whereHelperstring
	HTMLAttribution whereHelperstring
	CreatedAt       whereHelpernull_Time
	UpdatedAt       whereHelpernull_Time
}{
	ID:              whereHelperstring{field: "`google_place_photo_attributions`.`id`"},
	PhotoReference:  whereHelperstring{field: "`google_place_photo_attributions`.`photo_reference`"},
	HTMLAttribution: whereHelperstring{field: "`google_place_photo_attributions`.`html_attribution`"},
	CreatedAt:       whereHelpernull_Time{field: "`google_place_photo_attributions`.`created_at`"},
	UpdatedAt:       whereHelpernull_Time{field: "`google_place_photo_attributions`.`updated_at`"},
}

// GooglePlacePhotoAttributionRels is where relationship names are stored.
var GooglePlacePhotoAttributionRels = struct {
	PhotoReferenceGooglePlacePhotoReference string
}{
	PhotoReferenceGooglePlacePhotoReference: "PhotoReferenceGooglePlacePhotoReference",
}

// googlePlacePhotoAttributionR is where relationships are stored.
type googlePlacePhotoAttributionR struct {
	PhotoReferenceGooglePlacePhotoReference *GooglePlacePhotoReference `boil:"PhotoReferenceGooglePlacePhotoReference" json:"PhotoReferenceGooglePlacePhotoReference" toml:"PhotoReferenceGooglePlacePhotoReference" yaml:"PhotoReferenceGooglePlacePhotoReference"`
}

// NewStruct creates a new relationship struct
func (*googlePlacePhotoAttributionR) NewStruct() *googlePlacePhotoAttributionR {
	return &googlePlacePhotoAttributionR{}
}

func (r *googlePlacePhotoAttributionR) GetPhotoReferenceGooglePlacePhotoReference() *GooglePlacePhotoReference {
	if r == nil {
		return nil
	}
	return r.PhotoReferenceGooglePlacePhotoReference
}

// googlePlacePhotoAttributionL is where Load methods for each relationship are stored.
type googlePlacePhotoAttributionL struct{}

var (
	googlePlacePhotoAttributionAllColumns            = []string{"id", "photo_reference", "html_attribution", "created_at", "updated_at"}
	googlePlacePhotoAttributionColumnsWithoutDefault = []string{"photo_reference", "html_attribution"}
	googlePlacePhotoAttributionColumnsWithDefault    = []string{"id", "created_at", "updated_at"}
	googlePlacePhotoAttributionPrimaryKeyColumns     = []string{"id"}
	googlePlacePhotoAttributionGeneratedColumns      = []string{}
)

type (
	// GooglePlacePhotoAttributionSlice is an alias for a slice of pointers to GooglePlacePhotoAttribution.
	// This should almost always be used instead of []GooglePlacePhotoAttribution.
	GooglePlacePhotoAttributionSlice []*GooglePlacePhotoAttribution
	// GooglePlacePhotoAttributionHook is the signature for custom GooglePlacePhotoAttribution hook methods
	GooglePlacePhotoAttributionHook func(context.Context, boil.ContextExecutor, *GooglePlacePhotoAttribution) error

	googlePlacePhotoAttributionQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	googlePlacePhotoAttributionType                 = reflect.TypeOf(&GooglePlacePhotoAttribution{})
	googlePlacePhotoAttributionMapping              = queries.MakeStructMapping(googlePlacePhotoAttributionType)
	googlePlacePhotoAttributionPrimaryKeyMapping, _ = queries.BindMapping(googlePlacePhotoAttributionType, googlePlacePhotoAttributionMapping, googlePlacePhotoAttributionPrimaryKeyColumns)
	googlePlacePhotoAttributionInsertCacheMut       sync.RWMutex
	googlePlacePhotoAttributionInsertCache          = make(map[string]insertCache)
	googlePlacePhotoAttributionUpdateCacheMut       sync.RWMutex
	googlePlacePhotoAttributionUpdateCache          = make(map[string]updateCache)
	googlePlacePhotoAttributionUpsertCacheMut       sync.RWMutex
	googlePlacePhotoAttributionUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var googlePlacePhotoAttributionAfterSelectHooks []GooglePlacePhotoAttributionHook

var googlePlacePhotoAttributionBeforeInsertHooks []GooglePlacePhotoAttributionHook
var googlePlacePhotoAttributionAfterInsertHooks []GooglePlacePhotoAttributionHook

var googlePlacePhotoAttributionBeforeUpdateHooks []GooglePlacePhotoAttributionHook
var googlePlacePhotoAttributionAfterUpdateHooks []GooglePlacePhotoAttributionHook

var googlePlacePhotoAttributionBeforeDeleteHooks []GooglePlacePhotoAttributionHook
var googlePlacePhotoAttributionAfterDeleteHooks []GooglePlacePhotoAttributionHook

var googlePlacePhotoAttributionBeforeUpsertHooks []GooglePlacePhotoAttributionHook
var googlePlacePhotoAttributionAfterUpsertHooks []GooglePlacePhotoAttributionHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *GooglePlacePhotoAttribution) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlacePhotoAttributionAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *GooglePlacePhotoAttribution) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlacePhotoAttributionBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *GooglePlacePhotoAttribution) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlacePhotoAttributionAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *GooglePlacePhotoAttribution) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlacePhotoAttributionBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *GooglePlacePhotoAttribution) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlacePhotoAttributionAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *GooglePlacePhotoAttribution) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlacePhotoAttributionBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *GooglePlacePhotoAttribution) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlacePhotoAttributionAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *GooglePlacePhotoAttribution) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlacePhotoAttributionBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *GooglePlacePhotoAttribution) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlacePhotoAttributionAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddGooglePlacePhotoAttributionHook registers your hook function for all future operations.
func AddGooglePlacePhotoAttributionHook(hookPoint boil.HookPoint, googlePlacePhotoAttributionHook GooglePlacePhotoAttributionHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		googlePlacePhotoAttributionAfterSelectHooks = append(googlePlacePhotoAttributionAfterSelectHooks, googlePlacePhotoAttributionHook)
	case boil.BeforeInsertHook:
		googlePlacePhotoAttributionBeforeInsertHooks = append(googlePlacePhotoAttributionBeforeInsertHooks, googlePlacePhotoAttributionHook)
	case boil.AfterInsertHook:
		googlePlacePhotoAttributionAfterInsertHooks = append(googlePlacePhotoAttributionAfterInsertHooks, googlePlacePhotoAttributionHook)
	case boil.BeforeUpdateHook:
		googlePlacePhotoAttributionBeforeUpdateHooks = append(googlePlacePhotoAttributionBeforeUpdateHooks, googlePlacePhotoAttributionHook)
	case boil.AfterUpdateHook:
		googlePlacePhotoAttributionAfterUpdateHooks = append(googlePlacePhotoAttributionAfterUpdateHooks, googlePlacePhotoAttributionHook)
	case boil.BeforeDeleteHook:
		googlePlacePhotoAttributionBeforeDeleteHooks = append(googlePlacePhotoAttributionBeforeDeleteHooks, googlePlacePhotoAttributionHook)
	case boil.AfterDeleteHook:
		googlePlacePhotoAttributionAfterDeleteHooks = append(googlePlacePhotoAttributionAfterDeleteHooks, googlePlacePhotoAttributionHook)
	case boil.BeforeUpsertHook:
		googlePlacePhotoAttributionBeforeUpsertHooks = append(googlePlacePhotoAttributionBeforeUpsertHooks, googlePlacePhotoAttributionHook)
	case boil.AfterUpsertHook:
		googlePlacePhotoAttributionAfterUpsertHooks = append(googlePlacePhotoAttributionAfterUpsertHooks, googlePlacePhotoAttributionHook)
	}
}

// One returns a single googlePlacePhotoAttribution record from the query.
func (q googlePlacePhotoAttributionQuery) One(ctx context.Context, exec boil.ContextExecutor) (*GooglePlacePhotoAttribution, error) {
	o := &GooglePlacePhotoAttribution{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "entities: failed to execute a one query for google_place_photo_attributions")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all GooglePlacePhotoAttribution records from the query.
func (q googlePlacePhotoAttributionQuery) All(ctx context.Context, exec boil.ContextExecutor) (GooglePlacePhotoAttributionSlice, error) {
	var o []*GooglePlacePhotoAttribution

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "entities: failed to assign all query results to GooglePlacePhotoAttribution slice")
	}

	if len(googlePlacePhotoAttributionAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all GooglePlacePhotoAttribution records in the query.
func (q googlePlacePhotoAttributionQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "entities: failed to count google_place_photo_attributions rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q googlePlacePhotoAttributionQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "entities: failed to check if google_place_photo_attributions exists")
	}

	return count > 0, nil
}

// PhotoReferenceGooglePlacePhotoReference pointed to by the foreign key.
func (o *GooglePlacePhotoAttribution) PhotoReferenceGooglePlacePhotoReference(mods ...qm.QueryMod) googlePlacePhotoReferenceQuery {
	queryMods := []qm.QueryMod{
		qm.Where("`photo_reference` = ?", o.PhotoReference),
	}

	queryMods = append(queryMods, mods...)

	return GooglePlacePhotoReferences(queryMods...)
}

// LoadPhotoReferenceGooglePlacePhotoReference allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (googlePlacePhotoAttributionL) LoadPhotoReferenceGooglePlacePhotoReference(ctx context.Context, e boil.ContextExecutor, singular bool, maybeGooglePlacePhotoAttribution interface{}, mods queries.Applicator) error {
	var slice []*GooglePlacePhotoAttribution
	var object *GooglePlacePhotoAttribution

	if singular {
		var ok bool
		object, ok = maybeGooglePlacePhotoAttribution.(*GooglePlacePhotoAttribution)
		if !ok {
			object = new(GooglePlacePhotoAttribution)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeGooglePlacePhotoAttribution)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeGooglePlacePhotoAttribution))
			}
		}
	} else {
		s, ok := maybeGooglePlacePhotoAttribution.(*[]*GooglePlacePhotoAttribution)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeGooglePlacePhotoAttribution)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeGooglePlacePhotoAttribution))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &googlePlacePhotoAttributionR{}
		}
		args = append(args, object.PhotoReference)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &googlePlacePhotoAttributionR{}
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
		foreign.R.PhotoReferenceGooglePlacePhotoAttributions = append(foreign.R.PhotoReferenceGooglePlacePhotoAttributions, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.PhotoReference == foreign.PhotoReference {
				local.R.PhotoReferenceGooglePlacePhotoReference = foreign
				if foreign.R == nil {
					foreign.R = &googlePlacePhotoReferenceR{}
				}
				foreign.R.PhotoReferenceGooglePlacePhotoAttributions = append(foreign.R.PhotoReferenceGooglePlacePhotoAttributions, local)
				break
			}
		}
	}

	return nil
}

// SetPhotoReferenceGooglePlacePhotoReference of the googlePlacePhotoAttribution to the related item.
// Sets o.R.PhotoReferenceGooglePlacePhotoReference to related.
// Adds o to related.R.PhotoReferenceGooglePlacePhotoAttributions.
func (o *GooglePlacePhotoAttribution) SetPhotoReferenceGooglePlacePhotoReference(ctx context.Context, exec boil.ContextExecutor, insert bool, related *GooglePlacePhotoReference) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE `google_place_photo_attributions` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, []string{"photo_reference"}),
		strmangle.WhereClause("`", "`", 0, googlePlacePhotoAttributionPrimaryKeyColumns),
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
		o.R = &googlePlacePhotoAttributionR{
			PhotoReferenceGooglePlacePhotoReference: related,
		}
	} else {
		o.R.PhotoReferenceGooglePlacePhotoReference = related
	}

	if related.R == nil {
		related.R = &googlePlacePhotoReferenceR{
			PhotoReferenceGooglePlacePhotoAttributions: GooglePlacePhotoAttributionSlice{o},
		}
	} else {
		related.R.PhotoReferenceGooglePlacePhotoAttributions = append(related.R.PhotoReferenceGooglePlacePhotoAttributions, o)
	}

	return nil
}

// GooglePlacePhotoAttributions retrieves all the records using an executor.
func GooglePlacePhotoAttributions(mods ...qm.QueryMod) googlePlacePhotoAttributionQuery {
	mods = append(mods, qm.From("`google_place_photo_attributions`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`google_place_photo_attributions`.*"})
	}

	return googlePlacePhotoAttributionQuery{q}
}

// FindGooglePlacePhotoAttribution retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindGooglePlacePhotoAttribution(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*GooglePlacePhotoAttribution, error) {
	googlePlacePhotoAttributionObj := &GooglePlacePhotoAttribution{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `google_place_photo_attributions` where `id`=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, googlePlacePhotoAttributionObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "entities: unable to select from google_place_photo_attributions")
	}

	if err = googlePlacePhotoAttributionObj.doAfterSelectHooks(ctx, exec); err != nil {
		return googlePlacePhotoAttributionObj, err
	}

	return googlePlacePhotoAttributionObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *GooglePlacePhotoAttribution) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("entities: no google_place_photo_attributions provided for insertion")
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

	nzDefaults := queries.NonZeroDefaultSet(googlePlacePhotoAttributionColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	googlePlacePhotoAttributionInsertCacheMut.RLock()
	cache, cached := googlePlacePhotoAttributionInsertCache[key]
	googlePlacePhotoAttributionInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			googlePlacePhotoAttributionAllColumns,
			googlePlacePhotoAttributionColumnsWithDefault,
			googlePlacePhotoAttributionColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(googlePlacePhotoAttributionType, googlePlacePhotoAttributionMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(googlePlacePhotoAttributionType, googlePlacePhotoAttributionMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `google_place_photo_attributions` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `google_place_photo_attributions` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `google_place_photo_attributions` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, googlePlacePhotoAttributionPrimaryKeyColumns))
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
		return errors.Wrap(err, "entities: unable to insert into google_place_photo_attributions")
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
		return errors.Wrap(err, "entities: unable to populate default values for google_place_photo_attributions")
	}

CacheNoHooks:
	if !cached {
		googlePlacePhotoAttributionInsertCacheMut.Lock()
		googlePlacePhotoAttributionInsertCache[key] = cache
		googlePlacePhotoAttributionInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the GooglePlacePhotoAttribution.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *GooglePlacePhotoAttribution) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		queries.SetScanner(&o.UpdatedAt, currTime)
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	googlePlacePhotoAttributionUpdateCacheMut.RLock()
	cache, cached := googlePlacePhotoAttributionUpdateCache[key]
	googlePlacePhotoAttributionUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			googlePlacePhotoAttributionAllColumns,
			googlePlacePhotoAttributionPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("entities: unable to update google_place_photo_attributions, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `google_place_photo_attributions` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, googlePlacePhotoAttributionPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(googlePlacePhotoAttributionType, googlePlacePhotoAttributionMapping, append(wl, googlePlacePhotoAttributionPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "entities: unable to update google_place_photo_attributions row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "entities: failed to get rows affected by update for google_place_photo_attributions")
	}

	if !cached {
		googlePlacePhotoAttributionUpdateCacheMut.Lock()
		googlePlacePhotoAttributionUpdateCache[key] = cache
		googlePlacePhotoAttributionUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q googlePlacePhotoAttributionQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "entities: unable to update all for google_place_photo_attributions")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "entities: unable to retrieve rows affected for google_place_photo_attributions")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o GooglePlacePhotoAttributionSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), googlePlacePhotoAttributionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `google_place_photo_attributions` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, googlePlacePhotoAttributionPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "entities: unable to update all in googlePlacePhotoAttribution slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "entities: unable to retrieve rows affected all in update all googlePlacePhotoAttribution")
	}
	return rowsAff, nil
}

var mySQLGooglePlacePhotoAttributionUniqueColumns = []string{
	"id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *GooglePlacePhotoAttribution) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("entities: no google_place_photo_attributions provided for upsert")
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

	nzDefaults := queries.NonZeroDefaultSet(googlePlacePhotoAttributionColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLGooglePlacePhotoAttributionUniqueColumns, o)

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

	googlePlacePhotoAttributionUpsertCacheMut.RLock()
	cache, cached := googlePlacePhotoAttributionUpsertCache[key]
	googlePlacePhotoAttributionUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			googlePlacePhotoAttributionAllColumns,
			googlePlacePhotoAttributionColumnsWithDefault,
			googlePlacePhotoAttributionColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			googlePlacePhotoAttributionAllColumns,
			googlePlacePhotoAttributionPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("entities: unable to upsert google_place_photo_attributions, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`google_place_photo_attributions`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `google_place_photo_attributions` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(googlePlacePhotoAttributionType, googlePlacePhotoAttributionMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(googlePlacePhotoAttributionType, googlePlacePhotoAttributionMapping, ret)
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
		return errors.Wrap(err, "entities: unable to upsert for google_place_photo_attributions")
	}

	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(googlePlacePhotoAttributionType, googlePlacePhotoAttributionMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "entities: unable to retrieve unique values for google_place_photo_attributions")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "entities: unable to populate default values for google_place_photo_attributions")
	}

CacheNoHooks:
	if !cached {
		googlePlacePhotoAttributionUpsertCacheMut.Lock()
		googlePlacePhotoAttributionUpsertCache[key] = cache
		googlePlacePhotoAttributionUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single GooglePlacePhotoAttribution record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *GooglePlacePhotoAttribution) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("entities: no GooglePlacePhotoAttribution provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), googlePlacePhotoAttributionPrimaryKeyMapping)
	sql := "DELETE FROM `google_place_photo_attributions` WHERE `id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "entities: unable to delete from google_place_photo_attributions")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "entities: failed to get rows affected by delete for google_place_photo_attributions")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q googlePlacePhotoAttributionQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("entities: no googlePlacePhotoAttributionQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "entities: unable to delete all from google_place_photo_attributions")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "entities: failed to get rows affected by deleteall for google_place_photo_attributions")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o GooglePlacePhotoAttributionSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(googlePlacePhotoAttributionBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), googlePlacePhotoAttributionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `google_place_photo_attributions` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, googlePlacePhotoAttributionPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "entities: unable to delete all from googlePlacePhotoAttribution slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "entities: failed to get rows affected by deleteall for google_place_photo_attributions")
	}

	if len(googlePlacePhotoAttributionAfterDeleteHooks) != 0 {
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
func (o *GooglePlacePhotoAttribution) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindGooglePlacePhotoAttribution(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *GooglePlacePhotoAttributionSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := GooglePlacePhotoAttributionSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), googlePlacePhotoAttributionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `google_place_photo_attributions`.* FROM `google_place_photo_attributions` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, googlePlacePhotoAttributionPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "entities: unable to reload all in GooglePlacePhotoAttributionSlice")
	}

	*o = slice

	return nil
}

// GooglePlacePhotoAttributionExists checks if the GooglePlacePhotoAttribution row exists.
func GooglePlacePhotoAttributionExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `google_place_photo_attributions` where `id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "entities: unable to check if google_place_photo_attributions exists")
	}

	return exists, nil
}

// Exists checks if the GooglePlacePhotoAttribution row exists.
func (o *GooglePlacePhotoAttribution) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return GooglePlacePhotoAttributionExists(ctx, exec, o.ID)
}
