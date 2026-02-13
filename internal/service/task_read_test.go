//nolint:testpackage // Tests exercise unexported filter helper behavior directly.
package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
			assert.Equal(t, tt.want, got, "exprReferencesStateDone() mismatch")
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
			assert.Equal(t, tt.want, got, "exprReferencesID() mismatch")
		})
	}
}

func TestStripRecentModifier_RemovesRecentPredicate(t *testing.T) {
	t.Parallel()

	expr := nlp.FilterBinary{
		Op:    nlp.FilterAnd,
		Left:  nlp.Predicate{Kind: nlp.PredProject, Text: "work"},
		Right: nlp.Predicate{Kind: nlp.PredRecent, Text: "10"},
	}

	stripped, recent, limit, err := stripRecentModifier(expr, false)
	require.NoError(t, err, "stripRecentModifier() error")
	assert.True(t, recent, "recent should be detected")
	assert.EqualValues(t, 10, limit, "recent limit mismatch")

	pred, ok := stripped.(nlp.Predicate)
	assert.True(t, ok, "stripped expression should collapse to predicate")
	assert.Equal(t, nlp.PredProject, pred.Kind, "predicate kind mismatch")
	assert.Equal(t, "work", pred.Text, "predicate text mismatch")
}

func TestStripRecentModifier_NotRecentReturnsError(t *testing.T) {
	t.Parallel()

	_, _, _, err := stripRecentModifier(
		nlp.FilterNot{Expr: nlp.Predicate{Kind: nlp.PredRecent, Text: ""}},
		false,
	)
	require.Error(t, err, "stripRecentModifier() should reject negated recent")
	assert.Contains(t, err.Error(), "cannot be negated", "error should mention negation")
}

func TestStripRecentModifier_ConflictingLimitsReturnsError(t *testing.T) {
	t.Parallel()

	_, _, _, err := stripRecentModifier(
		nlp.FilterBinary{
			Op:    nlp.FilterAnd,
			Left:  nlp.Predicate{Kind: nlp.PredRecent, Text: "5"},
			Right: nlp.Predicate{Kind: nlp.PredRecent, Text: "10"},
		},
		false,
	)
	require.Error(t, err, "stripRecentModifier() should reject conflicting limits")
	assert.Contains(t, err.Error(), "conflicting recent limits", "error should mention conflict")
}
