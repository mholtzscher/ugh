package store

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mholtzscher/ugh/internal/nlp"
)

type filterSQLBuilder struct {
	args []any
}

func (b *filterSQLBuilder) Build(expr nlp.FilterExpr) (string, []any, error) {
	clause, err := b.buildExpr(expr)
	if err != nil {
		return "", nil, err
	}
	return clause, b.args, nil
}

func (b *filterSQLBuilder) buildExpr(expr nlp.FilterExpr) (string, error) {
	switch typed := expr.(type) {
	case nlp.Predicate:
		return b.buildPredicate(typed)
	case nlp.FilterBinary:
		left, err := b.buildExpr(typed.Left)
		if err != nil {
			return "", err
		}
		right, err := b.buildExpr(typed.Right)
		if err != nil {
			return "", err
		}
		op := "AND"
		if typed.Op == nlp.FilterOr {
			op = "OR"
		}
		return fmt.Sprintf("(%s %s %s)", left, op, right), nil
	case nlp.FilterNot:
		inner, err := b.buildExpr(typed.Expr)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("(NOT %s)", inner), nil
	default:
		return "", fmt.Errorf("unsupported filter expression type %T", expr)
	}
}

func (b *filterSQLBuilder) buildPredicate(pred nlp.Predicate) (string, error) {
	value := strings.TrimSpace(pred.Text)

	switch pred.Kind {
	case nlp.PredState:
		b.args = append(b.args, value)
		return "(t.state = ?)", nil
	case nlp.PredDue:
		if value == "" {
			return "(t.due_on IS NOT NULL AND t.due_on != '')", nil
		}
		b.args = append(b.args, value)
		return "(t.due_on = ?)", nil
	case nlp.PredProject:
		b.args = append(b.args, value)
		return `EXISTS (
			SELECT 1
			FROM task_project_links tpl
			JOIN projects p ON p.id = tpl.project_id
			WHERE tpl.task_id = t.id AND p.name = ?
		)`, nil
	case nlp.PredContext:
		b.args = append(b.args, value)
		return `EXISTS (
			SELECT 1
			FROM task_context_links tcl
			JOIN contexts c ON c.id = tcl.context_id
			WHERE tcl.task_id = t.id AND c.name = ?
		)`, nil
	case nlp.PredText:
		if value == "" {
			return "1=1", nil
		}
		like := "%" + value + "%"
		b.args = append(b.args, like, like, like, like, like, like)
		return `(
			t.title LIKE ?
			OR t.notes LIKE ?
			OR EXISTS (
				SELECT 1
				FROM task_project_links tpl
				JOIN projects p ON p.id = tpl.project_id
				WHERE tpl.task_id = t.id AND p.name LIKE ?
			)
			OR EXISTS (
				SELECT 1
				FROM task_context_links tcl
				JOIN contexts c ON c.id = tcl.context_id
				WHERE tcl.task_id = t.id AND c.name LIKE ?
			)
			OR EXISTS (
				SELECT 1
				FROM task_meta m
				WHERE m.task_id = t.id AND (
					m.key LIKE ?
					OR m.value LIKE ?
				)
			)
		)`, nil
	case nlp.PredID:
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil || id <= 0 {
			return "", fmt.Errorf("invalid id predicate %q", pred.Text)
		}
		b.args = append(b.args, id)
		return "(t.id = ?)", nil
	default:
		return "", fmt.Errorf("unsupported predicate kind %v", pred.Kind)
	}
}
