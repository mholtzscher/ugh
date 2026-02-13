package store

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	sq "github.com/Masterminds/squirrel"

	"github.com/mholtzscher/ugh/internal/nlp"
)

type filterSQLBuilder struct{}

func (b *filterSQLBuilder) Build(expr nlp.FilterExpr) (string, []any, error) {
	sqlizer, err := b.buildExpr(expr)
	if err != nil {
		return "", nil, err
	}

	clause, args, err := sqlizer.ToSql()
	if err != nil {
		return "", nil, err
	}
	return clause, args, nil
}

func (b *filterSQLBuilder) buildExpr(expr nlp.FilterExpr) (sq.Sqlizer, error) {
	switch typed := expr.(type) {
	case nlp.Predicate:
		return b.buildPredicate(typed)
	case nlp.FilterBinary:
		left, err := b.buildExpr(typed.Left)
		if err != nil {
			return nil, err
		}
		right, err := b.buildExpr(typed.Right)
		if err != nil {
			return nil, err
		}
		if typed.Op == nlp.FilterOr {
			return sq.Or{left, right}, nil
		}
		return sq.And{left, right}, nil
	case nlp.FilterNot:
		inner, err := b.buildExpr(typed.Expr)
		if err != nil {
			return nil, err
		}
		return sq.Expr("(NOT (?))", inner), nil
	default:
		return nil, fmt.Errorf("unsupported filter expression type %T", expr)
	}
}

func (b *filterSQLBuilder) buildPredicate(pred nlp.Predicate) (sq.Sqlizer, error) {
	value := strings.TrimSpace(pred.Text)

	switch pred.Kind {
	case nlp.PredState:
		return sq.Eq{"t.state": value}, nil
	case nlp.PredDue:
		if value == nlp.FilterWildcard {
			return sq.Expr("(t.due_on IS NOT NULL AND t.due_on != '')"), nil
		}
		return sq.Eq{"t.due_on": value}, nil
	case nlp.PredProject:
		if value != nlp.FilterWildcard {
			return sq.Expr(
				"EXISTS (SELECT 1 FROM json_each(t.projects_json) WHERE value = ?)",
				value,
			), nil
		}
		return sq.Expr("json_array_length(t.projects_json) > 0"), nil
	case nlp.PredContext:
		if value != nlp.FilterWildcard {
			return sq.Expr(
				"EXISTS (SELECT 1 FROM json_each(t.contexts_json) WHERE value = ?)",
				value,
			), nil
		}
		return sq.Expr("json_array_length(t.contexts_json) > 0"), nil
	case nlp.PredText:
		if value == "" {
			return sq.Expr("1=1"), nil
		}
		like := "%" + value + "%"

		return sq.Or{
			sq.Like{"t.title": like},
			sq.Like{"t.notes": like},
			sq.Expr("EXISTS (SELECT 1 FROM json_each(t.projects_json) WHERE value LIKE ?)", like),
			sq.Expr("EXISTS (SELECT 1 FROM json_each(t.contexts_json) WHERE value LIKE ?)", like),
			sq.Expr("EXISTS (SELECT 1 FROM json_each(t.meta_json) WHERE key LIKE ? OR value LIKE ?)", like, like),
		}, nil
	case nlp.PredID:
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil || id <= 0 {
			return nil, fmt.Errorf("invalid id predicate %q", pred.Text)
		}
		return sq.Eq{"t.id": id}, nil
	case nlp.PredRecent:
		return nil, errors.New("recent modifier must be stripped before SQL build")
	default:
		return nil, fmt.Errorf("unsupported predicate kind %v", pred.Kind)
	}
}
