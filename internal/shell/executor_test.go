package shell_test

import (
	"context"
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
