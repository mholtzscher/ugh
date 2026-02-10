package compile_test

import (
	"testing"
	"time"

	"github.com/mholtzscher/ugh/internal/nlp"
	"github.com/mholtzscher/ugh/internal/nlp/compile"
)

func TestBuildCreatePlan(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`add buy milk tomorrow #home @errands waiting:alex`, nlp.ParseOptions{
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

func TestBuildFilterPlan(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`find state:now and project:work and text:paper`, nlp.ParseOptions{})
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
	if len(plan.Filter.States) != 1 || plan.Filter.States[0] != "now" {
		t.Fatalf("states = %v, want [now]", plan.Filter.States)
	}
	if len(plan.Filter.Projects) != 1 || plan.Filter.Projects[0] != "work" {
		t.Fatalf("projects = %v, want [work]", plan.Filter.Projects)
	}
	if len(plan.Filter.Search) != 1 || plan.Filter.Search[0] != "paper" {
		t.Fatalf("search = %v, want [paper]", plan.Filter.Search)
	}
}

func TestBuildFilterPlanDueDateToday(t *testing.T) {
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
	if plan.Filter == nil {
		t.Fatal("filter request is nil")
	}
	if !plan.Filter.DueOnly {
		t.Fatal("DueOnly = false, want true")
	}
	if plan.Filter.DueOn != "2026-02-10" {
		t.Fatalf("DueOn = %q, want %q", plan.Filter.DueOn, "2026-02-10")
	}
}

func TestBuildFilterPlanDueDateTomorrow(t *testing.T) {
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
	if plan.Filter == nil {
		t.Fatal("filter request is nil")
	}
	if !plan.Filter.DueOnly {
		t.Fatal("DueOnly = false, want true")
	}
	if plan.Filter.DueOn != "2026-02-11" {
		t.Fatalf("DueOn = %q, want %q", plan.Filter.DueOn, "2026-02-11")
	}
}

func TestBuildFilterPlanDueDateExact(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`find due:2026-02-15`, nlp.ParseOptions{})
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
	if !plan.Filter.DueOnly {
		t.Fatal("DueOnly = false, want true")
	}
	if plan.Filter.DueOn != "2026-02-15" {
		t.Fatalf("DueOn = %q, want %q", plan.Filter.DueOn, "2026-02-15")
	}
}

func TestBuildFilterPlanDueDateCombined(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 10, 10, 0, 0, 0, time.UTC)
	parsed, err := nlp.Parse(`find state:now and due:tomorrow and project:work`, nlp.ParseOptions{Now: now})
	if err != nil {
		t.Fatalf("Parse(filter) error = %v", err)
	}

	plan, err := compile.Build(parsed, compile.BuildOptions{Now: now})
	if err != nil {
		t.Fatalf("Build(filter) error = %v", err)
	}
	if plan.Filter == nil {
		t.Fatal("filter request is nil")
	}
	if len(plan.Filter.States) != 1 || plan.Filter.States[0] != "now" {
		t.Fatalf("states = %v, want [now]", plan.Filter.States)
	}
	if len(plan.Filter.Projects) != 1 || plan.Filter.Projects[0] != "work" {
		t.Fatalf("projects = %v, want [work]", plan.Filter.Projects)
	}
	if !plan.Filter.DueOnly {
		t.Fatal("DueOnly = false, want true")
	}
	if plan.Filter.DueOn != "2026-02-11" {
		t.Fatalf("DueOn = %q, want %q", plan.Filter.DueOn, "2026-02-11")
	}
}

func TestBuildFilterMultiplePredicates(t *testing.T) {
	t.Parallel()

	parsed, err := nlp.Parse(`find state:now and state:done`, nlp.ParseOptions{})
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
	// Multiple states are accumulated
	if len(plan.Filter.States) != 2 || plan.Filter.States[0] != "now" || plan.Filter.States[1] != "done" {
		t.Fatalf("states = %v, want [now done]", plan.Filter.States)
	}
}
