// Code generated by SQLBoiler 4.15.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package generated

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

// GooglePlaceType is an object representing the database table.
type GooglePlaceType struct {
	ID            string    `boil:"id" json:"id" toml:"id" yaml:"id"`
	GooglePlaceID string    `boil:"google_place_id" json:"google_place_id" toml:"google_place_id" yaml:"google_place_id"`
	Type          string    `boil:"type" json:"type" toml:"type" yaml:"type"`
	OrderNum      int       `boil:"order_num" json:"order_num" toml:"order_num" yaml:"order_num"`
	CreatedAt     null.Time `boil:"created_at" json:"created_at,omitempty" toml:"created_at" yaml:"created_at,omitempty"`
	UpdatedAt     null.Time `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`

	R *googlePlaceTypeR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L googlePlaceTypeL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var GooglePlaceTypeColumns = struct {
	ID            string
	GooglePlaceID string
	Type          string
	OrderNum      string
	CreatedAt     string
	UpdatedAt     string
}{
	ID:            "id",
	GooglePlaceID: "google_place_id",
	Type:          "type",
	OrderNum:      "order_num",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

var GooglePlaceTypeTableColumns = struct {
	ID            string
	GooglePlaceID string
	Type          string
	OrderNum      string
	CreatedAt     string
	UpdatedAt     string
}{
	ID:            "google_place_types.id",
	GooglePlaceID: "google_place_types.google_place_id",
	Type:          "google_place_types.type",
	OrderNum:      "google_place_types.order_num",
	CreatedAt:     "google_place_types.created_at",
	UpdatedAt:     "google_place_types.updated_at",
}

// Generated where

var GooglePlaceTypeWhere = struct {
	ID            whereHelperstring
	GooglePlaceID whereHelperstring
	Type          whereHelperstring
	OrderNum      whereHelperint
	CreatedAt     whereHelpernull_Time
	UpdatedAt     whereHelpernull_Time
}{
	ID:            whereHelperstring{field: "`google_place_types`.`id`"},
	GooglePlaceID: whereHelperstring{field: "`google_place_types`.`google_place_id`"},
	Type:          whereHelperstring{field: "`google_place_types`.`type`"},
	OrderNum:      whereHelperint{field: "`google_place_types`.`order_num`"},
	CreatedAt:     whereHelpernull_Time{field: "`google_place_types`.`created_at`"},
	UpdatedAt:     whereHelpernull_Time{field: "`google_place_types`.`updated_at`"},
}

// GooglePlaceTypeRels is where relationship names are stored.
var GooglePlaceTypeRels = struct {
	GooglePlace string
}{
	GooglePlace: "GooglePlace",
}

// googlePlaceTypeR is where relationships are stored.
type googlePlaceTypeR struct {
	GooglePlace *GooglePlace `boil:"GooglePlace" json:"GooglePlace" toml:"GooglePlace" yaml:"GooglePlace"`
}

// NewStruct creates a new relationship struct
func (*googlePlaceTypeR) NewStruct() *googlePlaceTypeR {
	return &googlePlaceTypeR{}
}

func (r *googlePlaceTypeR) GetGooglePlace() *GooglePlace {
	if r == nil {
		return nil
	}
	return r.GooglePlace
}

// googlePlaceTypeL is where Load methods for each relationship are stored.
type googlePlaceTypeL struct{}

var (
	googlePlaceTypeAllColumns            = []string{"id", "google_place_id", "type", "order_num", "created_at", "updated_at"}
	googlePlaceTypeColumnsWithoutDefault = []string{"google_place_id", "type", "order_num"}
	googlePlaceTypeColumnsWithDefault    = []string{"id", "created_at", "updated_at"}
	googlePlaceTypePrimaryKeyColumns     = []string{"id"}
	googlePlaceTypeGeneratedColumns      = []string{}
)

type (
	// GooglePlaceTypeSlice is an alias for a slice of pointers to GooglePlaceType.
	// This should almost always be used instead of []GooglePlaceType.
	GooglePlaceTypeSlice []*GooglePlaceType
	// GooglePlaceTypeHook is the signature for custom GooglePlaceType hook methods
	GooglePlaceTypeHook func(context.Context, boil.ContextExecutor, *GooglePlaceType) error

	googlePlaceTypeQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	googlePlaceTypeType                 = reflect.TypeOf(&GooglePlaceType{})
	googlePlaceTypeMapping              = queries.MakeStructMapping(googlePlaceTypeType)
	googlePlaceTypePrimaryKeyMapping, _ = queries.BindMapping(googlePlaceTypeType, googlePlaceTypeMapping, googlePlaceTypePrimaryKeyColumns)
	googlePlaceTypeInsertCacheMut       sync.RWMutex
	googlePlaceTypeInsertCache          = make(map[string]insertCache)
	googlePlaceTypeUpdateCacheMut       sync.RWMutex
	googlePlaceTypeUpdateCache          = make(map[string]updateCache)
	googlePlaceTypeUpsertCacheMut       sync.RWMutex
	googlePlaceTypeUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var googlePlaceTypeAfterSelectHooks []GooglePlaceTypeHook

var googlePlaceTypeBeforeInsertHooks []GooglePlaceTypeHook
var googlePlaceTypeAfterInsertHooks []GooglePlaceTypeHook

var googlePlaceTypeBeforeUpdateHooks []GooglePlaceTypeHook
var googlePlaceTypeAfterUpdateHooks []GooglePlaceTypeHook

var googlePlaceTypeBeforeDeleteHooks []GooglePlaceTypeHook
var googlePlaceTypeAfterDeleteHooks []GooglePlaceTypeHook

var googlePlaceTypeBeforeUpsertHooks []GooglePlaceTypeHook
var googlePlaceTypeAfterUpsertHooks []GooglePlaceTypeHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *GooglePlaceType) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlaceTypeAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *GooglePlaceType) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlaceTypeBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *GooglePlaceType) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlaceTypeAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *GooglePlaceType) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlaceTypeBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *GooglePlaceType) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlaceTypeAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *GooglePlaceType) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlaceTypeBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *GooglePlaceType) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlaceTypeAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *GooglePlaceType) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlaceTypeBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *GooglePlaceType) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range googlePlaceTypeAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddGooglePlaceTypeHook registers your hook function for all future operations.
func AddGooglePlaceTypeHook(hookPoint boil.HookPoint, googlePlaceTypeHook GooglePlaceTypeHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		googlePlaceTypeAfterSelectHooks = append(googlePlaceTypeAfterSelectHooks, googlePlaceTypeHook)
	case boil.BeforeInsertHook:
		googlePlaceTypeBeforeInsertHooks = append(googlePlaceTypeBeforeInsertHooks, googlePlaceTypeHook)
	case boil.AfterInsertHook:
		googlePlaceTypeAfterInsertHooks = append(googlePlaceTypeAfterInsertHooks, googlePlaceTypeHook)
	case boil.BeforeUpdateHook:
		googlePlaceTypeBeforeUpdateHooks = append(googlePlaceTypeBeforeUpdateHooks, googlePlaceTypeHook)
	case boil.AfterUpdateHook:
		googlePlaceTypeAfterUpdateHooks = append(googlePlaceTypeAfterUpdateHooks, googlePlaceTypeHook)
	case boil.BeforeDeleteHook:
		googlePlaceTypeBeforeDeleteHooks = append(googlePlaceTypeBeforeDeleteHooks, googlePlaceTypeHook)
	case boil.AfterDeleteHook:
		googlePlaceTypeAfterDeleteHooks = append(googlePlaceTypeAfterDeleteHooks, googlePlaceTypeHook)
	case boil.BeforeUpsertHook:
		googlePlaceTypeBeforeUpsertHooks = append(googlePlaceTypeBeforeUpsertHooks, googlePlaceTypeHook)
	case boil.AfterUpsertHook:
		googlePlaceTypeAfterUpsertHooks = append(googlePlaceTypeAfterUpsertHooks, googlePlaceTypeHook)
	}
}

// One returns a single googlePlaceType record from the query.
func (q googlePlaceTypeQuery) One(ctx context.Context, exec boil.ContextExecutor) (*GooglePlaceType, error) {
	o := &GooglePlaceType{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "generated: failed to execute a one query for google_place_types")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all GooglePlaceType records from the query.
func (q googlePlaceTypeQuery) All(ctx context.Context, exec boil.ContextExecutor) (GooglePlaceTypeSlice, error) {
	var o []*GooglePlaceType

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "generated: failed to assign all query results to GooglePlaceType slice")
	}

	if len(googlePlaceTypeAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all GooglePlaceType records in the query.
func (q googlePlaceTypeQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "generated: failed to count google_place_types rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q googlePlaceTypeQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "generated: failed to check if google_place_types exists")
	}

	return count > 0, nil
}

// GooglePlace pointed to by the foreign key.
func (o *GooglePlaceType) GooglePlace(mods ...qm.QueryMod) googlePlaceQuery {
	queryMods := []qm.QueryMod{
		qm.Where("`google_place_id` = ?", o.GooglePlaceID),
	}

	queryMods = append(queryMods, mods...)

	return GooglePlaces(queryMods...)
}

// LoadGooglePlace allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (googlePlaceTypeL) LoadGooglePlace(ctx context.Context, e boil.ContextExecutor, singular bool, maybeGooglePlaceType interface{}, mods queries.Applicator) error {
	var slice []*GooglePlaceType
	var object *GooglePlaceType

	if singular {
		var ok bool
		object, ok = maybeGooglePlaceType.(*GooglePlaceType)
		if !ok {
			object = new(GooglePlaceType)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeGooglePlaceType)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeGooglePlaceType))
			}
		}
	} else {
		s, ok := maybeGooglePlaceType.(*[]*GooglePlaceType)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeGooglePlaceType)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeGooglePlaceType))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &googlePlaceTypeR{}
		}
		args = append(args, object.GooglePlaceID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &googlePlaceTypeR{}
			}

			for _, a := range args {
				if a == obj.GooglePlaceID {
					continue Outer
				}
			}

			args = append(args, obj.GooglePlaceID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`google_places`),
		qm.WhereIn(`google_places.google_place_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load GooglePlace")
	}

	var resultSlice []*GooglePlace
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice GooglePlace")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for google_places")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for google_places")
	}

	if len(googlePlaceAfterSelectHooks) != 0 {
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
		object.R.GooglePlace = foreign
		if foreign.R == nil {
			foreign.R = &googlePlaceR{}
		}
		foreign.R.GooglePlaceTypes = append(foreign.R.GooglePlaceTypes, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.GooglePlaceID == foreign.GooglePlaceID {
				local.R.GooglePlace = foreign
				if foreign.R == nil {
					foreign.R = &googlePlaceR{}
				}
				foreign.R.GooglePlaceTypes = append(foreign.R.GooglePlaceTypes, local)
				break
			}
		}
	}

	return nil
}

// SetGooglePlace of the googlePlaceType to the related item.
// Sets o.R.GooglePlace to related.
// Adds o to related.R.GooglePlaceTypes.
func (o *GooglePlaceType) SetGooglePlace(ctx context.Context, exec boil.ContextExecutor, insert bool, related *GooglePlace) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE `google_place_types` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, []string{"google_place_id"}),
		strmangle.WhereClause("`", "`", 0, googlePlaceTypePrimaryKeyColumns),
	)
	values := []interface{}{related.GooglePlaceID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.GooglePlaceID = related.GooglePlaceID
	if o.R == nil {
		o.R = &googlePlaceTypeR{
			GooglePlace: related,
		}
	} else {
		o.R.GooglePlace = related
	}

	if related.R == nil {
		related.R = &googlePlaceR{
			GooglePlaceTypes: GooglePlaceTypeSlice{o},
		}
	} else {
		related.R.GooglePlaceTypes = append(related.R.GooglePlaceTypes, o)
	}

	return nil
}

// GooglePlaceTypes retrieves all the records using an executor.
func GooglePlaceTypes(mods ...qm.QueryMod) googlePlaceTypeQuery {
	mods = append(mods, qm.From("`google_place_types`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`google_place_types`.*"})
	}

	return googlePlaceTypeQuery{q}
}

// FindGooglePlaceType retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindGooglePlaceType(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*GooglePlaceType, error) {
	googlePlaceTypeObj := &GooglePlaceType{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `google_place_types` where `id`=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, googlePlaceTypeObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "generated: unable to select from google_place_types")
	}

	if err = googlePlaceTypeObj.doAfterSelectHooks(ctx, exec); err != nil {
		return googlePlaceTypeObj, err
	}

	return googlePlaceTypeObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *GooglePlaceType) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("generated: no google_place_types provided for insertion")
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

	nzDefaults := queries.NonZeroDefaultSet(googlePlaceTypeColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	googlePlaceTypeInsertCacheMut.RLock()
	cache, cached := googlePlaceTypeInsertCache[key]
	googlePlaceTypeInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			googlePlaceTypeAllColumns,
			googlePlaceTypeColumnsWithDefault,
			googlePlaceTypeColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(googlePlaceTypeType, googlePlaceTypeMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(googlePlaceTypeType, googlePlaceTypeMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `google_place_types` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `google_place_types` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `google_place_types` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, googlePlaceTypePrimaryKeyColumns))
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
		return errors.Wrap(err, "generated: unable to insert into google_place_types")
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
		return errors.Wrap(err, "generated: unable to populate default values for google_place_types")
	}

CacheNoHooks:
	if !cached {
		googlePlaceTypeInsertCacheMut.Lock()
		googlePlaceTypeInsertCache[key] = cache
		googlePlaceTypeInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the GooglePlaceType.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *GooglePlaceType) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		queries.SetScanner(&o.UpdatedAt, currTime)
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	googlePlaceTypeUpdateCacheMut.RLock()
	cache, cached := googlePlaceTypeUpdateCache[key]
	googlePlaceTypeUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			googlePlaceTypeAllColumns,
			googlePlaceTypePrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("generated: unable to update google_place_types, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `google_place_types` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, googlePlaceTypePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(googlePlaceTypeType, googlePlaceTypeMapping, append(wl, googlePlaceTypePrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "generated: unable to update google_place_types row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "generated: failed to get rows affected by update for google_place_types")
	}

	if !cached {
		googlePlaceTypeUpdateCacheMut.Lock()
		googlePlaceTypeUpdateCache[key] = cache
		googlePlaceTypeUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q googlePlaceTypeQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "generated: unable to update all for google_place_types")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "generated: unable to retrieve rows affected for google_place_types")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o GooglePlaceTypeSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("generated: update all requires at least one column argument")
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), googlePlaceTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `google_place_types` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, googlePlaceTypePrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "generated: unable to update all in googlePlaceType slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "generated: unable to retrieve rows affected all in update all googlePlaceType")
	}
	return rowsAff, nil
}

var mySQLGooglePlaceTypeUniqueColumns = []string{
	"id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *GooglePlaceType) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("generated: no google_place_types provided for upsert")
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

	nzDefaults := queries.NonZeroDefaultSet(googlePlaceTypeColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLGooglePlaceTypeUniqueColumns, o)

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

	googlePlaceTypeUpsertCacheMut.RLock()
	cache, cached := googlePlaceTypeUpsertCache[key]
	googlePlaceTypeUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			googlePlaceTypeAllColumns,
			googlePlaceTypeColumnsWithDefault,
			googlePlaceTypeColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			googlePlaceTypeAllColumns,
			googlePlaceTypePrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("generated: unable to upsert google_place_types, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`google_place_types`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `google_place_types` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(googlePlaceTypeType, googlePlaceTypeMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(googlePlaceTypeType, googlePlaceTypeMapping, ret)
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
		return errors.Wrap(err, "generated: unable to upsert for google_place_types")
	}

	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(googlePlaceTypeType, googlePlaceTypeMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "generated: unable to retrieve unique values for google_place_types")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "generated: unable to populate default values for google_place_types")
	}

CacheNoHooks:
	if !cached {
		googlePlaceTypeUpsertCacheMut.Lock()
		googlePlaceTypeUpsertCache[key] = cache
		googlePlaceTypeUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single GooglePlaceType record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *GooglePlaceType) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("generated: no GooglePlaceType provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), googlePlaceTypePrimaryKeyMapping)
	sql := "DELETE FROM `google_place_types` WHERE `id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "generated: unable to delete from google_place_types")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "generated: failed to get rows affected by delete for google_place_types")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q googlePlaceTypeQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("generated: no googlePlaceTypeQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "generated: unable to delete all from google_place_types")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "generated: failed to get rows affected by deleteall for google_place_types")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o GooglePlaceTypeSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(googlePlaceTypeBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), googlePlaceTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `google_place_types` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, googlePlaceTypePrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "generated: unable to delete all from googlePlaceType slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "generated: failed to get rows affected by deleteall for google_place_types")
	}

	if len(googlePlaceTypeAfterDeleteHooks) != 0 {
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
func (o *GooglePlaceType) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindGooglePlaceType(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *GooglePlaceTypeSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := GooglePlaceTypeSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), googlePlaceTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `google_place_types`.* FROM `google_place_types` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, googlePlaceTypePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "generated: unable to reload all in GooglePlaceTypeSlice")
	}

	*o = slice

	return nil
}

// GooglePlaceTypeExists checks if the GooglePlaceType row exists.
func GooglePlaceTypeExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `google_place_types` where `id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "generated: unable to check if google_place_types exists")
	}

	return exists, nil
}

// Exists checks if the GooglePlaceType row exists.
func (o *GooglePlaceType) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return GooglePlaceTypeExists(ctx, exec, o.ID)
}