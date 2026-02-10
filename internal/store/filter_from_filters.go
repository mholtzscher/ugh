package store

import (
	"strings"

	"github.com/mholtzscher/ugh/internal/nlp"
)

func buildExprFromFilters(filters Filters) nlp.FilterExpr {
	parts := make([]nlp.FilterExpr, 0)

	if stateExpr := buildOrPredicates(nlp.PredState, filters.States); stateExpr != nil {
		parts = append(parts, stateExpr)
	}
	if projectExpr := buildOrPredicates(nlp.PredProject, filters.Projects); projectExpr != nil {
		parts = append(parts, projectExpr)
	}
	if contextExpr := buildOrPredicates(nlp.PredContext, filters.Contexts); contextExpr != nil {
		parts = append(parts, contextExpr)
	}

	for _, term := range filters.Search {
		trimmed := strings.TrimSpace(term)
		if trimmed == "" {
			continue
		}
		parts = append(parts, nlp.Predicate{Kind: nlp.PredText, Text: trimmed})
	}

	if due := strings.TrimSpace(filters.DueOn); due != "" {
		parts = append(parts, nlp.Predicate{Kind: nlp.PredDue, Text: due})
	} else if filters.DueSetOnly {
		parts = append(parts, nlp.Predicate{Kind: nlp.PredDue, Text: ""})
	}

	return foldWithOperator(parts, nlp.FilterAnd)
}

func buildOrPredicates(kind nlp.PredicateKind, values []string) nlp.FilterExpr {
	preds := make([]nlp.FilterExpr, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		preds = append(preds, nlp.Predicate{Kind: kind, Text: trimmed})
	}
	return foldWithOperator(preds, nlp.FilterOr)
}

func foldWithOperator(exprs []nlp.FilterExpr, op nlp.FilterBoolOp) nlp.FilterExpr {
	if len(exprs) == 0 {
		return nil
	}
	if len(exprs) == 1 {
		return exprs[0]
	}

	result := exprs[0]
	for _, expr := range exprs[1:] {
		result = nlp.FilterBinary{Op: op, Left: result, Right: expr}
	}
	return result
}
