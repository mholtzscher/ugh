//nolint:testpackage // Tests exercise unexported filter compilation helpers directly.
package store

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mholtzscher/ugh/internal/nlp"
)

func TestListTasksByExpr_BooleanSemantics(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := openTestStore(t)

	_, err := s.CreateTask(ctx, &Task{Title: "Now Task", State: StateNow, Projects: []string{"work"}})
	require.NoError(t, err, "CreateTask(now) error")
	_, err = s.CreateTask(ctx, &Task{Title: "Waiting Task", State: StateWaiting, Projects: []string{"home"}})
	require.NoError(t, err, "CreateTask(waiting) error")
	_, err = s.CreateTask(ctx, &Task{Title: "Done Task", State: StateDone, Projects: []string{"work"}})
	require.NoError(t, err, "CreateTask(done) error")

	tasks, err := s.ListTasksByExpr(ctx, nlp.FilterBinary{
		Op:    nlp.FilterOr,
		Left:  nlp.Predicate{Kind: nlp.PredState, Text: "now"},
		Right: nlp.Predicate{Kind: nlp.PredState, Text: "waiting"},
	}, ListTasksByExprOptions{ExcludeDone: true})
	require.NoError(t, err, "ListTasksByExpr(or) error")
	require.Len(t, tasks, 2, "ListTasksByExpr(or) count mismatch")

	tasks, err = s.ListTasksByExpr(ctx, nlp.FilterNot{
		Expr: nlp.Predicate{Kind: nlp.PredState, Text: "done"},
	}, ListTasksByExprOptions{})
	require.NoError(t, err, "ListTasksByExpr(not) error")
	require.Len(t, tasks, 2, "ListTasksByExpr(not) count mismatch")
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
	require.NoError(t, err, "CreateTask(now work) error")

	waitingHomeDue := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	_, err = s.CreateTask(ctx, &Task{
		Title:    "Waiting Home",
		State:    StateWaiting,
		Projects: []string{"home"},
		DueOn:    &waitingHomeDue,
	})
	require.NoError(t, err, "CreateTask(waiting home) error")

	waitingWorkDue := time.Date(2026, time.January, 3, 0, 0, 0, 0, time.UTC)
	waitingWork, err := s.CreateTask(ctx, &Task{
		Title:    "Waiting Work",
		State:    StateWaiting,
		Projects: []string{"work"},
		DueOn:    &waitingWorkDue,
	})
	require.NoError(t, err, "CreateTask(waiting work) error")

	laterWorkDue := time.Date(2026, time.January, 4, 0, 0, 0, 0, time.UTC)
	_, err = s.CreateTask(ctx, &Task{
		Title:    "Later Work",
		State:    StateLater,
		Projects: []string{"work"},
		DueOn:    &laterWorkDue,
	})
	require.NoError(t, err, "CreateTask(later work) error")

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
	require.NoError(t, err, "ListTasksByExpr(nested) error")

	gotIDs := taskIDs(tasks)
	wantIDs := []int64{nowWork.ID, waitingWork.ID}
	assert.Equal(t, wantIDs, gotIDs, "ListTasksByExpr(nested) ids mismatch")
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
	require.NoError(t, err, "Open() error")
	t.Cleanup(func() {
		_ = s.Close()
	})
	return s
}
