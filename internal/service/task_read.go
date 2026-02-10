package service

import (
	"context"
	"strings"

	"github.com/mholtzscher/ugh/internal/store"
)

func (s *TaskService) ListTasks(ctx context.Context, req ListTasksRequest) ([]*store.Task, error) {
	// If specific IDs requested, fetch those tasks directly
	if len(req.IDs) > 0 {
		var tasks []*store.Task
		for _, id := range req.IDs {
			task, err := s.GetTask(ctx, id)
			if err != nil {
				return nil, err
			}
			if task != nil {
				tasks = append(tasks, task)
			}
		}
		return tasks, nil
	}

	filters := store.Filters{
		All:        req.All,
		DoneOnly:   req.DoneOnly,
		TodoOnly:   req.TodoOnly,
		States:     uniqueStrings(req.States),
		Projects:   uniqueStrings(req.Projects),
		Contexts:   uniqueStrings(req.Contexts),
		Search:     uniqueStrings(req.Search),
		DueSetOnly: req.DueOnly,
		DueOn:      req.DueOn,
	}

	if !filters.All && !filters.DoneOnly && !filters.TodoOnly {
		filters.TodoOnly = true
	}

	return s.store.ListTasks(ctx, filters)
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(values))
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v == "" || seen[v] {
			continue
		}
		seen[v] = true
		result = append(result, v)
	}
	return result
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
