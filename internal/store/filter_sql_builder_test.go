//nolint:testpackage // Tests exercise unexported SQL builder behavior directly.
package store

import (
	"strings"
	"testing"

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
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	if !strings.Contains(clause, "t.state = ?") {
		t.Fatalf("clause = %q, want state predicate", clause)
	}
	if !strings.Contains(clause, "NOT") {
		t.Fatalf("clause = %q, want NOT operator", clause)
	}
	if !strings.Contains(clause, "task_project_links") {
		t.Fatalf("clause = %q, want project EXISTS subquery", clause)
	}

	if len(args) != 2 {
		t.Fatalf("args len = %d, want 2", len(args))
	}
	if args[0] != "now" || args[1] != "work" {
		t.Fatalf("args = %#v, want [now work]", args)
	}
}

func TestFilterSQLBuilder_TextSearchAddsLikeArgs(t *testing.T) {
	t.Parallel()

	b := &filterSQLBuilder{}
	clause, args, err := b.Build(nlp.Predicate{Kind: nlp.PredText, Text: "paper"})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	if !strings.Contains(clause, "t.title LIKE ?") {
		t.Fatalf("clause = %q, want title LIKE", clause)
	}
	if len(args) != 6 {
		t.Fatalf("args len = %d, want 6", len(args))
	}
	for i, arg := range args {
		if arg != "%paper%" {
			t.Fatalf("args[%d] = %#v, want %%paper%%", i, arg)
		}
	}
}

func TestFilterSQLBuilder_InvalidIDReturnsError(t *testing.T) {
	t.Parallel()

	b := &filterSQLBuilder{}
	_, _, err := b.Build(nlp.Predicate{Kind: nlp.PredID, Text: "abc"})
	if err == nil {
		t.Fatal("Build() error = nil, want invalid id error")
	}
}

func TestFilterSQLBuilder_EmptyDuePredicateHasNoArgs(t *testing.T) {
	t.Parallel()

	b := &filterSQLBuilder{}
	clause, args, err := b.Build(nlp.Predicate{Kind: nlp.PredDue, Text: ""})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	if !strings.Contains(clause, "t.due_on IS NOT NULL") {
		t.Fatalf("clause = %q, want IS NOT NULL", clause)
	}
	if !strings.Contains(clause, "t.due_on != ''") {
		t.Fatalf("clause = %q, want non-empty due check", clause)
	}
	if len(args) != 0 {
		t.Fatalf("args len = %d, want 0", len(args))
	}
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
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	if !strings.Contains(clause, "NOT") {
		t.Fatalf("clause = %q, want NOT", clause)
	}
	if !strings.Contains(clause, "t.state = ? OR t.state = ?") {
		t.Fatalf("clause = %q, want nested OR", clause)
	}
	if len(args) != 2 {
		t.Fatalf("args len = %d, want 2", len(args))
	}
	if args[0] != "now" || args[1] != "waiting" {
		t.Fatalf("args = %#v, want [now waiting]", args)
	}
}
