package shell_test

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"testing"

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
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if !contains(svc.lastCreate.Projects, "work") {
		t.Fatalf("projects = %#v, want injected 'work'", svc.lastCreate.Projects)
	}
	if !contains(svc.lastCreate.Contexts, "phone") {
		t.Fatalf("contexts = %#v, want injected 'phone'", svc.lastCreate.Contexts)
	}
}

func TestExecuteCreateExplicitProjectSkipsStickyProject(t *testing.T) {
	t.Parallel()

	svc := &recordingService{}
	exec := shell.NewExecutor(svc, &shell.SessionState{ContextProject: "work"}, true)

	_, err := exec.Execute(context.Background(), "add buy milk #personal")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if contains(svc.lastCreate.Projects, "work") {
		t.Fatalf("projects = %#v, did not expect injected 'work'", svc.lastCreate.Projects)
	}
	if !contains(svc.lastCreate.Projects, "personal") {
		t.Fatalf("projects = %#v, want explicit 'personal'", svc.lastCreate.Projects)
	}
}

func TestExecuteFilterInjectsMissingPredicates(t *testing.T) {
	t.Parallel()

	svc := &recordingService{}
	exec := shell.NewExecutor(svc, &shell.SessionState{ContextProject: "work", ContextContext: "phone"}, true)

	_, err := exec.Execute(context.Background(), "find done")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if !hasPredicateKind(svc.lastFilter.Filter, nlp.PredProject) {
		t.Fatalf("filter = %#v, want project predicate", svc.lastFilter.Filter)
	}
	if !hasPredicateKind(svc.lastFilter.Filter, nlp.PredContext) {
		t.Fatalf("filter = %#v, want context predicate", svc.lastFilter.Filter)
	}
}

func TestExecuteUpdateInjectsMissingTags(t *testing.T) {
	t.Parallel()

	svc := &recordingService{}
	exec := shell.NewExecutor(svc, &shell.SessionState{ContextProject: "work", ContextContext: "phone"}, true)

	_, err := exec.Execute(context.Background(), "set #7 title: hello")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if !contains(svc.lastUpdate.AddProjects, "work") {
		t.Fatalf("add projects = %#v, want injected 'work'", svc.lastUpdate.AddProjects)
	}
	if !contains(svc.lastUpdate.AddContexts, "phone") {
		t.Fatalf("add contexts = %#v, want injected 'phone'", svc.lastUpdate.AddContexts)
	}
}

func TestExecuteFilterStickyProjectWrapsEntireOrExpression(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "test.sqlite")
	st, err := store.Open(ctx, store.Options{Path: dbPath})
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })

	svc := service.NewTaskService(st)

	_, err = svc.CreateTask(ctx, service.CreateTaskRequest{Title: "now-home", State: "now", Projects: []string{"home"}})
	if err != nil {
		t.Fatalf("create now-home: %v", err)
	}
	waitingWork, err := svc.CreateTask(ctx, service.CreateTaskRequest{
		Title:    "waiting-work",
		State:    "waiting",
		Projects: []string{"work"},
	})
	if err != nil {
		t.Fatalf("create waiting-work: %v", err)
	}

	exec := shell.NewExecutor(svc, &shell.SessionState{ContextProject: "work"}, true)
	result, err := exec.Execute(ctx, "find state:now or state:waiting")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if len(result.TaskIDs) != 1 {
		t.Fatalf("task IDs = %#v, want exactly one task", result.TaskIDs)
	}
	if result.TaskIDs[0] != waitingWork.ID {
		t.Fatalf("task IDs = %#v, want [%d]", result.TaskIDs, waitingWork.ID)
	}
}

func TestExecuteUpdateStickyProjectStillAllowsExplicitRemoval(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "test.sqlite")
	st, err := store.Open(ctx, store.Options{Path: dbPath})
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })

	svc := service.NewTaskService(st)
	task, err := svc.CreateTask(ctx, service.CreateTaskRequest{Title: "cleanup", Projects: []string{"work"}})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	exec := shell.NewExecutor(svc, &shell.SessionState{ContextProject: "work"}, true)
	_, err = exec.Execute(ctx, fmt.Sprintf("set #%d -project:work", task.ID))
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	updated, err := svc.GetTask(ctx, task.ID)
	if err != nil {
		t.Fatalf("get updated task: %v", err)
	}
	if contains(updated.Projects, "work") {
		t.Fatalf("projects = %#v, expected explicit removal to win", updated.Projects)
	}
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
	if err == nil {
		t.Fatal("expected error for empty command, got nil")
	}
}

func TestExecuteWhitespaceOnlyCommand(t *testing.T) {
	t.Parallel()

	svc := &recordingService{}
	exec := shell.NewExecutor(svc, &shell.SessionState{}, true)

	_, err := exec.Execute(context.Background(), "   \t\n  ")
	if err == nil {
		t.Fatal("expected error for whitespace-only command, got nil")
	}
}

func TestExecuteUpdatesLastTaskID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "test.sqlite")
	st, err := store.Open(ctx, store.Options{Path: dbPath})
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })

	svc := service.NewTaskService(st)
	state := &shell.SessionState{}
	exec := shell.NewExecutor(svc, state, true)

	// Create first task
	result1, err := exec.Execute(ctx, "add first task")
	if err != nil {
		t.Fatalf("execute first: %v", err)
	}

	if len(state.LastTaskIDs) != 1 {
		t.Fatalf("LastTaskIDs = %v, want 1 entry", state.LastTaskIDs)
	}
	if state.LastTaskIDs[0] != result1.TaskIDs[0] {
		t.Errorf("LastTaskIDs[0] = %d, want %d", state.LastTaskIDs[0], result1.TaskIDs[0])
	}

	// Create second task
	result2, err := exec.Execute(ctx, "add second task")
	if err != nil {
		t.Fatalf("execute second: %v", err)
	}

	if len(state.LastTaskIDs) != 1 {
		t.Fatalf("LastTaskIDs = %v, want 1 entry after second create", state.LastTaskIDs)
	}
	if state.LastTaskIDs[0] != result2.TaskIDs[0] {
		t.Errorf("LastTaskIDs[0] = %d, want %d", state.LastTaskIDs[0], result2.TaskIDs[0])
	}
}

func TestExecuteFilterUpdatesLastTaskIDs(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "test.sqlite")
	st, err := store.Open(ctx, store.Options{Path: dbPath})
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })

	svc := service.NewTaskService(st)

	// Create multiple tasks
	var taskIDs []int64
	for i := 1; i <= 3; i++ {
		task, createErr := svc.CreateTask(ctx, service.CreateTaskRequest{
			Title: fmt.Sprintf("task %d", i),
			State: "now",
		})
		if createErr != nil {
			t.Fatalf("create task %d: %v", i, createErr)
		}
		taskIDs = append(taskIDs, task.ID)
	}

	state := &shell.SessionState{}
	exec := shell.NewExecutor(svc, state, true)

	_, err = exec.Execute(ctx, "find state:now")
	if err != nil {
		t.Fatalf("execute filter: %v", err)
	}

	if len(state.LastTaskIDs) != 3 {
		t.Fatalf("LastTaskIDs = %v, want 3 entries", state.LastTaskIDs)
	}

	// Check that all task IDs are present
	for _, id := range taskIDs {
		if !slices.Contains(state.LastTaskIDs, id) {
			t.Errorf("task ID %d not found in LastTaskIDs", id)
		}
	}
}

func TestExecuteUpdateSetsSelectedTaskID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "test.sqlite")
	st, err := store.Open(ctx, store.Options{Path: dbPath})
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })

	svc := service.NewTaskService(st)
	task, err := svc.CreateTask(ctx, service.CreateTaskRequest{Title: "task to update"})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	// Set SelectedTaskID first (required to use "selected" target)
	state := &shell.SessionState{
		SelectedTaskID: &task.ID,
	}
	exec := shell.NewExecutor(svc, state, true)

	// Update with "selected" target should keep SelectedTaskID set
	_, err = exec.Execute(ctx, "set selected title:updated")
	if err != nil {
		t.Fatalf("execute update: %v", err)
	}

	if state.SelectedTaskID == nil {
		t.Fatal("SelectedTaskID is nil, expected it to be set")
	}
	if *state.SelectedTaskID != task.ID {
		t.Errorf("SelectedTaskID = %d, want %d", *state.SelectedTaskID, task.ID)
	}
}

func TestExecuteUpdateByIDDoesNotSetSelectedTaskID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "test.sqlite")
	st, err := store.Open(ctx, store.Options{Path: dbPath})
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })

	svc := service.NewTaskService(st)
	task, err := svc.CreateTask(ctx, service.CreateTaskRequest{Title: "task to update"})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	state := &shell.SessionState{}
	exec := shell.NewExecutor(svc, state, true)

	// Update by explicit ID should NOT set SelectedTaskID
	_, err = exec.Execute(ctx, fmt.Sprintf("set #%d title:updated", task.ID))
	if err != nil {
		t.Fatalf("execute update: %v", err)
	}

	if state.SelectedTaskID != nil {
		t.Errorf("SelectedTaskID = %d, expected nil when updating by ID", *state.SelectedTaskID)
	}
}

func TestExecutePronounSubstitution(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "test.sqlite")
	st, err := store.Open(ctx, store.Options{Path: dbPath})
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })

	svc := service.NewTaskService(st)

	// Create two tasks
	task1, err := svc.CreateTask(ctx, service.CreateTaskRequest{Title: "first task"})
	if err != nil {
		t.Fatalf("create first task: %v", err)
	}
	task2, err := svc.CreateTask(ctx, service.CreateTaskRequest{Title: "second task"})
	if err != nil {
		t.Fatalf("create second task: %v", err)
	}

	state := &shell.SessionState{
		LastTaskIDs: []int64{task1.ID, task2.ID},
	}
	exec := shell.NewExecutor(svc, state, true)

	// "it" and "this" should refer to last task (task2)
	_, err = exec.Execute(ctx, "set it title:updated via it")
	if err != nil {
		t.Fatalf("execute with 'it': %v", err)
	}

	updated, err := svc.GetTask(ctx, task2.ID)
	if err != nil {
		t.Fatalf("get updated task: %v", err)
	}
	if updated.Title != "updated via it" {
		t.Errorf("task2 title = %q, want 'updated via it'", updated.Title)
	}

	// "last" should also refer to task2
	_, err = exec.Execute(ctx, "set last title:updated via last")
	if err != nil {
		t.Fatalf("execute with 'last': %v", err)
	}

	updated, err = svc.GetTask(ctx, task2.ID)
	if err != nil {
		t.Fatalf("get updated task: %v", err)
	}
	if updated.Title != "updated via last" {
		t.Errorf("task2 title = %q, want 'updated via last'", updated.Title)
	}
}

func TestExecuteThatPronounSubstitution(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "test.sqlite")
	st, err := store.Open(ctx, store.Options{Path: dbPath})
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })

	svc := service.NewTaskService(st)

	// Create two tasks
	task1, err := svc.CreateTask(ctx, service.CreateTaskRequest{Title: "first task"})
	if err != nil {
		t.Fatalf("create first task: %v", err)
	}
	task2, err := svc.CreateTask(ctx, service.CreateTaskRequest{Title: "second task"})
	if err != nil {
		t.Fatalf("create second task: %v", err)
	}

	state := &shell.SessionState{
		LastTaskIDs: []int64{task1.ID, task2.ID},
	}
	exec := shell.NewExecutor(svc, state, true)

	// "that" should refer to second-to-last task (task1)
	_, err = exec.Execute(ctx, "set that title:updated via that")
	if err != nil {
		t.Fatalf("execute with 'that': %v", err)
	}

	updated, err := svc.GetTask(ctx, task1.ID)
	if err != nil {
		t.Fatalf("get updated task: %v", err)
	}
	if updated.Title != "updated via that" {
		t.Errorf("task1 title = %q, want 'updated via that'", updated.Title)
	}
}

func TestExecuteSelectedPronounSubstitution(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "test.sqlite")
	st, err := store.Open(ctx, store.Options{Path: dbPath})
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })

	svc := service.NewTaskService(st)
	task, err := svc.CreateTask(ctx, service.CreateTaskRequest{Title: "task"})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	state := &shell.SessionState{
		SelectedTaskID: &task.ID,
	}
	exec := shell.NewExecutor(svc, state, true)

	// "selected" should refer to selected task
	_, err = exec.Execute(ctx, "set selected title:updated via selected")
	if err != nil {
		t.Fatalf("execute with 'selected': %v", err)
	}

	updated, err := svc.GetTask(ctx, task.ID)
	if err != nil {
		t.Fatalf("get updated task: %v", err)
	}
	if updated.Title != "updated via selected" {
		t.Errorf("task title = %q, want 'updated via selected'", updated.Title)
	}
}

func TestExecuteCreateWithoutContext(t *testing.T) {
	t.Parallel()

	svc := &recordingService{}
	exec := shell.NewExecutor(svc, &shell.SessionState{}, true)

	_, err := exec.Execute(context.Background(), "add buy milk")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if len(svc.lastCreate.Projects) != 0 {
		t.Errorf("projects = %v, want empty", svc.lastCreate.Projects)
	}
	if len(svc.lastCreate.Contexts) != 0 {
		t.Errorf("contexts = %v, want empty", svc.lastCreate.Contexts)
	}
}

func TestExecuteFilterWithoutContext(t *testing.T) {
	t.Parallel()

	svc := &recordingService{}
	exec := shell.NewExecutor(svc, &shell.SessionState{}, true)

	_, err := exec.Execute(context.Background(), "find state:now")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	// Without context, should not inject project/context predicates
	hasProject := hasPredicateKind(svc.lastFilter.Filter, nlp.PredProject)
	hasContext := hasPredicateKind(svc.lastFilter.Filter, nlp.PredContext)

	if hasProject {
		t.Error("filter has project predicate, expected none without context")
	}
	if hasContext {
		t.Error("filter has context predicate, expected none without context")
	}
}

func TestExecuteInjectContextOnlyProject(t *testing.T) {
	t.Parallel()

	svc := &recordingService{}
	exec := shell.NewExecutor(svc, &shell.SessionState{ContextProject: "work"}, true)

	_, err := exec.Execute(context.Background(), "add buy milk")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if !contains(svc.lastCreate.Projects, "work") {
		t.Errorf("projects = %v, want 'work'", svc.lastCreate.Projects)
	}
	if len(svc.lastCreate.Contexts) != 0 {
		t.Errorf("contexts = %v, want empty when only project context set", svc.lastCreate.Contexts)
	}
}

func TestExecuteInjectContextOnlyContext(t *testing.T) {
	t.Parallel()

	svc := &recordingService{}
	exec := shell.NewExecutor(svc, &shell.SessionState{ContextContext: "phone"}, true)

	_, err := exec.Execute(context.Background(), "add buy milk")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if len(svc.lastCreate.Projects) != 0 {
		t.Errorf("projects = %v, want empty when only context context set", svc.lastCreate.Projects)
	}
	if !contains(svc.lastCreate.Contexts, "phone") {
		t.Errorf("contexts = %v, want 'phone'", svc.lastCreate.Contexts)
	}
}
