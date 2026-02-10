//nolint:testpackage // Tests exercise unexported filter compilation helpers directly.
package store

import (
	"context"
	"path/filepath"
	"reflect"
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

func TestListTasksByExpr_NestedBooleanFilterIDs(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := openTestStore(t)

	nowWorkDue := time.Date(2026, time.January, 2, 0, 0, 0, 0, time.UTC)
	nowWork, err := s.CreateTask(ctx, &Task{
		Title:    "Now Work",
		State:    StateNow,
		Projects: []string{"work"},
		DueOn:    &nowWorkDue,
	})
	if err != nil {
		t.Fatalf("CreateTask(now work) error = %v", err)
	}

	waitingHomeDue := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	_, err = s.CreateTask(ctx, &Task{
		Title:    "Waiting Home",
		State:    StateWaiting,
		Projects: []string{"home"},
		DueOn:    &waitingHomeDue,
	})
	if err != nil {
		t.Fatalf("CreateTask(waiting home) error = %v", err)
	}

	waitingWorkDue := time.Date(2026, time.January, 3, 0, 0, 0, 0, time.UTC)
	waitingWork, err := s.CreateTask(ctx, &Task{
		Title:    "Waiting Work",
		State:    StateWaiting,
		Projects: []string{"work"},
		DueOn:    &waitingWorkDue,
	})
	if err != nil {
		t.Fatalf("CreateTask(waiting work) error = %v", err)
	}

	laterWorkDue := time.Date(2026, time.January, 4, 0, 0, 0, 0, time.UTC)
	_, err = s.CreateTask(ctx, &Task{
		Title:    "Later Work",
		State:    StateLater,
		Projects: []string{"work"},
		DueOn:    &laterWorkDue,
	})
	if err != nil {
		t.Fatalf("CreateTask(later work) error = %v", err)
	}

	filter := nlp.FilterBinary{
		Op: nlp.FilterAnd,
		Left: nlp.FilterBinary{
			Op:    nlp.FilterOr,
			Left:  nlp.Predicate{Kind: nlp.PredState, Text: "now"},
			Right: nlp.Predicate{Kind: nlp.PredState, Text: "waiting"},
		},
		Right: nlp.FilterNot{
			Expr: nlp.Predicate{Kind: nlp.PredProject, Text: "home"},
		},
	}

	tasks, err := s.ListTasksByExpr(ctx, filter, ListTasksByExprOptions{ExcludeDone: true})
	if err != nil {
		t.Fatalf("ListTasksByExpr(nested) error = %v", err)
	}

	gotIDs := taskIDs(tasks)
	wantIDs := []int64{nowWork.ID, waitingWork.ID}
	if !reflect.DeepEqual(gotIDs, wantIDs) {
		t.Fatalf("ListTasksByExpr(nested) ids = %v, want %v", gotIDs, wantIDs)
	}
}

func taskIDs(tasks []*Task) []int64 {
	ids := make([]int64, 0, len(tasks))
	for _, task := range tasks {
		ids = append(ids, task.ID)
	}
	return ids
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
