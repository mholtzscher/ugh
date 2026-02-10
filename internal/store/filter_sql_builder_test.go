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
