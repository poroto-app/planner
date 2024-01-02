// Code generated by SQLBoiler 4.15.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package entities

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/randomize"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testPlanCandidateSetCategories(t *testing.T) {
	t.Parallel()

	query := PlanCandidateSetCategories()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testPlanCandidateSetCategoriesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &PlanCandidateSetCategory{}
	if err = randomize.Struct(seed, o, planCandidateSetCategoryDBTypes, true, planCandidateSetCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PlanCandidateSetCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := PlanCandidateSetCategories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testPlanCandidateSetCategoriesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &PlanCandidateSetCategory{}
	if err = randomize.Struct(seed, o, planCandidateSetCategoryDBTypes, true, planCandidateSetCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PlanCandidateSetCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := PlanCandidateSetCategories().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := PlanCandidateSetCategories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testPlanCandidateSetCategoriesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &PlanCandidateSetCategory{}
	if err = randomize.Struct(seed, o, planCandidateSetCategoryDBTypes, true, planCandidateSetCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PlanCandidateSetCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := PlanCandidateSetCategorySlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := PlanCandidateSetCategories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testPlanCandidateSetCategoriesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &PlanCandidateSetCategory{}
	if err = randomize.Struct(seed, o, planCandidateSetCategoryDBTypes, true, planCandidateSetCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PlanCandidateSetCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := PlanCandidateSetCategoryExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if PlanCandidateSetCategory exists: %s", err)
	}
	if !e {
		t.Errorf("Expected PlanCandidateSetCategoryExists to return true, but got false.")
	}
}

func testPlanCandidateSetCategoriesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &PlanCandidateSetCategory{}
	if err = randomize.Struct(seed, o, planCandidateSetCategoryDBTypes, true, planCandidateSetCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PlanCandidateSetCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	planCandidateSetCategoryFound, err := FindPlanCandidateSetCategory(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if planCandidateSetCategoryFound == nil {
		t.Error("want a record, got nil")
	}
}

func testPlanCandidateSetCategoriesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &PlanCandidateSetCategory{}
	if err = randomize.Struct(seed, o, planCandidateSetCategoryDBTypes, true, planCandidateSetCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PlanCandidateSetCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = PlanCandidateSetCategories().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testPlanCandidateSetCategoriesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &PlanCandidateSetCategory{}
	if err = randomize.Struct(seed, o, planCandidateSetCategoryDBTypes, true, planCandidateSetCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PlanCandidateSetCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := PlanCandidateSetCategories().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testPlanCandidateSetCategoriesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	planCandidateSetCategoryOne := &PlanCandidateSetCategory{}
	planCandidateSetCategoryTwo := &PlanCandidateSetCategory{}
	if err = randomize.Struct(seed, planCandidateSetCategoryOne, planCandidateSetCategoryDBTypes, false, planCandidateSetCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PlanCandidateSetCategory struct: %s", err)
	}
	if err = randomize.Struct(seed, planCandidateSetCategoryTwo, planCandidateSetCategoryDBTypes, false, planCandidateSetCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PlanCandidateSetCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = planCandidateSetCategoryOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = planCandidateSetCategoryTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := PlanCandidateSetCategories().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testPlanCandidateSetCategoriesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	planCandidateSetCategoryOne := &PlanCandidateSetCategory{}
	planCandidateSetCategoryTwo := &PlanCandidateSetCategory{}
	if err = randomize.Struct(seed, planCandidateSetCategoryOne, planCandidateSetCategoryDBTypes, false, planCandidateSetCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PlanCandidateSetCategory struct: %s", err)
	}
	if err = randomize.Struct(seed, planCandidateSetCategoryTwo, planCandidateSetCategoryDBTypes, false, planCandidateSetCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PlanCandidateSetCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = planCandidateSetCategoryOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = planCandidateSetCategoryTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := PlanCandidateSetCategories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func planCandidateSetCategoryBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *PlanCandidateSetCategory) error {
	*o = PlanCandidateSetCategory{}
	return nil
}

func planCandidateSetCategoryAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *PlanCandidateSetCategory) error {
	*o = PlanCandidateSetCategory{}
	return nil
}

func planCandidateSetCategoryAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *PlanCandidateSetCategory) error {
	*o = PlanCandidateSetCategory{}
	return nil
}

func planCandidateSetCategoryBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *PlanCandidateSetCategory) error {
	*o = PlanCandidateSetCategory{}
	return nil
}

func planCandidateSetCategoryAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *PlanCandidateSetCategory) error {
	*o = PlanCandidateSetCategory{}
	return nil
}

func planCandidateSetCategoryBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *PlanCandidateSetCategory) error {
	*o = PlanCandidateSetCategory{}
	return nil
}

func planCandidateSetCategoryAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *PlanCandidateSetCategory) error {
	*o = PlanCandidateSetCategory{}
	return nil
}

func planCandidateSetCategoryBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *PlanCandidateSetCategory) error {
	*o = PlanCandidateSetCategory{}
	return nil
}

func planCandidateSetCategoryAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *PlanCandidateSetCategory) error {
	*o = PlanCandidateSetCategory{}
	return nil
}

func testPlanCandidateSetCategoriesHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &PlanCandidateSetCategory{}
	o := &PlanCandidateSetCategory{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, planCandidateSetCategoryDBTypes, false); err != nil {
		t.Errorf("Unable to randomize PlanCandidateSetCategory object: %s", err)
	}

	AddPlanCandidateSetCategoryHook(boil.BeforeInsertHook, planCandidateSetCategoryBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	planCandidateSetCategoryBeforeInsertHooks = []PlanCandidateSetCategoryHook{}

	AddPlanCandidateSetCategoryHook(boil.AfterInsertHook, planCandidateSetCategoryAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	planCandidateSetCategoryAfterInsertHooks = []PlanCandidateSetCategoryHook{}

	AddPlanCandidateSetCategoryHook(boil.AfterSelectHook, planCandidateSetCategoryAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	planCandidateSetCategoryAfterSelectHooks = []PlanCandidateSetCategoryHook{}

	AddPlanCandidateSetCategoryHook(boil.BeforeUpdateHook, planCandidateSetCategoryBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	planCandidateSetCategoryBeforeUpdateHooks = []PlanCandidateSetCategoryHook{}

	AddPlanCandidateSetCategoryHook(boil.AfterUpdateHook, planCandidateSetCategoryAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	planCandidateSetCategoryAfterUpdateHooks = []PlanCandidateSetCategoryHook{}

	AddPlanCandidateSetCategoryHook(boil.BeforeDeleteHook, planCandidateSetCategoryBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	planCandidateSetCategoryBeforeDeleteHooks = []PlanCandidateSetCategoryHook{}

	AddPlanCandidateSetCategoryHook(boil.AfterDeleteHook, planCandidateSetCategoryAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	planCandidateSetCategoryAfterDeleteHooks = []PlanCandidateSetCategoryHook{}

	AddPlanCandidateSetCategoryHook(boil.BeforeUpsertHook, planCandidateSetCategoryBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	planCandidateSetCategoryBeforeUpsertHooks = []PlanCandidateSetCategoryHook{}

	AddPlanCandidateSetCategoryHook(boil.AfterUpsertHook, planCandidateSetCategoryAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	planCandidateSetCategoryAfterUpsertHooks = []PlanCandidateSetCategoryHook{}
}

func testPlanCandidateSetCategoriesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &PlanCandidateSetCategory{}
	if err = randomize.Struct(seed, o, planCandidateSetCategoryDBTypes, true, planCandidateSetCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PlanCandidateSetCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := PlanCandidateSetCategories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testPlanCandidateSetCategoriesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &PlanCandidateSetCategory{}
	if err = randomize.Struct(seed, o, planCandidateSetCategoryDBTypes, true); err != nil {
		t.Errorf("Unable to randomize PlanCandidateSetCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(planCandidateSetCategoryColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := PlanCandidateSetCategories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testPlanCandidateSetCategoryToOnePlanCandidateSetUsingPlanCandidateSet(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local PlanCandidateSetCategory
	var foreign PlanCandidateSet

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, planCandidateSetCategoryDBTypes, true, planCandidateSetCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PlanCandidateSetCategory struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, planCandidateSetDBTypes, false, planCandidateSetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PlanCandidateSet struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	queries.Assign(&local.PlanCandidateSetID, foreign.ID)
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.PlanCandidateSet().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if !queries.Equal(check.ID, foreign.ID) {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	ranAfterSelectHook := false
	AddPlanCandidateSetHook(boil.AfterSelectHook, func(ctx context.Context, e boil.ContextExecutor, o *PlanCandidateSet) error {
		ranAfterSelectHook = true
		return nil
	})

	slice := PlanCandidateSetCategorySlice{&local}
	if err = local.L.LoadPlanCandidateSet(ctx, tx, false, (*[]*PlanCandidateSetCategory)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.PlanCandidateSet == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.PlanCandidateSet = nil
	if err = local.L.LoadPlanCandidateSet(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.PlanCandidateSet == nil {
		t.Error("struct should have been eager loaded")
	}

	if !ranAfterSelectHook {
		t.Error("failed to run AfterSelect hook for relationship")
	}
}

func testPlanCandidateSetCategoryToOneSetOpPlanCandidateSetUsingPlanCandidateSet(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a PlanCandidateSetCategory
	var b, c PlanCandidateSet

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, planCandidateSetCategoryDBTypes, false, strmangle.SetComplement(planCandidateSetCategoryPrimaryKeyColumns, planCandidateSetCategoryColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, planCandidateSetDBTypes, false, strmangle.SetComplement(planCandidateSetPrimaryKeyColumns, planCandidateSetColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, planCandidateSetDBTypes, false, strmangle.SetComplement(planCandidateSetPrimaryKeyColumns, planCandidateSetColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*PlanCandidateSet{&b, &c} {
		err = a.SetPlanCandidateSet(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.PlanCandidateSet != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.PlanCandidateSetCategories[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if !queries.Equal(a.PlanCandidateSetID, x.ID) {
			t.Error("foreign key was wrong value", a.PlanCandidateSetID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.PlanCandidateSetID))
		reflect.Indirect(reflect.ValueOf(&a.PlanCandidateSetID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if !queries.Equal(a.PlanCandidateSetID, x.ID) {
			t.Error("foreign key was wrong value", a.PlanCandidateSetID, x.ID)
		}
	}
}

func testPlanCandidateSetCategoryToOneRemoveOpPlanCandidateSetUsingPlanCandidateSet(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a PlanCandidateSetCategory
	var b PlanCandidateSet

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, planCandidateSetCategoryDBTypes, false, strmangle.SetComplement(planCandidateSetCategoryPrimaryKeyColumns, planCandidateSetCategoryColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, planCandidateSetDBTypes, false, strmangle.SetComplement(planCandidateSetPrimaryKeyColumns, planCandidateSetColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = a.SetPlanCandidateSet(ctx, tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemovePlanCandidateSet(ctx, tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.PlanCandidateSet().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.PlanCandidateSet != nil {
		t.Error("R struct entry should be nil")
	}

	if !queries.IsValuerNil(a.PlanCandidateSetID) {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.PlanCandidateSetCategories) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testPlanCandidateSetCategoriesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &PlanCandidateSetCategory{}
	if err = randomize.Struct(seed, o, planCandidateSetCategoryDBTypes, true, planCandidateSetCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PlanCandidateSetCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testPlanCandidateSetCategoriesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &PlanCandidateSetCategory{}
	if err = randomize.Struct(seed, o, planCandidateSetCategoryDBTypes, true, planCandidateSetCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PlanCandidateSetCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := PlanCandidateSetCategorySlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testPlanCandidateSetCategoriesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &PlanCandidateSetCategory{}
	if err = randomize.Struct(seed, o, planCandidateSetCategoryDBTypes, true, planCandidateSetCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PlanCandidateSetCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := PlanCandidateSetCategories().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	planCandidateSetCategoryDBTypes = map[string]string{`ID`: `char`, `PlanCandidateSetID`: `char`, `Category`: `varchar`, `IsSelected`: `tinyint`, `CreatedAt`: `timestamp`, `UpdatedAt`: `timestamp`}
	_                               = bytes.MinRead
)

func testPlanCandidateSetCategoriesUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(planCandidateSetCategoryPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(planCandidateSetCategoryAllColumns) == len(planCandidateSetCategoryPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &PlanCandidateSetCategory{}
	if err = randomize.Struct(seed, o, planCandidateSetCategoryDBTypes, true, planCandidateSetCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PlanCandidateSetCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := PlanCandidateSetCategories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, planCandidateSetCategoryDBTypes, true, planCandidateSetCategoryPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize PlanCandidateSetCategory struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testPlanCandidateSetCategoriesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(planCandidateSetCategoryAllColumns) == len(planCandidateSetCategoryPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &PlanCandidateSetCategory{}
	if err = randomize.Struct(seed, o, planCandidateSetCategoryDBTypes, true, planCandidateSetCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize PlanCandidateSetCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := PlanCandidateSetCategories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, planCandidateSetCategoryDBTypes, true, planCandidateSetCategoryPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize PlanCandidateSetCategory struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(planCandidateSetCategoryAllColumns, planCandidateSetCategoryPrimaryKeyColumns) {
		fields = planCandidateSetCategoryAllColumns
	} else {
		fields = strmangle.SetComplement(
			planCandidateSetCategoryAllColumns,
			planCandidateSetCategoryPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := PlanCandidateSetCategorySlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testPlanCandidateSetCategoriesUpsert(t *testing.T) {
	t.Parallel()

	if len(planCandidateSetCategoryAllColumns) == len(planCandidateSetCategoryPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}
	if len(mySQLPlanCandidateSetCategoryUniqueColumns) == 0 {
		t.Skip("Skipping table with no unique columns to conflict on")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := PlanCandidateSetCategory{}
	if err = randomize.Struct(seed, &o, planCandidateSetCategoryDBTypes, false); err != nil {
		t.Errorf("Unable to randomize PlanCandidateSetCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert PlanCandidateSetCategory: %s", err)
	}

	count, err := PlanCandidateSetCategories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, planCandidateSetCategoryDBTypes, false, planCandidateSetCategoryPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize PlanCandidateSetCategory struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert PlanCandidateSetCategory: %s", err)
	}

	count, err = PlanCandidateSetCategories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
