//nolint:testpackage // Tests validate internal filter behavior directly.
package tui

import "testing"

func TestListFiltersToListTasksRequest(t *testing.T) {
	filters := listFilters{
		completion: completionTodo,
		state:      "now",
		project:    "work",
		context:    "office",
		search:     "alpha",
	}

	req := filters.toListTasksRequest()
	if !req.TodoOnly {
		t.Fatal("expected TodoOnly request")
	}
	if req.State != "now" || req.Project != "work" || req.Context != "office" || req.Search != "alpha" {
		t.Fatalf("unexpected request filters: %+v", req)
	}
}

func TestCycleCompletion(t *testing.T) {
	f := listFilters{completion: completionTodo}
	f = f.cycleCompletion()
	if f.completion != completionAll {
		t.Fatalf("expected completionAll, got %v", f.completion)
	}
	f = f.cycleCompletion()
	if f.completion != completionDone {
		t.Fatalf("expected completionDone, got %v", f.completion)
	}
	f = f.cycleCompletion()
	if f.completion != completionTodo {
		t.Fatalf("expected completionTodo, got %v", f.completion)
	}
}

func TestStatusTextIncludesSearch(t *testing.T) {
	f := listFilters{completion: completionAll, search: "pilot"}
	text := f.statusText()
	if text == "" {
		t.Fatal("expected non-empty status text")
	}
	if text != "mode:all search:pilot" {
		t.Fatalf("unexpected status text: %q", text)
	}
}
