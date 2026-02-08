//nolint:testpackage // Tests validate internal filter behavior directly.
package tui

import "testing"

func TestListFiltersToListTasksRequest(t *testing.T) {
	tests := []struct {
		name     string
		state    string
		wantAll  bool
		wantDone bool
		wantTodo bool
	}{
		{name: "state tab uses todo", state: "now", wantTodo: true},
		{name: "done tab uses done only", state: "done", wantDone: true},
		{name: "all tab uses all", state: "", wantAll: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filters := listFilters{
				state:   tt.state,
				project: "work",
				context: "office",
				search:  "alpha",
			}

			req := filters.toListTasksRequest()
			if req.All != tt.wantAll || req.DoneOnly != tt.wantDone || req.TodoOnly != tt.wantTodo {
				t.Fatalf(
					"unexpected completion flags: all=%v doneOnly=%v todoOnly=%v",
					req.All,
					req.DoneOnly,
					req.TodoOnly,
				)
			}
			if req.State != tt.state || req.Project != "work" || req.Context != "office" || req.Search != "alpha" {
				t.Fatalf("unexpected request filters: %+v", req)
			}
		})
	}
}

func TestListFiltersToListTagsRequest(t *testing.T) {
	tests := []struct {
		name     string
		state    string
		wantAll  bool
		wantDone bool
		wantTodo bool
	}{
		{name: "state tab uses todo", state: "now", wantTodo: true},
		{name: "done tab uses done only", state: "done", wantDone: true},
		{name: "all tab uses all", state: "", wantAll: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := (listFilters{state: tt.state}).toListTagsRequest()
			if req.All != tt.wantAll || req.DoneOnly != tt.wantDone || req.TodoOnly != tt.wantTodo {
				t.Fatalf(
					"unexpected completion flags: all=%v doneOnly=%v todoOnly=%v",
					req.All,
					req.DoneOnly,
					req.TodoOnly,
				)
			}
		})
	}
}

func TestStatusTextIncludesSearch(t *testing.T) {
	f := listFilters{search: "pilot"}
	text := f.statusText()
	if text == "" {
		t.Fatal("expected non-empty status text")
	}
	if text != "search:pilot" {
		t.Fatalf("unexpected status text: %q", text)
	}
}

func TestStatusTextIncludesDueOnly(t *testing.T) {
	f := listFilters{dueOnly: true}
	text := f.statusText()
	if text != "due:only" {
		t.Fatalf("unexpected status text: %q", text)
	}
}
