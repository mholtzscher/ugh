package compile_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mholtzscher/ugh/internal/nlp"
	"github.com/mholtzscher/ugh/internal/nlp/compile"
)

func TestBuildCreatePlan(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`add buy milk due:tomorrow #home @errands waiting:alex`, nlp.ParseOptions{
		Now: time.Date(2026, 2, 8, 10, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err, "Parse(create) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{Now: time.Date(2026, 2, 8, 10, 0, 0, 0, time.UTC)})
	require.NoError(t, err, "Build(create) error")
	require.NotNil(t, plan.Create, "create request is nil")
	require.Equal(t, "buy milk", plan.Create.Title, "title mismatch")
	require.Equal(t, "2026-02-09", plan.Create.DueOn, "due mismatch")
	require.Equal(t, "alex", plan.Create.WaitingFor, "waiting_for mismatch")
	require.Equal(t, []string{"home"}, plan.Create.Projects, "projects mismatch")
	require.Equal(t, []string{"errands"}, plan.Create.Contexts, "contexts mismatch")
}

func TestBuildCreatePlanNormalizesNaturalLanguageDueDate(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 10, 10, 0, 0, 0, time.UTC)
	parsed, err := nlp.Parse(`add buy milk due:next monday`, nlp.ParseOptions{Now: now})
	require.NoError(t, err, "Parse(create) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{Now: now})
	require.NoError(t, err, "Build(create) error")
	require.NotNil(t, plan.Create, "create request is nil")
	require.Equal(t, "2026-02-16", plan.Create.DueOn, "due mismatch")
}

func TestBuildUpdatePlanResolvesSelectedTarget(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set selected state:now +project:work !due`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(update) error")

	id := int64(42)
	plan, err := compile.Build(parsed, compile.BuildOptions{SelectedTaskID: &id})
	require.NoError(t, err, "Build(update) error")
	require.NotNil(t, plan.Update, "update request is nil")
	require.Equal(t, int64(42), plan.Update.ID, "update id mismatch")
	require.NotNil(t, plan.Update.State, "state is nil")
	require.Equal(t, "now", *plan.Update.State, "state mismatch")
	require.True(t, plan.Update.ClearDueOn, "ClearDueOn should be true")
	require.Equal(t, []string{"work"}, plan.Update.AddProjects, "add projects mismatch")
}

func TestBuildFilterPlanBuildsBooleanExpression(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`find state:todo and project:work and text:paper`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(filter) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	require.NoError(t, err, "Build(filter) error")
	require.NotNil(t, plan.Filter, "filter request is nil")
	require.NotNil(t, plan.Filter.Filter, "filter expression is nil")

	binary, ok := plan.Filter.Filter.(nlp.FilterBinary)
	require.True(t, ok, "filter type should be FilterBinary, got %T", plan.Filter.Filter)
	require.Equal(t, nlp.FilterAnd, binary.Op, "filter should be AND")

	left, ok := binary.Left.(nlp.Predicate)
	require.True(t, ok, "left should be Predicate, got %T", binary.Left)
	require.Equal(t, nlp.PredState, left.Kind, "left predicate kind mismatch")
	require.Equal(t, "inbox", left.Text, "left predicate text mismatch")

	right, ok := binary.Right.(nlp.FilterBinary)
	require.True(t, ok, "right should be FilterBinary, got %T", binary.Right)
	require.Equal(t, nlp.FilterAnd, right.Op, "right should be AND")

	project, ok := right.Left.(nlp.Predicate)
	require.True(t, ok, "project should be Predicate, got %T", right.Left)
	require.Equal(t, nlp.PredProject, project.Kind, "project predicate kind mismatch")
	require.Equal(t, "work", project.Text, "project predicate text mismatch")

	search, ok := right.Right.(nlp.Predicate)
	require.True(t, ok, "search should be Predicate, got %T", right.Right)
	require.Equal(t, nlp.PredText, search.Kind, "search predicate kind mismatch")
	require.Equal(t, "paper", search.Text, "search predicate text mismatch")
}

func TestBuildFilterPlanNormalizesDueDate(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 10, 10, 0, 0, 0, time.UTC)
	parsed, err := nlp.Parse(`find due:tomorrow`, nlp.ParseOptions{Now: now})
	require.NoError(t, err, "Parse(filter) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{Now: now})
	require.NoError(t, err, "Build(filter) error")

	pred, ok := plan.Filter.Filter.(nlp.Predicate)
	require.True(t, ok, "filter type should be Predicate, got %T", plan.Filter.Filter)
	require.Equal(t, nlp.PredDue, pred.Kind, "predicate kind mismatch")
	require.Equal(t, "2026-02-11", pred.Text, "due predicate text mismatch")
}

func TestBuildFilterPlanSupportsOrAndNot(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`find state:now or not state:done`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(filter) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	require.NoError(t, err, "Build(filter) error")

	binary, ok := plan.Filter.Filter.(nlp.FilterBinary)
	require.True(t, ok, "filter type should be FilterBinary, got %T", plan.Filter.Filter)
	require.Equal(t, nlp.FilterOr, binary.Op, "filter should be OR")

	rightNot, ok := binary.Right.(nlp.FilterNot)
	require.True(t, ok, "right should be FilterNot, got %T", binary.Right)

	pred, ok := rightNot.Expr.(nlp.Predicate)
	require.True(t, ok, "not predicate should be Predicate, got %T", rightNot.Expr)
	require.Equal(t, nlp.PredState, pred.Kind, "predicate kind mismatch")
	require.Equal(t, "done", pred.Text, "predicate text mismatch")
}

func TestBuildFilterPlanInvalidIDReturnsError(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`find id:abc`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(filter) error")

	_, err = compile.Build(parsed, compile.BuildOptions{})
	require.Error(t, err, "Build(filter) should return invalid id error")
}

func TestNormalizeFilterExpr_AllowsWildcardSetPredicate(t *testing.T) {
	t.Parallel()

	expr, err := compile.NormalizeFilterExpr(
		nlp.Predicate{Kind: nlp.PredDue, Text: nlp.FilterWildcard},
		compile.BuildOptions{Now: time.Date(2026, 2, 10, 10, 0, 0, 0, time.UTC)},
	)
	require.NoError(t, err, "NormalizeFilterExpr() error")

	pred, ok := expr.(nlp.Predicate)
	require.True(t, ok, "expr type should be Predicate, got %T", expr)
	require.Equal(t, nlp.PredDue, pred.Kind, "predicate kind mismatch")
	require.Equal(t, nlp.FilterWildcard, pred.Text, "predicate text mismatch")
}

func TestNormalizeFilterExpr_AllowsWildcardProjectPredicate(t *testing.T) {
	t.Parallel()

	expr, err := compile.NormalizeFilterExpr(
		nlp.Predicate{Kind: nlp.PredProject, Text: nlp.FilterWildcard},
		compile.BuildOptions{Now: time.Date(2026, 2, 10, 10, 0, 0, 0, time.UTC)},
	)
	require.NoError(t, err, "NormalizeFilterExpr() error")

	pred, ok := expr.(nlp.Predicate)
	require.True(t, ok, "expr type should be Predicate, got %T", expr)
	require.Equal(t, nlp.PredProject, pred.Kind, "predicate kind mismatch")
	require.Equal(t, nlp.FilterWildcard, pred.Text, "predicate text mismatch")
}

func TestNormalizeFilterExpr_WildcardStateReturnsError(t *testing.T) {
	t.Parallel()

	_, err := compile.NormalizeFilterExpr(
		nlp.Predicate{Kind: nlp.PredState, Text: nlp.FilterWildcard},
		compile.BuildOptions{Now: time.Date(2026, 2, 10, 10, 0, 0, 0, time.UTC)},
	)
	require.Error(t, err, "NormalizeFilterExpr() should reject wildcard for state")
	assert.Contains(t, err.Error(), "wildcard", "error should mention wildcard")
}

func TestBuildFilterPlanNormalizesTodayDate(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 10, 10, 0, 0, 0, time.UTC)
	parsed, err := nlp.Parse(`find due:today`, nlp.ParseOptions{Now: now})
	require.NoError(t, err, "Parse(filter) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{Now: now})
	require.NoError(t, err, "Build(filter) error")

	pred, ok := plan.Filter.Filter.(nlp.Predicate)
	require.True(t, ok, "filter type should be Predicate, got %T", plan.Filter.Filter)
	require.Equal(t, nlp.PredDue, pred.Kind, "predicate kind mismatch")
	require.Equal(t, "2026-02-10", pred.Text, "due predicate text mismatch")
}

func TestBuildFilterPlanNormalizesNextWeekDate(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 10, 10, 0, 0, 0, time.UTC)
	parsed, err := nlp.Parse(`find due:next-week`, nlp.ParseOptions{Now: now})
	require.NoError(t, err, "Parse(filter) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{Now: now})
	require.NoError(t, err, "Build(filter) error")

	pred, ok := plan.Filter.Filter.(nlp.Predicate)
	require.True(t, ok, "filter type should be Predicate, got %T", plan.Filter.Filter)
	require.Equal(t, nlp.PredDue, pred.Kind, "predicate kind mismatch")
	require.Equal(t, "2026-02-17", pred.Text, "due predicate text mismatch")
}

func TestBuildFilterPlanNormalizesNaturalLanguageDate(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 10, 10, 0, 0, 0, time.UTC)
	parsed, err := nlp.Parse(`find due:next monday`, nlp.ParseOptions{Now: now})
	require.NoError(t, err, "Parse(filter) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{Now: now})
	require.NoError(t, err, "Build(filter) error")

	pred, ok := plan.Filter.Filter.(nlp.Predicate)
	require.True(t, ok, "filter type should be Predicate, got %T", plan.Filter.Filter)
	require.Equal(t, nlp.PredDue, pred.Kind, "predicate kind mismatch")
	require.Equal(t, "2026-02-16", pred.Text, "due predicate text mismatch")
}

func TestBuildFilterPlanNaturalLanguageDateStopsAtAndOperator(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 10, 10, 0, 0, 0, time.UTC)
	parsed, err := nlp.Parse(`find due:in 3 days and state:now`, nlp.ParseOptions{Now: now})
	require.NoError(t, err, "Parse(filter) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{Now: now})
	require.NoError(t, err, "Build(filter) error")

	binary, ok := plan.Filter.Filter.(nlp.FilterBinary)
	require.True(t, ok, "filter type should be FilterBinary, got %T", plan.Filter.Filter)
	require.Equal(t, nlp.FilterAnd, binary.Op, "filter should be AND")

	duePred, ok := binary.Left.(nlp.Predicate)
	require.True(t, ok, "left should be Predicate, got %T", binary.Left)
	require.Equal(t, nlp.PredDue, duePred.Kind, "left predicate kind mismatch")
	require.Equal(t, "2026-02-13", duePred.Text, "left due predicate text mismatch")

	statePred, ok := binary.Right.(nlp.Predicate)
	require.True(t, ok, "right should be Predicate, got %T", binary.Right)
	require.Equal(t, nlp.PredState, statePred.Kind, "right predicate kind mismatch")
	require.Equal(t, "now", statePred.Text, "right predicate text mismatch")
}

func TestBuildFilterPlanNormalizesExplicitDate(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 10, 10, 0, 0, 0, time.UTC)
	parsed, err := nlp.Parse(`find due:2026-03-15`, nlp.ParseOptions{Now: now})
	require.NoError(t, err, "Parse(filter) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{Now: now})
	require.NoError(t, err, "Build(filter) error")

	pred, ok := plan.Filter.Filter.(nlp.Predicate)
	require.True(t, ok, "filter type should be Predicate, got %T", plan.Filter.Filter)
	require.Equal(t, nlp.PredDue, pred.Kind, "predicate kind mismatch")
	require.Equal(t, "2026-03-15", pred.Text, "due predicate text mismatch")
}

func TestBuildFilterPlanInvalidDueDateReturnsError(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 10, 10, 0, 0, 0, time.UTC)
	parsed, err := nlp.Parse(`find due:invalid-date`, nlp.ParseOptions{Now: now})
	require.NoError(t, err, "Parse(filter) error")

	_, err = compile.Build(parsed, compile.BuildOptions{Now: now})
	require.Error(t, err, "Build(filter) should return invalid date error")
}

func TestBuildFilterPlanNormalizesTodoState(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`find state:todo`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(filter) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	require.NoError(t, err, "Build(filter) error")

	pred, ok := plan.Filter.Filter.(nlp.Predicate)
	require.True(t, ok, "filter type should be Predicate, got %T", plan.Filter.Filter)
	require.Equal(t, nlp.PredState, pred.Kind, "predicate kind mismatch")
	require.Equal(t, "inbox", pred.Text, "state predicate text mismatch")
}

func TestBuildFilterPlanInvalidStateReturnsError(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`find state:invalid`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(filter) error")

	_, err = compile.Build(parsed, compile.BuildOptions{})
	require.Error(t, err, "Build(filter) should return invalid state error")
}

func TestBuildUpdatePlanRemovesProjectsAndContexts(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 -project:old -context:deprecated`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(update) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	require.NoError(t, err, "Build(update) error")
	require.NotNil(t, plan.Update, "update request is nil")
	require.Equal(t, int64(42), plan.Update.ID, "update id mismatch")
	require.Equal(t, []string{"old"}, plan.Update.RemoveProjects, "remove projects mismatch")
	require.Equal(t, []string{"deprecated"}, plan.Update.RemoveContexts, "remove contexts mismatch")
}

func TestBuildUpdatePlanRemovesMetaKeys(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 -meta:key1 -meta:key2:value`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(update) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	require.NoError(t, err, "Build(update) error")
	require.NotNil(t, plan.Update, "update request is nil")
	require.Len(t, plan.Update.RemoveMetaKeys, 2, "remove meta keys count mismatch")
	require.Equal(t, "key1", plan.Update.RemoveMetaKeys[0], "first meta key mismatch")
	require.Equal(t, "key2", plan.Update.RemoveMetaKeys[1], "second meta key mismatch")
}

func TestBuildUpdatePlanClearsDueAndWaiting(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 !due !waiting`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(update) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	require.NoError(t, err, "Build(update) error")
	require.NotNil(t, plan.Update, "update request is nil")
	require.True(t, plan.Update.ClearDueOn, "ClearDueOn should be true")
	require.Nil(t, plan.Update.DueOn, "DueOn should be nil when clearing")
	require.True(t, plan.Update.ClearWaitingFor, "ClearWaitingFor should be true")
	require.Nil(t, plan.Update.WaitingFor, "WaitingFor should be nil when clearing")
}

func TestBuildUpdatePlanClearsNotes(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 !notes`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(update) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	require.NoError(t, err, "Build(update) error")
	require.NotNil(t, plan.Update, "update request is nil")
	require.NotNil(t, plan.Update.Notes, "notes should not be nil")
	require.Empty(t, *plan.Update.Notes, "notes should be empty string")
}

func TestBuildCreatePlanWithMultipleProjectsAndContexts(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`add task #project1 #project2 @context1 @context2`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(create) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	require.NoError(t, err, "Build(create) error")
	require.NotNil(t, plan.Create, "create request is nil")
	require.Len(t, plan.Create.Projects, 2, "projects count mismatch")
	require.Len(t, plan.Create.Contexts, 2, "contexts count mismatch")
}

func TestBuildCreatePlanClearsFields(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`add task !notes !due !waiting !projects !contexts`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(create) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	require.NoError(t, err, "Build(create) error")
	require.NotNil(t, plan.Create, "create request is nil")
	assert.Empty(t, plan.Create.Notes, "notes should be empty")
	assert.Empty(t, plan.Create.DueOn, "due should be empty")
	assert.Empty(t, plan.Create.WaitingFor, "waiting should be empty")
	assert.Empty(t, plan.Create.Projects, "projects should be empty")
	assert.Empty(t, plan.Create.Contexts, "contexts should be empty")
}

func TestBuildCreatePlanRemovesNotSupported(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`add task -project:foo`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(create) error")

	_, err = compile.Build(parsed, compile.BuildOptions{})
	require.Error(t, err, "Build(create) should return remove not supported error")
}

func TestBuildUpdatePlanSetsMetaValues(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 meta:key:value`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(update) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	require.NoError(t, err, "Build(update) error")
	require.NotNil(t, plan.Update, "update request is nil")
	require.Len(t, plan.Update.SetMeta, 1, "set meta count mismatch")
	require.Equal(t, "value", plan.Update.SetMeta["key"], "meta[key] value mismatch")
}

func TestBuildUpdatePlanAddsMetaValues(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 +meta:priority:high`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(update) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	require.NoError(t, err, "Build(update) error")
	require.NotNil(t, plan.Update, "update request is nil")
	require.Len(t, plan.Update.SetMeta, 1, "set meta count mismatch")
	require.Equal(t, "high", plan.Update.SetMeta["priority"], "meta[priority] value mismatch")
}

func TestBuildCreatePlanSetsMetaValues(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`add task meta:priority:high`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(create) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	require.NoError(t, err, "Build(create) error")
	require.NotNil(t, plan.Create, "create request is nil")
	require.Len(t, plan.Create.Meta, 1, "meta count mismatch")
}

func TestBuildUpdatePlanRequiresSelectedTaskID(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set selected state:done`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(update) error")

	_, err = compile.Build(parsed, compile.BuildOptions{})
	require.Error(t, err, "Build(update) should return selected task ID required error")
}

func TestBuildUpdatePlanRemovesWithMultipleValues(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 -project:p1,p2,p3`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(update) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	require.NoError(t, err, "Build(update) error")
	require.NotNil(t, plan.Update, "update request is nil")
	require.Len(t, plan.Update.RemoveProjects, 3, "remove projects count mismatch")
}

func TestBuildCreatePlanDeduplicatesProjectsAndContexts(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`add task #foo #foo @bar @bar`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(create) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	require.NoError(t, err, "Build(create) error")
	require.NotNil(t, plan.Create, "create request is nil")
	require.Len(t, plan.Create.Projects, 1, "projects count mismatch")
	require.Len(t, plan.Create.Contexts, 1, "contexts count mismatch")
}

func TestBuildUpdatePlanDeduplicatesRemovals(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 -project:foo -project:foo -context:bar -context:bar`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(update) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	require.NoError(t, err, "Build(update) error")
	require.NotNil(t, plan.Update, "update request is nil")
	require.Len(t, plan.Update.RemoveProjects, 1, "remove projects count mismatch")
	require.Len(t, plan.Update.RemoveContexts, 1, "remove contexts count mismatch")
}

func TestNormalizeFilterExprHandlesNil(t *testing.T) {
	t.Parallel()

	expr, err := compile.NormalizeFilterExpr(nil, compile.BuildOptions{})
	require.NoError(t, err, "NormalizeFilterExpr() error")
	require.Nil(t, expr, "expr should be nil")
}

func TestBuildCreatePlanSetsState(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`add task state:now`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(create) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	require.NoError(t, err, "Build(create) error")
	require.NotNil(t, plan.Create, "create request is nil")
	require.Equal(t, "now", plan.Create.State, "state mismatch")
}

func TestBuildCreatePlanDefaultsToInboxState(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`add task`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(create) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	require.NoError(t, err, "Build(create) error")
	require.NotNil(t, plan.Create, "create request is nil")
	require.Equal(t, "inbox", plan.Create.State, "state mismatch")
}

func TestBuildFilterPlanHandlesIDFilter(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`find id:123`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(filter) error")

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	require.NoError(t, err, "Build(filter) error")
	require.NotNil(t, plan.Filter, "filter request is nil")

	pred, ok := plan.Filter.Filter.(nlp.Predicate)
	require.True(t, ok, "filter type should be Predicate, got %T", plan.Filter.Filter)
	require.Equal(t, nlp.PredID, pred.Kind, "predicate kind mismatch")
	require.Equal(t, "123", pred.Text, "predicate text mismatch")
}

func TestBuildUpdatePlanRejectsSetProjects(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 projects:foo`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(update) error")

	_, err = compile.Build(parsed, compile.BuildOptions{})
	require.Error(t, err, "Build(update) should return set projects not supported error")
}

func TestBuildUpdatePlanRejectsSetContexts(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 contexts:foo`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(update) error")

	_, err = compile.Build(parsed, compile.BuildOptions{})
	require.Error(t, err, "Build(update) should return set contexts not supported error")
}

func TestBuildUpdatePlanRejectsClearProjects(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 !projects`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(update) error")

	_, err = compile.Build(parsed, compile.BuildOptions{})
	require.Error(t, err, "Build(update) should return clear projects not supported error")
}

func TestBuildUpdatePlanRejectsClearContexts(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 !contexts`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(update) error")

	_, err = compile.Build(parsed, compile.BuildOptions{})
	require.Error(t, err, "Build(update) should return clear contexts not supported error")
}

func TestBuildUpdatePlanRejectsClearMeta(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 !meta`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(update) error")

	_, err = compile.Build(parsed, compile.BuildOptions{})
	require.Error(t, err, "Build(update) should return clear meta not supported error")
}

func TestBuildCreatePlanRejectsClearState(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`add task !state`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(create) error")

	_, err = compile.Build(parsed, compile.BuildOptions{})
	require.Error(t, err, "Build(create) should return clear state not supported error")
}

func TestBuildCreatePlanRejectsClearTitle(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`add task !title`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(create) error")

	_, err = compile.Build(parsed, compile.BuildOptions{})
	require.Error(t, err, "Build(create) should return clear title not supported error")
}

func TestBuildUpdatePlanRejectsClearTitle(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 !title`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(update) error")

	_, err = compile.Build(parsed, compile.BuildOptions{})
	require.Error(t, err, "Build(update) should return clear title not supported error")
}

func TestBuildUpdatePlanRejectsClearState(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 !state`, nlp.ParseOptions{})
	require.NoError(t, err, "Parse(update) error")

	_, err = compile.Build(parsed, compile.BuildOptions{})
	require.Error(t, err, "Build(update) should return clear state not supported error")
}
