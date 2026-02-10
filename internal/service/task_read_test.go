//nolint:testpackage // Tests exercise unexported filter helper behavior directly.
package service

import (
	"testing"

	"github.com/mholtzscher/ugh/internal/nlp"
)

func TestExprReferencesStateDone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		expr nlp.FilterExpr
		want bool
	}{
		{
			name: "single done state",
			expr: nlp.Predicate{Kind: nlp.PredState, Text: "done"},
			want: true,
		},
		{
			name: "single non-done state",
			expr: nlp.Predicate{Kind: nlp.PredState, Text: "now"},
			want: false,
		},
		{
			name: "not done",
			expr: nlp.FilterNot{Expr: nlp.Predicate{Kind: nlp.PredState, Text: "done"}},
			want: true,
		},
		{
			name: "nested binary includes done",
			expr: nlp.FilterBinary{
				Op:   nlp.FilterAnd,
				Left: nlp.Predicate{Kind: nlp.PredProject, Text: "work"},
				Right: nlp.FilterBinary{
					Op:    nlp.FilterOr,
					Left:  nlp.Predicate{Kind: nlp.PredText, Text: "paper"},
					Right: nlp.Predicate{Kind: nlp.PredState, Text: "done"},
				},
			},
			want: true,
		},
		{
			name: "nil expression",
			expr: nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := exprReferencesStateDone(tt.expr)
			if got != tt.want {
				t.Fatalf("exprReferencesStateDone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExprReferencesID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		expr nlp.FilterExpr
		want bool
	}{
		{
			name: "single id predicate",
			expr: nlp.Predicate{Kind: nlp.PredID, Text: "42"},
			want: true,
		},
		{
			name: "non-id predicate",
			expr: nlp.Predicate{Kind: nlp.PredState, Text: "done"},
			want: false,
		},
		{
			name: "nested contains id",
			expr: nlp.FilterBinary{
				Op:    nlp.FilterAnd,
				Left:  nlp.Predicate{Kind: nlp.PredProject, Text: "work"},
				Right: nlp.FilterNot{Expr: nlp.Predicate{Kind: nlp.PredID, Text: "1"}},
			},
			want: true,
		},
		{
			name: "nil expression",
			expr: nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := exprReferencesID(tt.expr)
			if got != tt.want {
				t.Fatalf("exprReferencesID() = %v, want %v", got, tt.want)
			}
		})
	}
}
