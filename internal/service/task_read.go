package service

import (
	"context"
	"strings"

	"github.com/mholtzscher/ugh/internal/domain"
	"github.com/mholtzscher/ugh/internal/nlp"
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

	if req.Filter != nil {
		if req.All {
			return s.store.ListTasksByExpr(ctx, req.Filter, store.ListTasksByExprOptions{})
		}
		if req.DoneOnly {
			return s.store.ListTasksByExpr(ctx, req.Filter, store.ListTasksByExprOptions{OnlyDone: true})
		}
		if req.TodoOnly {
			return s.store.ListTasksByExpr(ctx, req.Filter, store.ListTasksByExprOptions{ExcludeDone: true})
		}

		excludeDone := !exprReferencesStateDone(req.Filter) && !exprReferencesID(req.Filter)
		return s.store.ListTasksByExpr(ctx, req.Filter, store.ListTasksByExprOptions{ExcludeDone: excludeDone})
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

func exprReferencesStateDone(expr nlp.FilterExpr) bool {
	switch typed := expr.(type) {
	case nlp.Predicate:
		return typed.Kind == nlp.PredState && strings.EqualFold(strings.TrimSpace(typed.Text), domain.TaskStateDone)
	case nlp.FilterBinary:
		return exprReferencesStateDone(typed.Left) || exprReferencesStateDone(typed.Right)
	case nlp.FilterNot:
		return exprReferencesStateDone(typed.Expr)
	default:
		return false
	}
}

func exprReferencesID(expr nlp.FilterExpr) bool {
	switch typed := expr.(type) {
	case nlp.Predicate:
		return typed.Kind == nlp.PredID
	case nlp.FilterBinary:
		return exprReferencesID(typed.Left) || exprReferencesID(typed.Right)
	case nlp.FilterNot:
		return exprReferencesID(typed.Expr)
	default:
		return false
	}
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
