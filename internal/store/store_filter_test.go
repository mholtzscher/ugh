//nolint:testpackage // Tests exercise unexported filter compilation helpers directly.
package store

import (
	"context"
	"path/filepath"
	"testing"
	"time"

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

func TestListTasks_UsesSQLForLegacyMultiValueFilters(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := openTestStore(t)

	due := time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC)

	_, err := s.CreateTask(ctx, &Task{Title: "foo bar", State: StateNow, Projects: []string{"work"}, DueOn: &due})
	if err != nil {
		t.Fatalf("CreateTask(1) error = %v", err)
	}
	_, err = s.CreateTask(ctx, &Task{Title: "foo", State: StateWaiting, Projects: []string{"home"}, DueOn: &due})
	if err != nil {
		t.Fatalf("CreateTask(2) error = %v", err)
	}
	_, err = s.CreateTask(ctx, &Task{Title: "foo bar", State: StateDone, Projects: []string{"work"}, DueOn: &due})
	if err != nil {
		t.Fatalf("CreateTask(3) error = %v", err)
	}
	_, err = s.CreateTask(ctx, &Task{Title: "foo bar", State: StateNow, Projects: []string{"misc"}})
	if err != nil {
		t.Fatalf("CreateTask(4) error = %v", err)
	}

	tasks, err := s.ListTasks(ctx, Filters{
		TodoOnly:   true,
		States:     []string{"now", "waiting"},
		Projects:   []string{"work", "home"},
		Search:     []string{"foo", "bar"},
		DueSetOnly: true,
	})
	if err != nil {
		t.Fatalf("ListTasks() error = %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("ListTasks() count = %d, want 1", len(tasks))
	}
	if tasks[0].Title != "foo bar" || tasks[0].State != StateNow {
		t.Fatalf("task = %#v, want now task with title foo bar", tasks[0])
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
