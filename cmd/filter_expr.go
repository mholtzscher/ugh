package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mholtzscher/ugh/internal/nlp"
	"github.com/mholtzscher/ugh/internal/nlp/compile"
)

type listFilterOptions struct {
	Where   string
	State   string
	Project string
	Context string
	Search  string
	DueSet  bool
}

func buildListFilterExpr(opts listFilterOptions) (nlp.FilterExpr, error) {
	whereExpr, err := parseWhereExpr(opts.Where)
	if err != nil {
		return nil, err
	}

	expr := andExpr(
		whereExpr,
		stateExpr(opts.State),
		projectExpr(opts.Project),
		contextExpr(opts.Context),
		textExpr(opts.Search),
		dueSetExpr(opts.DueSet),
	)

	return compile.NormalizeFilterExpr(expr, compile.BuildOptions{Now: time.Now()})
}

func parseWhereExpr(where string) (nlp.FilterExpr, error) {
	where = strings.TrimSpace(where)
	if where == "" {
		var emptyExpr nlp.FilterExpr
		return emptyExpr, nil
	}

	parsed, err := nlp.Parse("find "+where, nlp.ParseOptions{Mode: nlp.ModeFilter, Now: time.Now()})
	if err != nil {
		return nil, fmt.Errorf("parse --where: %w", err)
	}

	filterCmd, ok := parsed.Command.(*nlp.FilterCommand)
	if !ok || filterCmd.Expr == nil {
		return nil, errors.New("parse --where: expected filter expression")
	}

	return filterCmd.Expr, nil
}

func andExpr(exprs ...nlp.FilterExpr) nlp.FilterExpr {
	filtered := make([]nlp.FilterExpr, 0, len(exprs))
	for _, expr := range exprs {
		if expr == nil {
			continue
		}
		filtered = append(filtered, expr)
	}

	if len(filtered) == 0 {
		return nil
	}
	result := filtered[0]
	for _, expr := range filtered[1:] {
		result = nlp.FilterBinary{Op: nlp.FilterAnd, Left: result, Right: expr}
	}
	return result
}

func stateExpr(value string) nlp.FilterExpr {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return nlp.Predicate{Kind: nlp.PredState, Text: value}
}

func projectExpr(value string) nlp.FilterExpr {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return nlp.Predicate{Kind: nlp.PredProject, Text: value}
}

func contextExpr(value string) nlp.FilterExpr {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return nlp.Predicate{Kind: nlp.PredContext, Text: value}
}

func textExpr(value string) nlp.FilterExpr {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return nlp.Predicate{Kind: nlp.PredText, Text: value}
}

func dueSetExpr(enabled bool) nlp.FilterExpr {
	if !enabled {
		return nil
	}
	return nlp.Predicate{Kind: nlp.PredDue, Text: ""}
}
