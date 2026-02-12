//nolint:testpackage // Tests exercise unexported SQL builder behavior directly.
package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mholtzscher/ugh/internal/nlp"
)

func TestFilterSQLBuilder_BuildBooleanExpression(t *testing.T) {
	t.Parallel()

	b := &filterSQLBuilder{}
	expr := nlp.FilterBinary{
		Op: nlp.FilterOr,
		Left: nlp.Predicate{
			Kind: nlp.PredState,
			Text: "now",
		},
		Right: nlp.FilterNot{
			Expr: nlp.Predicate{Kind: nlp.PredProject, Text: "work"},
		},
	}

	clause, args, err := b.Build(expr)
	require.NoError(t, err, "Build() error")

	assert.Contains(t, clause, "t.state = ?", "clause should contain state predicate")
	assert.Contains(t, clause, "NOT", "clause should contain NOT operator")
	assert.Contains(t, clause, "task_project_links", "clause should contain project EXISTS subquery")
	require.Len(t, args, 2, "args len mismatch")
	assert.Equal(t, "now", args[0], "args[0] mismatch")
	assert.Equal(t, "work", args[1], "args[1] mismatch")
}

func TestFilterSQLBuilder_TextSearchAddsLikeArgs(t *testing.T) {
	t.Parallel()

	b := &filterSQLBuilder{}
	clause, args, err := b.Build(nlp.Predicate{Kind: nlp.PredText, Text: "paper"})
	require.NoError(t, err, "Build() error")

	assert.Contains(t, clause, "t.title LIKE ?", "clause should contain title LIKE")
	require.Len(t, args, 6, "args len mismatch")
	for i, arg := range args {
		assert.Equal(t, "%paper%", arg, "args[%d] mismatch", i)
	}
}

func TestFilterSQLBuilder_InvalidIDReturnsError(t *testing.T) {
	t.Parallel()

	b := &filterSQLBuilder{}
	_, _, err := b.Build(nlp.Predicate{Kind: nlp.PredID, Text: "abc"})
	require.Error(t, err, "Build() should return invalid id error")
}

func TestFilterSQLBuilder_WildcardDuePredicateHasNoArgs(t *testing.T) {
	t.Parallel()

	b := &filterSQLBuilder{}
	clause, args, err := b.Build(nlp.Predicate{Kind: nlp.PredDue, Text: nlp.FilterWildcard})
	require.NoError(t, err, "Build() error")

	assert.Contains(t, clause, "t.due_on IS NOT NULL", "clause should contain IS NOT NULL")
	assert.Contains(t, clause, "t.due_on != ''", "clause should contain non-empty due check")
	assert.Empty(t, args, "args should be empty")
}

func TestFilterSQLBuilder_WildcardProjectPredicateHasNoNameArg(t *testing.T) {
	t.Parallel()

	b := &filterSQLBuilder{}
	clause, args, err := b.Build(nlp.Predicate{Kind: nlp.PredProject, Text: nlp.FilterWildcard})
	require.NoError(t, err, "Build() error")

	assert.Contains(t, clause, "task_project_links", "clause should include project link subquery")
	assert.NotContains(t, clause, "p.name = ?", "clause should not constrain project name")
	assert.Empty(t, args, "args should be empty")
}

func TestFilterSQLBuilder_WildcardContextPredicateHasNoNameArg(t *testing.T) {
	t.Parallel()

	b := &filterSQLBuilder{}
	clause, args, err := b.Build(nlp.Predicate{Kind: nlp.PredContext, Text: nlp.FilterWildcard})
	require.NoError(t, err, "Build() error")

	assert.Contains(t, clause, "task_context_links", "clause should include context link subquery")
	assert.NotContains(t, clause, "c.name = ?", "clause should not constrain context name")
	assert.Empty(t, args, "args should be empty")
}

func TestFilterSQLBuilder_NotWrapsNestedOrExpression(t *testing.T) {
	t.Parallel()

	b := &filterSQLBuilder{}
	expr := nlp.FilterNot{
		Expr: nlp.FilterBinary{
			Op:    nlp.FilterOr,
			Left:  nlp.Predicate{Kind: nlp.PredState, Text: "now"},
			Right: nlp.Predicate{Kind: nlp.PredState, Text: "waiting"},
		},
	}

	clause, args, err := b.Build(expr)
	require.NoError(t, err, "Build() error")

	assert.Contains(t, clause, "NOT", "clause should contain NOT")
	assert.Contains(t, clause, "t.state = ? OR t.state = ?", "clause should contain nested OR")
	require.Len(t, args, 2, "args len mismatch")
	assert.Equal(t, "now", args[0], "args[0] mismatch")
	assert.Equal(t, "waiting", args[1], "args[1] mismatch")
}
