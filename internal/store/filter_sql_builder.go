package store

import (
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
		if value == "" {
			return sq.Expr("(t.due_on IS NOT NULL AND t.due_on != '')"), nil
		}
		return sq.Eq{"t.due_on": value}, nil
	case nlp.PredProject:
		subquery := sq.Select("1").
			From("task_project_links tpl").
			Join("projects p ON p.id = tpl.project_id").
			Where(sq.Expr("tpl.task_id = t.id")).
			Where(sq.Eq{"p.name": value})
		return sq.Expr("EXISTS (?)", subquery), nil
	case nlp.PredContext:
		subquery := sq.Select("1").
			From("task_context_links tcl").
			Join("contexts c ON c.id = tcl.context_id").
			Where(sq.Expr("tcl.task_id = t.id")).
			Where(sq.Eq{"c.name": value})
		return sq.Expr("EXISTS (?)", subquery), nil
	case nlp.PredText:
		if value == "" {
			return sq.Expr("1=1"), nil
		}
		like := "%" + value + "%"

		projectSubquery := sq.Select("1").
			From("task_project_links tpl").
			Join("projects p ON p.id = tpl.project_id").
			Where(sq.Expr("tpl.task_id = t.id")).
			Where(sq.Like{"p.name": like})

		contextSubquery := sq.Select("1").
			From("task_context_links tcl").
			Join("contexts c ON c.id = tcl.context_id").
			Where(sq.Expr("tcl.task_id = t.id")).
			Where(sq.Like{"c.name": like})

		metaSubquery := sq.Select("1").
			From("task_meta m").
			Where(sq.Expr("m.task_id = t.id")).
			Where(sq.Or{sq.Like{"m.key": like}, sq.Like{"m.value": like}})

		return sq.Or{
			sq.Like{"t.title": like},
			sq.Like{"t.notes": like},
			sq.Expr("EXISTS (?)", projectSubquery),
			sq.Expr("EXISTS (?)", contextSubquery),
			sq.Expr("EXISTS (?)", metaSubquery),
		}, nil
	case nlp.PredID:
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil || id <= 0 {
			return nil, fmt.Errorf("invalid id predicate %q", pred.Text)
		}
		return sq.Eq{"t.id": id}, nil
	default:
		return nil, fmt.Errorf("unsupported predicate kind %v", pred.Kind)
	}
}
