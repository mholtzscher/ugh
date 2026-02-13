package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mholtzscher/ugh/internal/domain"
	"github.com/mholtzscher/ugh/internal/nlp"
	"github.com/mholtzscher/ugh/internal/store"
)

const defaultRecentLimit int64 = 20

func (s *TaskService) ListTasks(ctx context.Context, req ListTasksRequest) ([]*store.Task, error) {
	expr, recentEnabled, recentLimit, err := stripRecentModifier(req.Filter, false)
	if err != nil {
		return nil, err
	}

	effectiveLimit := req.Limit
	if effectiveLimit == 0 {
		effectiveLimit = recentLimit
	}
	if effectiveLimit == 0 && (recentEnabled || req.Recent) {
		effectiveLimit = defaultRecentLimit
	}

	opts := store.ListTasksByExprOptions{}
	opts.Recent = recentEnabled || req.Recent
	opts.Limit = effectiveLimit
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

//nolint:gocognit // Recursive AST stripping needs explicit per-node branching.
func stripRecentModifier(expr nlp.FilterExpr, inNot bool) (nlp.FilterExpr, bool, int64, error) {
	if expr == nil {
		return nil, false, 0, nil
	}

	switch typed := expr.(type) {
	case nlp.Predicate:
		if typed.Kind != nlp.PredRecent {
			return typed, false, 0, nil
		}
		if inNot {
			return nil, false, 0, errors.New("recent modifier cannot be negated")
		}
		limit, err := parseRecentLimit(typed.Text)
		if err != nil {
			return nil, false, 0, err
		}
		return nil, true, limit, nil
	case nlp.FilterBinary:
		leftExpr, leftRecent, leftLimit, err := stripRecentModifier(typed.Left, inNot)
		if err != nil {
			return nil, false, 0, err
		}
		rightExpr, rightRecent, rightLimit, err := stripRecentModifier(typed.Right, inNot)
		if err != nil {
			return nil, false, 0, err
		}
		limit, err := mergeRecentLimits(leftLimit, rightLimit)
		if err != nil {
			return nil, false, 0, err
		}
		recent := leftRecent || rightRecent
		switch {
		case leftExpr == nil && rightExpr == nil:
			return nil, recent, limit, nil
		case leftExpr == nil:
			return rightExpr, recent, limit, nil
		case rightExpr == nil:
			return leftExpr, recent, limit, nil
		default:
			return nlp.FilterBinary{Op: typed.Op, Left: leftExpr, Right: rightExpr}, recent, limit, nil
		}
	case nlp.FilterNot:
		inner, recent, limit, err := stripRecentModifier(typed.Expr, true)
		if err != nil {
			return nil, false, 0, err
		}
		if inner == nil {
			return nil, false, 0, errors.New("recent modifier cannot be negated")
		}
		return nlp.FilterNot{Expr: inner}, recent, limit, nil
	default:
		return nil, false, 0, fmt.Errorf("unsupported filter expression type %T", expr)
	}
}

func parseRecentLimit(value string) (int64, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return 0, nil
	}
	limit, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil || limit <= 0 {
		return 0, fmt.Errorf("invalid recent limit %q", value)
	}
	return limit, nil
}

func mergeRecentLimits(left int64, right int64) (int64, error) {
	if left == 0 {
		return right, nil
	}
	if right == 0 {
		return left, nil
	}
	if left != right {
		return 0, errors.New("conflicting recent limits")
	}
	return left, nil
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
