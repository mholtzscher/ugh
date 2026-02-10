//nolint:testpackage // Tests exercise unexported filtering and ordering helpers directly.
package store

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/mholtzscher/ugh/internal/nlp"
)

func TestFilterBySearchTermsMatchesTaskFields(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 9, 12, 0, 0, 0, time.UTC)
	tasks := []*Task{
		{
			ID:        1,
			Title:     "Alpha",
			Notes:     "Foo note",
			Projects:  []string{"work"},
			Contexts:  []string{"office"},
			Meta:      map[string]string{"priority": "urgent"},
			UpdatedAt: now,
		},
		{
			ID:        2,
			Title:     "Beta",
			Notes:     "Bar note",
			Projects:  []string{"home"},
			Contexts:  []string{"garden"},
			Meta:      map[string]string{"tag": "optional"},
			UpdatedAt: now,
		},
	}

	filtered := filterBySearchTerms(tasks, []string{"foo", "urgent"})
	if len(filtered) != 1 {
		t.Fatalf("filtered length = %d, want 1", len(filtered))
	}
	if filtered[0].ID != 1 {
		t.Fatalf("filtered[0].ID = %d, want 1", filtered[0].ID)
	}
}

func TestSortTasksForListMatchesQueryOrdering(t *testing.T) {
	t.Parallel()

	dueSoon := time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC)
	dueLater := time.Date(2026, 2, 12, 0, 0, 0, 0, time.UTC)
	newer := time.Date(2026, 2, 9, 15, 0, 0, 0, time.UTC)
	older := time.Date(2026, 2, 9, 12, 0, 0, 0, time.UTC)

	tasks := []*Task{
		{ID: 4, State: StateDone, DueOn: &dueSoon, UpdatedAt: newer},
		{ID: 3, State: StateNow, DueOn: nil, UpdatedAt: newer},
		{ID: 2, State: StateNow, DueOn: &dueLater, UpdatedAt: older},
		{ID: 1, State: StateNow, DueOn: &dueSoon, UpdatedAt: newer},
	}

	sortTasksForList(tasks)

	got := []int64{tasks[0].ID, tasks[1].ID, tasks[2].ID, tasks[3].ID}
	want := []int64{1, 2, 3, 4}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("order = %v, want %v", got, want)
		}
	}
}

func TestListTasksByExpr_BooleanSemantics(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := openTestStore(t)

	_, err := s.CreateTask(ctx, &Task{Title: "Now Task", State: StateNow, Projects: []string{"work"}})
	if err != nil {
		t.Fatalf("CreateTask(now) error = %v", err)
	}
	_, err = s.CreateTask(ctx, &Task{Title: "Waiting Task", State: StateWaiting, Projects: []string{"home"}})
	if err != nil {
		t.Fatalf("CreateTask(waiting) error = %v", err)
	}
	_, err = s.CreateTask(ctx, &Task{Title: "Done Task", State: StateDone, Projects: []string{"work"}})
	if err != nil {
		t.Fatalf("CreateTask(done) error = %v", err)
	}

	tasks, err := s.ListTasksByExpr(ctx, nlp.FilterBinary{
		Op:    nlp.FilterOr,
		Left:  nlp.Predicate{Kind: nlp.PredState, Text: "now"},
		Right: nlp.Predicate{Kind: nlp.PredState, Text: "waiting"},
	}, ListTasksByExprOptions{ExcludeDone: true})
	if err != nil {
		t.Fatalf("ListTasksByExpr(or) error = %v", err)
	}
	if len(tasks) != 2 {
		t.Fatalf("ListTasksByExpr(or) count = %d, want 2", len(tasks))
	}

	tasks, err = s.ListTasksByExpr(ctx, nlp.FilterNot{
		Expr: nlp.Predicate{Kind: nlp.PredState, Text: "done"},
	}, ListTasksByExprOptions{})
	if err != nil {
		t.Fatalf("ListTasksByExpr(not) error = %v", err)
	}
	if len(tasks) != 2 {
		t.Fatalf("ListTasksByExpr(not) count = %d, want 2", len(tasks))
	}
}

func openTestStore(t *testing.T) *Store {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "test.sqlite")
	s, err := Open(context.Background(), Options{Path: dbPath})
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	t.Cleanup(func() {
		_ = s.Close()
	})
	return s
}
