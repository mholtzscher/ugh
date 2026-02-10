package compile_test

import (
	"testing"
	"time"

	"github.com/mholtzscher/ugh/internal/nlp"
	"github.com/mholtzscher/ugh/internal/nlp/compile"
)

func TestBuildCreatePlan(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`add buy milk due:tomorrow #home @errands waiting:alex`, nlp.ParseOptions{
		Now: time.Date(2026, 2, 8, 10, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Parse(create) error = %v", err)
	}

	plan, err := compile.Build(parsed, compile.BuildOptions{Now: time.Date(2026, 2, 8, 10, 0, 0, 0, time.UTC)})
	if err != nil {
		t.Fatalf("Build(create) error = %v", err)
	}
	if plan.Create == nil {
		t.Fatal("create request is nil")
	}
	if plan.Create.Title != "buy milk" {
		t.Fatalf("title = %q, want %q", plan.Create.Title, "buy milk")
	}
	if plan.Create.DueOn != "2026-02-09" {
		t.Fatalf("due = %q, want %q", plan.Create.DueOn, "2026-02-09")
	}
	if plan.Create.WaitingFor != "alex" {
		t.Fatalf("waiting_for = %q, want %q", plan.Create.WaitingFor, "alex")
	}
	if len(plan.Create.Projects) != 1 || plan.Create.Projects[0] != "home" {
		t.Fatalf("projects = %#v, want [home]", plan.Create.Projects)
	}
	if len(plan.Create.Contexts) != 1 || plan.Create.Contexts[0] != "errands" {
		t.Fatalf("contexts = %#v, want [errands]", plan.Create.Contexts)
	}
}

func TestBuildUpdatePlanResolvesSelectedTarget(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set selected state:now +project:work !due`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(update) error = %v", err)
	}

	id := int64(42)
	plan, err := compile.Build(parsed, compile.BuildOptions{SelectedTaskID: &id})
	if err != nil {
		t.Fatalf("Build(update) error = %v", err)
	}
	if plan.Update == nil {
		t.Fatal("update request is nil")
	}
	if plan.Update.ID != 42 {
		t.Fatalf("update id = %d, want 42", plan.Update.ID)
	}
	if plan.Update.State == nil || *plan.Update.State != "now" {
		t.Fatalf("state = %#v, want now", plan.Update.State)
	}
	if !plan.Update.ClearDueOn {
		t.Fatal("ClearDueOn = false, want true")
	}
	if len(plan.Update.AddProjects) != 1 || plan.Update.AddProjects[0] != "work" {
		t.Fatalf("add projects = %#v, want [work]", plan.Update.AddProjects)
	}
}

func TestBuildFilterPlanBuildsBooleanExpression(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`find state:todo and project:work and text:paper`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(filter) error = %v", err)
	}

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	if err != nil {
		t.Fatalf("Build(filter) error = %v", err)
	}
	if plan.Filter == nil || plan.Filter.Filter == nil {
		t.Fatal("filter expression is nil")
	}

	binary, ok := plan.Filter.Filter.(nlp.FilterBinary)
	if !ok || binary.Op != nlp.FilterAnd {
		t.Fatalf("filter type = %T, want FilterBinary(AND)", plan.Filter.Filter)
	}

	left, ok := binary.Left.(nlp.Predicate)
	if !ok || left.Kind != nlp.PredState || left.Text != "inbox" {
		t.Fatalf("left predicate = %#v, want state:inbox", binary.Left)
	}

	right, ok := binary.Right.(nlp.FilterBinary)
	if !ok || right.Op != nlp.FilterAnd {
		t.Fatalf("right expression = %#v, want nested AND", binary.Right)
	}

	project, ok := right.Left.(nlp.Predicate)
	if !ok || project.Kind != nlp.PredProject || project.Text != "work" {
		t.Fatalf("project predicate = %#v, want project:work", right.Left)
	}

	search, ok := right.Right.(nlp.Predicate)
	if !ok || search.Kind != nlp.PredText || search.Text != "paper" {
		t.Fatalf("search predicate = %#v, want text:paper", right.Right)
	}
}

func TestBuildFilterPlanNormalizesDueDate(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 10, 10, 0, 0, 0, time.UTC)
	parsed, err := nlp.Parse(`find due:tomorrow`, nlp.ParseOptions{Now: now})
	if err != nil {
		t.Fatalf("Parse(filter) error = %v", err)
	}

	plan, err := compile.Build(parsed, compile.BuildOptions{Now: now})
	if err != nil {
		t.Fatalf("Build(filter) error = %v", err)
	}

	pred, ok := plan.Filter.Filter.(nlp.Predicate)
	if !ok {
		t.Fatalf("filter type = %T, want Predicate", plan.Filter.Filter)
	}
	if pred.Kind != nlp.PredDue || pred.Text != "2026-02-11" {
		t.Fatalf("due predicate = %#v, want due:2026-02-11", pred)
	}
}

func TestBuildFilterPlanSupportsOrAndNot(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`find state:now or not state:done`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(filter) error = %v", err)
	}

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	if err != nil {
		t.Fatalf("Build(filter) error = %v", err)
	}

	binary, ok := plan.Filter.Filter.(nlp.FilterBinary)
	if !ok || binary.Op != nlp.FilterOr {
		t.Fatalf("filter type = %T, want FilterBinary(OR)", plan.Filter.Filter)
	}

	rightNot, ok := binary.Right.(nlp.FilterNot)
	if !ok {
		t.Fatalf("right expression = %T, want FilterNot", binary.Right)
	}

	pred, ok := rightNot.Expr.(nlp.Predicate)
	if !ok || pred.Kind != nlp.PredState || pred.Text != "done" {
		t.Fatalf("not predicate = %#v, want state:done", rightNot.Expr)
	}
}

func TestBuildFilterPlanInvalidIDReturnsError(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`find id:abc`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(filter) error = %v", err)
	}

	_, err = compile.Build(parsed, compile.BuildOptions{})
	if err == nil {
		t.Fatal("Build(filter) error = nil, want invalid id error")
	}
}

func TestNormalizeFilterExpr_AllowsDueSetPredicate(t *testing.T) {
	t.Parallel()

	expr, err := compile.NormalizeFilterExpr(
		nlp.Predicate{Kind: nlp.PredDue, Text: ""},
		compile.BuildOptions{Now: time.Date(2026, 2, 10, 10, 0, 0, 0, time.UTC)},
	)
	if err != nil {
		t.Fatalf("NormalizeFilterExpr() error = %v", err)
	}

	pred, ok := expr.(nlp.Predicate)
	if !ok {
		t.Fatalf("expr type = %T, want Predicate", expr)
	}
	if pred.Kind != nlp.PredDue || pred.Text != "" {
		t.Fatalf("predicate = %#v, want due predicate with empty text", pred)
	}
}

func TestBuildFilterPlanNormalizesTodayDate(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 10, 10, 0, 0, 0, time.UTC)
	parsed, err := nlp.Parse(`find due:today`, nlp.ParseOptions{Now: now})
	if err != nil {
		t.Fatalf("Parse(filter) error = %v", err)
	}

	plan, err := compile.Build(parsed, compile.BuildOptions{Now: now})
	if err != nil {
		t.Fatalf("Build(filter) error = %v", err)
	}

	pred, ok := plan.Filter.Filter.(nlp.Predicate)
	if !ok {
		t.Fatalf("filter type = %T, want Predicate", plan.Filter.Filter)
	}
	if pred.Kind != nlp.PredDue || pred.Text != "2026-02-10" {
		t.Fatalf("due predicate = %#v, want due:2026-02-10", pred)
	}
}

func TestBuildFilterPlanNormalizesNextWeekDate(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 10, 10, 0, 0, 0, time.UTC)
	parsed, err := nlp.Parse(`find due:next-week`, nlp.ParseOptions{Now: now})
	if err != nil {
		t.Fatalf("Parse(filter) error = %v", err)
	}

	plan, err := compile.Build(parsed, compile.BuildOptions{Now: now})
	if err != nil {
		t.Fatalf("Build(filter) error = %v", err)
	}

	pred, ok := plan.Filter.Filter.(nlp.Predicate)
	if !ok {
		t.Fatalf("filter type = %T, want Predicate", plan.Filter.Filter)
	}
	if pred.Kind != nlp.PredDue || pred.Text != "2026-02-17" {
		t.Fatalf("due predicate = %#v, want due:2026-02-17", pred)
	}
}

func TestBuildFilterPlanNormalizesExplicitDate(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 10, 10, 0, 0, 0, time.UTC)
	parsed, err := nlp.Parse(`find due:2026-03-15`, nlp.ParseOptions{Now: now})
	if err != nil {
		t.Fatalf("Parse(filter) error = %v", err)
	}

	plan, err := compile.Build(parsed, compile.BuildOptions{Now: now})
	if err != nil {
		t.Fatalf("Build(filter) error = %v", err)
	}

	pred, ok := plan.Filter.Filter.(nlp.Predicate)
	if !ok {
		t.Fatalf("filter type = %T, want Predicate", plan.Filter.Filter)
	}
	if pred.Kind != nlp.PredDue || pred.Text != "2026-03-15" {
		t.Fatalf("due predicate = %#v, want due:2026-03-15", pred)
	}
}

func TestBuildFilterPlanInvalidDueDateReturnsError(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 10, 10, 0, 0, 0, time.UTC)
	parsed, err := nlp.Parse(`find due:invalid-date`, nlp.ParseOptions{Now: now})
	if err != nil {
		t.Fatalf("Parse(filter) error = %v", err)
	}

	_, err = compile.Build(parsed, compile.BuildOptions{Now: now})
	if err == nil {
		t.Fatal("Build(filter) error = nil, want invalid date error")
	}
}

func TestBuildFilterPlanNormalizesTodoState(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`find state:todo`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(filter) error = %v", err)
	}

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	if err != nil {
		t.Fatalf("Build(filter) error = %v", err)
	}

	pred, ok := plan.Filter.Filter.(nlp.Predicate)
	if !ok {
		t.Fatalf("filter type = %T, want Predicate", plan.Filter.Filter)
	}
	if pred.Kind != nlp.PredState || pred.Text != "inbox" {
		t.Fatalf("state predicate = %#v, want state:inbox", pred)
	}
}

func TestBuildFilterPlanInvalidStateReturnsError(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`find state:invalid`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(filter) error = %v", err)
	}

	_, err = compile.Build(parsed, compile.BuildOptions{})
	if err == nil {
		t.Fatal("Build(filter) error = nil, want invalid state error")
	}
}

func TestBuildUpdatePlanRemovesProjectsAndContexts(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 -project:old -context:deprecated`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(update) error = %v", err)
	}

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	if err != nil {
		t.Fatalf("Build(update) error = %v", err)
	}
	if plan.Update == nil {
		t.Fatal("update request is nil")
	}
	if plan.Update.ID != 42 {
		t.Fatalf("update id = %d, want 42", plan.Update.ID)
	}
	if len(plan.Update.RemoveProjects) != 1 || plan.Update.RemoveProjects[0] != "old" {
		t.Fatalf("remove projects = %#v, want [old]", plan.Update.RemoveProjects)
	}
	if len(plan.Update.RemoveContexts) != 1 || plan.Update.RemoveContexts[0] != "deprecated" {
		t.Fatalf("remove contexts = %#v, want [deprecated]", plan.Update.RemoveContexts)
	}
}

func TestBuildUpdatePlanRemovesMetaKeys(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 -meta:key1 -meta:key2:value`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(update) error = %v", err)
	}

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	if err != nil {
		t.Fatalf("Build(update) error = %v", err)
	}
	if plan.Update == nil {
		t.Fatal("update request is nil")
	}
	if len(plan.Update.RemoveMetaKeys) != 2 {
		t.Fatalf("remove meta keys count = %d, want 2", len(plan.Update.RemoveMetaKeys))
	}
	if plan.Update.RemoveMetaKeys[0] != "key1" {
		t.Fatalf("first meta key = %q, want key1", plan.Update.RemoveMetaKeys[0])
	}
	if plan.Update.RemoveMetaKeys[1] != "key2" {
		t.Fatalf("second meta key = %q, want key2", plan.Update.RemoveMetaKeys[1])
	}
}

func TestBuildUpdatePlanClearsDueAndWaiting(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 !due !waiting`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(update) error = %v", err)
	}

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	if err != nil {
		t.Fatalf("Build(update) error = %v", err)
	}
	if plan.Update == nil {
		t.Fatal("update request is nil")
	}
	if !plan.Update.ClearDueOn {
		t.Fatal("ClearDueOn = false, want true")
	}
	if plan.Update.DueOn != nil {
		t.Fatal("DueOn should be nil when clearing")
	}
	if !plan.Update.ClearWaitingFor {
		t.Fatal("ClearWaitingFor = false, want true")
	}
	if plan.Update.WaitingFor != nil {
		t.Fatal("WaitingFor should be nil when clearing")
	}
}

func TestBuildUpdatePlanClearsNotes(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 !notes`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(update) error = %v", err)
	}

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	if err != nil {
		t.Fatalf("Build(update) error = %v", err)
	}
	if plan.Update == nil {
		t.Fatal("update request is nil")
	}
	if plan.Update.Notes == nil || *plan.Update.Notes != "" {
		t.Fatalf("notes = %v, want empty string", plan.Update.Notes)
	}
}

func TestBuildCreatePlanWithMultipleProjectsAndContexts(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`add task #project1 #project2 @context1 @context2`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(create) error = %v", err)
	}

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	if err != nil {
		t.Fatalf("Build(create) error = %v", err)
	}
	if plan.Create == nil {
		t.Fatal("create request is nil")
	}
	if len(plan.Create.Projects) != 2 {
		t.Fatalf("projects count = %d, want 2", len(plan.Create.Projects))
	}
	if len(plan.Create.Contexts) != 2 {
		t.Fatalf("contexts count = %d, want 2", len(plan.Create.Contexts))
	}
}

func TestBuildCreatePlanClearsFields(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`add task !notes !due !waiting !projects !contexts`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(create) error = %v", err)
	}

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	if err != nil {
		t.Fatalf("Build(create) error = %v", err)
	}
	if plan.Create == nil {
		t.Fatal("create request is nil")
	}
	if plan.Create.Notes != "" {
		t.Fatalf("notes = %q, want empty", plan.Create.Notes)
	}
	if plan.Create.DueOn != "" {
		t.Fatalf("due = %q, want empty", plan.Create.DueOn)
	}
	if plan.Create.WaitingFor != "" {
		t.Fatalf("waiting = %q, want empty", plan.Create.WaitingFor)
	}
	if len(plan.Create.Projects) != 0 {
		t.Fatalf("projects = %#v, want empty", plan.Create.Projects)
	}
	if len(plan.Create.Contexts) != 0 {
		t.Fatalf("contexts = %#v, want empty", plan.Create.Contexts)
	}
}

func TestBuildCreatePlanRemovesNotSupported(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`add task -project:foo`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(create) error = %v", err)
	}

	_, err = compile.Build(parsed, compile.BuildOptions{})
	if err == nil {
		t.Fatal("Build(create) error = nil, want remove not supported error")
	}
}

func TestBuildUpdatePlanSetsMetaValues(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 meta:key:value`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(update) error = %v", err)
	}

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	if err != nil {
		t.Fatalf("Build(update) error = %v", err)
	}
	if plan.Update == nil {
		t.Fatal("update request is nil")
	}
	if len(plan.Update.SetMeta) != 1 {
		t.Fatalf("set meta count = %d, want 1", len(plan.Update.SetMeta))
	}
	if plan.Update.SetMeta["key"] != "value" {
		t.Fatalf("meta[key] = %q, want value", plan.Update.SetMeta["key"])
	}
}

func TestBuildUpdatePlanAddsMetaValues(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 +meta:priority:high`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(update) error = %v", err)
	}

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	if err != nil {
		t.Fatalf("Build(update) error = %v", err)
	}
	if plan.Update == nil {
		t.Fatal("update request is nil")
	}
	if len(plan.Update.SetMeta) != 1 {
		t.Fatalf("set meta count = %d, want 1", len(plan.Update.SetMeta))
	}
	if plan.Update.SetMeta["priority"] != "high" {
		t.Fatalf("meta[priority] = %q, want high", plan.Update.SetMeta["priority"])
	}
}

func TestBuildCreatePlanSetsMetaValues(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`add task meta:priority:high`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(create) error = %v", err)
	}

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	if err != nil {
		t.Fatalf("Build(create) error = %v", err)
	}
	if plan.Create == nil {
		t.Fatal("create request is nil")
	}
	if len(plan.Create.Meta) != 1 {
		t.Fatalf("meta count = %d, want 1", len(plan.Create.Meta))
	}
}

func TestBuildUpdatePlanRequiresSelectedTaskID(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set selected state:done`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(update) error = %v", err)
	}

	_, err = compile.Build(parsed, compile.BuildOptions{})
	if err == nil {
		t.Fatal("Build(update) error = nil, want selected task ID required error")
	}
}

func TestBuildUpdatePlanRemovesWithMultipleValues(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 -project:p1,p2,p3`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(update) error = %v", err)
	}

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	if err != nil {
		t.Fatalf("Build(update) error = %v", err)
	}
	if plan.Update == nil {
		t.Fatal("update request is nil")
	}
	if len(plan.Update.RemoveProjects) != 3 {
		t.Fatalf("remove projects count = %d, want 3", len(plan.Update.RemoveProjects))
	}
}

func TestBuildCreatePlanDeduplicatesProjectsAndContexts(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`add task #foo #foo @bar @bar`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(create) error = %v", err)
	}

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	if err != nil {
		t.Fatalf("Build(create) error = %v", err)
	}
	if plan.Create == nil {
		t.Fatal("create request is nil")
	}
	if len(plan.Create.Projects) != 1 {
		t.Fatalf("projects count = %d, want 1", len(plan.Create.Projects))
	}
	if len(plan.Create.Contexts) != 1 {
		t.Fatalf("contexts count = %d, want 1", len(plan.Create.Contexts))
	}
}

func TestBuildUpdatePlanDeduplicatesRemovals(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 -project:foo -project:foo -context:bar -context:bar`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(update) error = %v", err)
	}

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	if err != nil {
		t.Fatalf("Build(update) error = %v", err)
	}
	if plan.Update == nil {
		t.Fatal("update request is nil")
	}
	if len(plan.Update.RemoveProjects) != 1 {
		t.Fatalf("remove projects count = %d, want 1", len(plan.Update.RemoveProjects))
	}
	if len(plan.Update.RemoveContexts) != 1 {
		t.Fatalf("remove contexts count = %d, want 1", len(plan.Update.RemoveContexts))
	}
}

func TestNormalizeFilterExprHandlesNil(t *testing.T) {
	t.Parallel()

	expr, err := compile.NormalizeFilterExpr(nil, compile.BuildOptions{})
	if err != nil {
		t.Fatalf("NormalizeFilterExpr() error = %v", err)
	}
	if expr != nil {
		t.Fatalf("expr = %v, want nil", expr)
	}
}

func TestBuildCreatePlanSetsState(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`add task state:now`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(create) error = %v", err)
	}

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	if err != nil {
		t.Fatalf("Build(create) error = %v", err)
	}
	if plan.Create == nil {
		t.Fatal("create request is nil")
	}
	if plan.Create.State != "now" {
		t.Fatalf("state = %q, want now", plan.Create.State)
	}
}

func TestBuildCreatePlanDefaultsToInboxState(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`add task`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(create) error = %v", err)
	}

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	if err != nil {
		t.Fatalf("Build(create) error = %v", err)
	}
	if plan.Create == nil {
		t.Fatal("create request is nil")
	}
	if plan.Create.State != "inbox" {
		t.Fatalf("state = %q, want inbox", plan.Create.State)
	}
}

func TestBuildFilterPlanHandlesIDFilter(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`find id:123`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(filter) error = %v", err)
	}

	plan, err := compile.Build(parsed, compile.BuildOptions{})
	if err != nil {
		t.Fatalf("Build(filter) error = %v", err)
	}
	if plan.Filter == nil {
		t.Fatal("filter request is nil")
	}

	pred, ok := plan.Filter.Filter.(nlp.Predicate)
	if !ok {
		t.Fatalf("filter type = %T, want Predicate", plan.Filter.Filter)
	}
	if pred.Kind != nlp.PredID || pred.Text != "123" {
		t.Fatalf("id predicate = %#v, want id:123", pred)
	}
}

func TestBuildUpdatePlanRejectsSetProjects(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 projects:foo`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(update) error = %v", err)
	}

	_, err = compile.Build(parsed, compile.BuildOptions{})
	if err == nil {
		t.Fatal("Build(update) error = nil, want set projects not supported error")
	}
}

func TestBuildUpdatePlanRejectsSetContexts(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 contexts:foo`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(update) error = %v", err)
	}

	_, err = compile.Build(parsed, compile.BuildOptions{})
	if err == nil {
		t.Fatal("Build(update) error = nil, want set contexts not supported error")
	}
}

func TestBuildUpdatePlanRejectsClearProjects(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 !projects`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(update) error = %v", err)
	}

	_, err = compile.Build(parsed, compile.BuildOptions{})
	if err == nil {
		t.Fatal("Build(update) error = nil, want clear projects not supported error")
	}
}

func TestBuildUpdatePlanRejectsClearContexts(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 !contexts`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(update) error = %v", err)
	}

	_, err = compile.Build(parsed, compile.BuildOptions{})
	if err == nil {
		t.Fatal("Build(update) error = nil, want clear contexts not supported error")
	}
}

func TestBuildUpdatePlanRejectsClearMeta(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 !meta`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(update) error = %v", err)
	}

	_, err = compile.Build(parsed, compile.BuildOptions{})
	if err == nil {
		t.Fatal("Build(update) error = nil, want clear meta not supported error")
	}
}

func TestBuildCreatePlanRejectsClearState(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`add task !state`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(create) error = %v", err)
	}

	_, err = compile.Build(parsed, compile.BuildOptions{})
	if err == nil {
		t.Fatal("Build(create) error = nil, want clear state not supported error")
	}
}

func TestBuildCreatePlanRejectsClearTitle(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`add task !title`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(create) error = %v", err)
	}

	_, err = compile.Build(parsed, compile.BuildOptions{})
	if err == nil {
		t.Fatal("Build(create) error = nil, want clear title not supported error")
	}
}

func TestBuildUpdatePlanRejectsClearTitle(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 !title`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(update) error = %v", err)
	}

	_, err = compile.Build(parsed, compile.BuildOptions{})
	if err == nil {
		t.Fatal("Build(update) error = nil, want clear title not supported error")
	}
}

func TestBuildUpdatePlanRejectsClearState(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`set 42 !state`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(update) error = %v", err)
	}

	_, err = compile.Build(parsed, compile.BuildOptions{})
	if err == nil {
		t.Fatal("Build(update) error = nil, want clear state not supported error")
	}
}
