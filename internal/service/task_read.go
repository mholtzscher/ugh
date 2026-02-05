package service

import (
	"context"
	"strings"

	"github.com/mholtzscher/ugh/internal/store"
)

func (s *TaskService) ListTasks(ctx context.Context, req ListTasksRequest) ([]*store.Task, error) {
	filters := store.Filters{
		All:        req.All,
		DoneOnly:   req.DoneOnly,
		TodoOnly:   req.TodoOnly,
		State:      strings.TrimSpace(req.State),
		Project:    req.Project,
		Context:    req.Context,
		Search:     req.Search,
		DueSetOnly: req.DueOnly,
	}

	if !filters.All && !filters.DoneOnly && !filters.TodoOnly {
		filters.TodoOnly = true
	}

	return s.store.ListTasks(ctx, filters)
}

func (s *TaskService) GetTask(ctx context.Context, id int64) (*store.Task, error) {
	return s.store.GetTask(ctx, id)
}

func (s *TaskService) ListProjects(ctx context.Context, req ListTagsRequest) ([]store.NameCount, error) {
	onlyDone := req.DoneOnly
	excludeDone := req.TodoOnly
	if !req.All && !req.DoneOnly && !req.TodoOnly {
		excludeDone = true
	}
	return s.store.ListProjectCounts(ctx, onlyDone, excludeDone)
}

func (s *TaskService) ListContexts(ctx context.Context, req ListTagsRequest) ([]store.NameCount, error) {
	onlyDone := req.DoneOnly
	excludeDone := req.TodoOnly
	if !req.All && !req.DoneOnly && !req.TodoOnly {
		excludeDone = true
	}
	return s.store.ListContextCounts(ctx, onlyDone, excludeDone)
}
