package cmd

import (
	"strings"

	"github.com/mholtzscher/ugh/internal/nlp"
)

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
	value = strings.ToLower(value)
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

func dueSetExpr() nlp.FilterExpr {
	return nlp.Predicate{Kind: nlp.PredDue, Text: ""}
}
