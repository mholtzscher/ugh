//nolint:testpackage // Tests exercise unexported filter compilation helpers directly.
package store

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/mholtzscher/ugh/internal/nlp"
)

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
