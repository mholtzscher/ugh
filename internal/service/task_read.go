package service

import (
	"context"
	"strings"

	"github.com/mholtzscher/ugh/internal/domain"
	"github.com/mholtzscher/ugh/internal/nlp"
	"github.com/mholtzscher/ugh/internal/store"
)

func (s *TaskService) ListTasks(ctx context.Context, req ListTasksRequest) ([]*store.Task, error) {
	expr := req.Filter

	opts := store.ListTasksByExprOptions{}
	switch {
	case req.All:
		// no-op
	case req.DoneOnly:
		opts.OnlyDone = true
	case req.TodoOnly:
		opts.ExcludeDone = true
	case expr == nil:
		opts.ExcludeDone = true
	default:
		opts.ExcludeDone = !exprReferencesStateDone(expr) && !exprReferencesID(expr)
	}

	return s.store.ListTasksByExpr(ctx, expr, opts)
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
