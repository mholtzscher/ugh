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
