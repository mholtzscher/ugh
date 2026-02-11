package shell_test

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mholtzscher/ugh/internal/nlp"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/mholtzscher/ugh/internal/shell"
	"github.com/mholtzscher/ugh/internal/store"
)

func TestExecuteCreateQuotedHashStillInjectsStickyContext(t *testing.T) {
	t.Parallel()

	svc := &recordingService{}
	exec := shell.NewExecutor(svc, &shell.SessionState{ContextProject: "work", ContextContext: "phone"}, true)

	_, err := exec.Execute(context.Background(), `add buy milk "email #hashtag"`)
	require.NoError(t, err, "execute error")

	assert.True(t, contains(svc.lastCreate.Projects, "work"), "projects should contain injected 'work'")
	assert.True(t, contains(svc.lastCreate.Contexts, "phone"), "contexts should contain injected 'phone'")
}

func TestExecuteCreateExplicitProjectSkipsStickyProject(t *testing.T) {
	t.Parallel()

	svc := &recordingService{}
	exec := shell.NewExecutor(svc, &shell.SessionState{ContextProject: "work"}, true)

	_, err := exec.Execute(context.Background(), "add buy milk #personal")
	require.NoError(t, err, "execute error")

	assert.False(t, contains(svc.lastCreate.Projects, "work"), "projects should not contain injected 'work'")
	assert.True(t, contains(svc.lastCreate.Projects, "personal"), "projects should contain explicit 'personal'")
}

func TestExecuteFilterInjectsMissingPredicates(t *testing.T) {
	t.Parallel()

	svc := &recordingService{}
	exec := shell.NewExecutor(svc, &shell.SessionState{ContextProject: "work", ContextContext: "phone"}, true)

	_, err := exec.Execute(context.Background(), "find done")
	require.NoError(t, err, "execute error")

	assert.True(t, hasPredicateKind(svc.lastFilter.Filter, nlp.PredProject), "filter should have project predicate")
	assert.True(t, hasPredicateKind(svc.lastFilter.Filter, nlp.PredContext), "filter should have context predicate")
}

func TestExecuteUpdateInjectsMissingTags(t *testing.T) {
	t.Parallel()

	svc := &recordingService{}
	exec := shell.NewExecutor(svc, &shell.SessionState{ContextProject: "work", ContextContext: "phone"}, true)

	_, err := exec.Execute(context.Background(), "set #7 title: hello")
	require.NoError(t, err, "execute error")

	assert.True(t, contains(svc.lastUpdate.AddProjects, "work"), "add projects should contain injected 'work'")
	assert.True(t, contains(svc.lastUpdate.AddContexts, "phone"), "add contexts should contain injected 'phone'")
}

func TestExecuteFilterStickyProjectWrapsEntireOrExpression(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "test.sqlite")
	st, err := store.Open(ctx, store.Options{Path: dbPath})
	require.NoError(t, err, "open store error")
	t.Cleanup(func() { _ = st.Close() })

	svc := service.NewTaskService(st)

	_, err = svc.CreateTask(ctx, service.CreateTaskRequest{Title: "now-home", State: "now", Projects: []string{"home"}})
	require.NoError(t, err, "create now-home error")
	waitingWork, err := svc.CreateTask(ctx, service.CreateTaskRequest{
		Title:    "waiting-work",
		State:    "waiting",
		Projects: []string{"work"},
	})
	require.NoError(t, err, "create waiting-work error")

	exec := shell.NewExecutor(svc, &shell.SessionState{ContextProject: "work"}, true)
	result, err := exec.Execute(ctx, "find state:now or state:waiting")
	require.NoError(t, err, "execute error")

	require.Len(t, result.TaskIDs, 1, "task IDs should have exactly one task")
	require.Equal(t, waitingWork.ID, result.TaskIDs[0], "task ID mismatch")
}

func TestExecuteUpdateStickyProjectStillAllowsExplicitRemoval(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "test.sqlite")
	st, err := store.Open(ctx, store.Options{Path: dbPath})
	require.NoError(t, err, "open store error")
	t.Cleanup(func() { _ = st.Close() })

	svc := service.NewTaskService(st)
	task, err := svc.CreateTask(ctx, service.CreateTaskRequest{Title: "cleanup", Projects: []string{"work"}})
	require.NoError(t, err, "create task error")

	exec := shell.NewExecutor(svc, &shell.SessionState{ContextProject: "work"}, true)
	_, err = exec.Execute(ctx, fmt.Sprintf("set #%d -project:work", task.ID))
	require.NoError(t, err, "execute error")

	updated, err := svc.GetTask(ctx, task.ID)
	require.NoError(t, err, "get updated task error")
	assert.False(t, contains(updated.Projects, "work"), "projects should not contain 'work' after explicit removal")
}

type recordingService struct {
	lastCreate service.CreateTaskRequest
	lastUpdate service.UpdateTaskRequest
	lastFilter service.ListTasksRequest
}

func (s *recordingService) CreateTask(_ context.Context, req service.CreateTaskRequest) (*store.Task, error) {
	s.lastCreate = req
	return &store.Task{
		ID:       1,
		Title:    req.Title,
		Projects: req.Projects,
		Contexts: req.Contexts,
		State:    store.StateInbox,
	}, nil
}

func (s *recordingService) ListTasks(_ context.Context, req service.ListTasksRequest) ([]*store.Task, error) {
	s.lastFilter = req
	return []*store.Task{}, nil
}

func (s *recordingService) GetTask(_ context.Context, _ int64) (*store.Task, error) {
	return &store.Task{}, nil
}

func (s *recordingService) UpdateTask(_ context.Context, req service.UpdateTaskRequest) (*store.Task, error) {
	s.lastUpdate = req
	return &store.Task{ID: req.ID, Title: "updated", State: store.StateInbox}, nil
}

func (*recordingService) FullUpdateTask(_ context.Context, _ service.FullUpdateTaskRequest) (*store.Task, error) {
	return &store.Task{}, nil
}

func (*recordingService) SetDone(_ context.Context, _ []int64, _ bool) (int64, error) {
	return 0, nil
}

func (*recordingService) DeleteTasks(_ context.Context, _ []int64) (int64, error) {
	return 0, nil
}

func (*recordingService) ListProjects(_ context.Context, _ service.ListTagsRequest) ([]store.NameCount, error) {
	return []store.NameCount{}, nil
}

func (*recordingService) ListContexts(_ context.Context, _ service.ListTagsRequest) ([]store.NameCount, error) {
	return []store.NameCount{}, nil
}

func (*recordingService) Sync(_ context.Context) error {
	return nil
}

func (*recordingService) Push(_ context.Context) error {
	return nil
}

func (*recordingService) SyncStatus(_ context.Context) (*service.SyncStatus, error) {
	return &service.SyncStatus{}, nil
}

func (*recordingService) Close() error {
	return nil
}

func (*recordingService) RecordShellHistory(
	_ context.Context, _ string, _ bool, _ string, _ string,
) (*store.ShellHistory, error) {
	return &store.ShellHistory{}, nil
}

func (*recordingService) ListShellHistory(_ context.Context, _ int64) ([]*store.ShellHistory, error) {
	return []*store.ShellHistory{}, nil
}

func (*recordingService) SearchShellHistory(
	_ context.Context, _, _ string, _ *bool, _ int64,
) ([]*store.ShellHistory, error) {
	return []*store.ShellHistory{}, nil
}

func (*recordingService) ClearShellHistory(_ context.Context) error {
	return nil
}

func contains(values []string, want string) bool {
	return slices.Contains(values, want)
}

func hasPredicateKind(expr nlp.FilterExpr, kind nlp.PredicateKind) bool {
	switch typed := expr.(type) {
	case nlp.Predicate:
		return typed.Kind == kind
	case nlp.FilterBinary:
		return hasPredicateKind(typed.Left, kind) || hasPredicateKind(typed.Right, kind)
	case nlp.FilterNot:
		return hasPredicateKind(typed.Expr, kind)
	default:
		return false
	}
}

func TestExecuteEmptyCommand(t *testing.T) {
	t.Parallel()

	svc := &recordingService{}
	exec := shell.NewExecutor(svc, &shell.SessionState{}, true)

	_, err := exec.Execute(context.Background(), "")
	require.Error(t, err, "expected error for empty command")
}

func TestExecuteWhitespaceOnlyCommand(t *testing.T) {
	t.Parallel()

	svc := &recordingService{}
	exec := shell.NewExecutor(svc, &shell.SessionState{}, true)

	_, err := exec.Execute(context.Background(), "   \t\n  ")
	require.Error(t, err, "expected error for whitespace-only command")
}

func TestExecuteRejectsControlCharacters(t *testing.T) {
	t.Parallel()

	svc := &recordingService{}
	exec := shell.NewExecutor(svc, &shell.SessionState{}, true)

	_, err := exec.Execute(context.Background(), "\x16set #1 projects:work")
	require.Error(t, err, "expected error for control character input")
	assert.Contains(t, err.Error(), "control character", "error should explain invalid control character")
	assert.Contains(t, err.Error(), "U+0016", "error should include Unicode control code")
	assert.Contains(t, err.Error(), `\x16`, "error should include escaped control character")
	assert.Contains(t, err.Error(), "rune position 1", "error should include 1-based rune position")
}

func TestExecuteRejectsControlCharactersReportsRunePosition(t *testing.T) {
	t.Parallel()

	svc := &recordingService{}
	exec := shell.NewExecutor(svc, &shell.SessionState{}, true)

	_, err := exec.Execute(context.Background(), "é\x16set #1 projects:work")
	require.Error(t, err, "expected error for control character input")
	assert.Contains(t, err.Error(), "U+0016", "error should include Unicode control code")
	assert.Contains(t, err.Error(), "rune position 2", "error should report rune-based position")
}

func TestExecuteUpdatesLastTaskID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "test.sqlite")
	st, err := store.Open(ctx, store.Options{Path: dbPath})
	require.NoError(t, err, "open store error")
	t.Cleanup(func() { _ = st.Close() })

	svc := service.NewTaskService(st)
	state := &shell.SessionState{}
	exec := shell.NewExecutor(svc, state, true)

	// Create first task
	result1, err := exec.Execute(ctx, "add first task")
	require.NoError(t, err, "execute first error")

	require.Len(t, state.LastTaskIDs, 1, "LastTaskIDs should have 1 entry")
	assert.Equal(t, result1.TaskIDs[0], state.LastTaskIDs[0], "LastTaskIDs[0] mismatch")

	// Create second task
	result2, err := exec.Execute(ctx, "add second task")
	require.NoError(t, err, "execute second error")

	require.Len(t, state.LastTaskIDs, 1, "LastTaskIDs should have 1 entry after second create")
	assert.Equal(t, result2.TaskIDs[0], state.LastTaskIDs[0], "LastTaskIDs[0] mismatch")
}

func TestExecuteFilterUpdatesLastTaskIDs(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "test.sqlite")
	st, err := store.Open(ctx, store.Options{Path: dbPath})
	require.NoError(t, err, "open store error")
	t.Cleanup(func() { _ = st.Close() })

	svc := service.NewTaskService(st)

	// Create multiple tasks
	var taskIDs []int64
	for i := 1; i <= 3; i++ {
		task, createErr := svc.CreateTask(ctx, service.CreateTaskRequest{
			Title: fmt.Sprintf("task %d", i),
			State: "now",
		})
		require.NoError(t, createErr, "create task %d error", i)
		taskIDs = append(taskIDs, task.ID)
	}

	state := &shell.SessionState{}
	exec := shell.NewExecutor(svc, state, true)

	_, err = exec.Execute(ctx, "find state:now")
	require.NoError(t, err, "execute filter error")

	require.Len(t, state.LastTaskIDs, 3, "LastTaskIDs should have 3 entries")

	// Check that all task IDs are present
	for _, id := range taskIDs {
		assert.True(t, slices.Contains(state.LastTaskIDs, id), "task ID %d not found in LastTaskIDs", id)
	}
}

func TestExecuteUpdateSetsSelectedTaskID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "test.sqlite")
	st, err := store.Open(ctx, store.Options{Path: dbPath})
	require.NoError(t, err, "open store error")
	t.Cleanup(func() { _ = st.Close() })

	svc := service.NewTaskService(st)
	task, err := svc.CreateTask(ctx, service.CreateTaskRequest{Title: "task to update"})
	require.NoError(t, err, "create task error")

	// Set SelectedTaskID first (required to use "selected" target)
	state := &shell.SessionState{
		SelectedTaskID: &task.ID,
	}
	exec := shell.NewExecutor(svc, state, true)

	// Update with "selected" target should keep SelectedTaskID set
	_, err = exec.Execute(ctx, "set selected title:updated")
	require.NoError(t, err, "execute update error")

	require.NotNil(t, state.SelectedTaskID, "SelectedTaskID should be set")
	assert.Equal(t, task.ID, *state.SelectedTaskID, "SelectedTaskID mismatch")
}

func TestExecuteUpdateByIDDoesNotSetSelectedTaskID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "test.sqlite")
	st, err := store.Open(ctx, store.Options{Path: dbPath})
	require.NoError(t, err, "open store error")
	t.Cleanup(func() { _ = st.Close() })

	svc := service.NewTaskService(st)
	task, err := svc.CreateTask(ctx, service.CreateTaskRequest{Title: "task to update"})
	require.NoError(t, err, "create task error")

	state := &shell.SessionState{}
	exec := shell.NewExecutor(svc, state, true)

	// Update by explicit ID should NOT set SelectedTaskID
	_, err = exec.Execute(ctx, fmt.Sprintf("set #%d title:updated", task.ID))
	require.NoError(t, err, "execute update error")

	assert.Nil(t, state.SelectedTaskID, "SelectedTaskID should be nil when updating by ID")
}

func TestExecutePronounSubstitution(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "test.sqlite")
	st, err := store.Open(ctx, store.Options{Path: dbPath})
	require.NoError(t, err, "open store error")
	t.Cleanup(func() { _ = st.Close() })

	svc := service.NewTaskService(st)

	// Create two tasks
	task1, err := svc.CreateTask(ctx, service.CreateTaskRequest{Title: "first task"})
	require.NoError(t, err, "create first task error")
	task2, err := svc.CreateTask(ctx, service.CreateTaskRequest{Title: "second task"})
	require.NoError(t, err, "create second task error")

	state := &shell.SessionState{
		LastTaskIDs: []int64{task1.ID, task2.ID},
	}
	exec := shell.NewExecutor(svc, state, true)

	// "it" and "this" should refer to last task (task2)
	_, err = exec.Execute(ctx, "set it title:updated via it")
	require.NoError(t, err, "execute with 'it' error")

	updated, err := svc.GetTask(ctx, task2.ID)
	require.NoError(t, err, "get updated task error")
	assert.Equal(t, "updated via it", updated.Title, "task2 title mismatch")

	// "last" should also refer to task2
	_, err = exec.Execute(ctx, "set last title:updated via last")
	require.NoError(t, err, "execute with 'last' error")

	updated, err = svc.GetTask(ctx, task2.ID)
	require.NoError(t, err, "get updated task error")
	assert.Equal(t, "updated via last", updated.Title, "task2 title mismatch")
}

func TestExecuteThatPronounSubstitution(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "test.sqlite")
	st, err := store.Open(ctx, store.Options{Path: dbPath})
	require.NoError(t, err, "open store error")
	t.Cleanup(func() { _ = st.Close() })

	svc := service.NewTaskService(st)

	// Create two tasks
	task1, err := svc.CreateTask(ctx, service.CreateTaskRequest{Title: "first task"})
	require.NoError(t, err, "create first task error")
	task2, err := svc.CreateTask(ctx, service.CreateTaskRequest{Title: "second task"})
	require.NoError(t, err, "create second task error")

	state := &shell.SessionState{
		LastTaskIDs: []int64{task1.ID, task2.ID},
	}
	exec := shell.NewExecutor(svc, state, true)

	// "that" should refer to second-to-last task (task1)
	_, err = exec.Execute(ctx, "set that title:updated via that")
	require.NoError(t, err, "execute with 'that' error")

	updated, err := svc.GetTask(ctx, task1.ID)
	require.NoError(t, err, "get updated task error")
	assert.Equal(t, "updated via that", updated.Title, "task1 title mismatch")
}

func TestExecuteSelectedPronounSubstitution(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "test.sqlite")
	st, err := store.Open(ctx, store.Options{Path: dbPath})
	require.NoError(t, err, "open store error")
	t.Cleanup(func() { _ = st.Close() })

	svc := service.NewTaskService(st)
	task, err := svc.CreateTask(ctx, service.CreateTaskRequest{Title: "task"})
	require.NoError(t, err, "create task error")

	state := &shell.SessionState{
		SelectedTaskID: &task.ID,
	}
	exec := shell.NewExecutor(svc, state, true)

	// "selected" should refer to selected task
	_, err = exec.Execute(ctx, "set selected title:updated via selected")
	require.NoError(t, err, "execute with 'selected' error")

	updated, err := svc.GetTask(ctx, task.ID)
	require.NoError(t, err, "get updated task error")
	assert.Equal(t, "updated via selected", updated.Title, "task title mismatch")
}

func TestExecuteCreateWithoutContext(t *testing.T) {
	t.Parallel()

	svc := &recordingService{}
	exec := shell.NewExecutor(svc, &shell.SessionState{}, true)

	_, err := exec.Execute(context.Background(), "add buy milk")
	require.NoError(t, err, "execute error")

	assert.Empty(t, svc.lastCreate.Projects, "projects should be empty")
	assert.Empty(t, svc.lastCreate.Contexts, "contexts should be empty")
}

func TestExecuteFilterWithoutContext(t *testing.T) {
	t.Parallel()

	svc := &recordingService{}
	exec := shell.NewExecutor(svc, &shell.SessionState{}, true)

	_, err := exec.Execute(context.Background(), "find state:now")
	require.NoError(t, err, "execute error")

	// Without context, should not inject project/context predicates
	hasProject := hasPredicateKind(svc.lastFilter.Filter, nlp.PredProject)
	hasContext := hasPredicateKind(svc.lastFilter.Filter, nlp.PredContext)

	assert.False(t, hasProject, "filter should not have project predicate without context")
	assert.False(t, hasContext, "filter should not have context predicate without context")
}

func TestExecuteInjectContextOnlyProject(t *testing.T) {
	t.Parallel()

	svc := &recordingService{}
	exec := shell.NewExecutor(svc, &shell.SessionState{ContextProject: "work"}, true)

	_, err := exec.Execute(context.Background(), "add buy milk")
	require.NoError(t, err, "execute error")

	assert.True(t, contains(svc.lastCreate.Projects, "work"), "projects should contain 'work'")
	assert.Empty(t, svc.lastCreate.Contexts, "contexts should be empty when only project context set")
}

func TestExecuteInjectContextOnlyContext(t *testing.T) {
	t.Parallel()

	svc := &recordingService{}
	exec := shell.NewExecutor(svc, &shell.SessionState{ContextContext: "phone"}, true)

	_, err := exec.Execute(context.Background(), "add buy milk")
	require.NoError(t, err, "execute error")

	assert.Empty(t, svc.lastCreate.Projects, "projects should be empty when only context context set")
	assert.True(t, contains(svc.lastCreate.Contexts, "phone"), "contexts should contain 'phone'")
}

func TestExecuteViewShowsHelp(t *testing.T) {
	t.Parallel()

	svc := &recordingService{}
	exec := shell.NewExecutor(svc, &shell.SessionState{}, true)

	result, err := exec.Execute(context.Background(), "view")
	require.NoError(t, err, "execute error")
	require.NotNil(t, result, "result should not be nil")
	assert.Equal(t, "view", result.Intent, "intent mismatch")
	assert.Contains(t, result.Message, "Available Views:", "message should contain view help")
}

func TestExecuteViewRunsFilterQuery(t *testing.T) {
	t.Parallel()

	svc := &recordingService{}
	exec := shell.NewExecutor(svc, &shell.SessionState{}, true)

	result, err := exec.Execute(context.Background(), "view now")
	require.NoError(t, err, "execute error")
	require.NotNil(t, result, "result should not be nil")
	assert.Equal(t, "filter", result.Intent, "view now should execute as filter")
	assert.True(t, hasPredicateKind(svc.lastFilter.Filter, nlp.PredState), "filter should include state predicate")
}

func TestExecuteViewRespectsStickyContext(t *testing.T) {
	t.Parallel()

	svc := &recordingService{}
	exec := shell.NewExecutor(svc, &shell.SessionState{ContextProject: "work"}, true)

	_, err := exec.Execute(context.Background(), "view now")
	require.NoError(t, err, "execute error")
	assert.True(
		t,
		hasPredicateKind(svc.lastFilter.Filter, nlp.PredProject),
		"view filter should include sticky project",
	)
}

func TestExecuteContextSetShowAndClear(t *testing.T) {
	t.Parallel()

	selectedID := int64(9)
	state := &shell.SessionState{
		SelectedTaskID: &selectedID,
		LastTaskIDs:    []int64{3, 7},
	}
	exec := shell.NewExecutor(&recordingService{}, state, true)

	result, err := exec.Execute(context.Background(), "context #work")
	require.NoError(t, err, "set project context error")
	assert.Equal(t, "context", result.Intent, "intent mismatch")
	assert.Equal(t, "work", state.ContextProject, "project context mismatch")

	result, err = exec.Execute(context.Background(), "context @urgent")
	require.NoError(t, err, "set context filter error")
	assert.Equal(t, "context", result.Intent, "intent mismatch")
	assert.Equal(t, "urgent", state.ContextContext, "context filter mismatch")

	result, err = exec.Execute(context.Background(), "context")
	require.NoError(t, err, "show context error")
	assert.Contains(t, result.Message, "Selected: #9", "context output should include selected task")
	assert.Contains(t, result.Message, "Last: #3, #7", "context output should include last tasks")
	assert.Contains(t, result.Message, "Project: #work", "context output should include project")
	assert.Contains(t, result.Message, "Context: @urgent", "context output should include context")

	result, err = exec.Execute(context.Background(), "context clear")
	require.NoError(t, err, "clear context error")
	assert.Equal(t, "context", result.Intent, "intent mismatch")
	assert.Empty(t, state.ContextProject, "project context should be cleared")
	assert.Empty(t, state.ContextContext, "context filter should be cleared")
}
