package nlp_test

import (
	"testing"
	"time"

	"github.com/mholtzscher/ugh/internal/nlp"
)

func TestParseCreateCommand(t *testing.T) {
	t.Parallel()

	result, err := nlp.Parse(
		`buy milk tomorrow #home @errands waiting:alex`,
		nlp.ParseOptions{Now: time.Date(2026, 2, 8, 10, 0, 0, 0, time.UTC)},
	)
	if err != nil {
		t.Fatalf("Parse(create) error = %v", err)
	}
	if result.Intent != nlp.IntentCreate {
		t.Fatalf("Parse(create) intent = %v, want %v", result.Intent, nlp.IntentCreate)
	}

	cmd, ok := result.Command.(nlp.CreateCommand)
	if !ok {
		t.Fatalf("command type = %T, want CreateCommand", result.Command)
	}
	if cmd.Title != "buy milk" {
		t.Fatalf("title = %q, want %q", cmd.Title, "buy milk")
	}
	if len(cmd.Ops) != 4 {
		t.Fatalf("ops len = %d, want 4", len(cmd.Ops))
	}
}

func TestParseUpdateCommand(t *testing.T) {
	t.Parallel()

	result, err := nlp.Parse(`set selected state:now +project:work !due`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(update) error = %v", err)
	}
	if result.Intent != nlp.IntentUpdate {
		t.Fatalf("Parse(update) intent = %v, want %v", result.Intent, nlp.IntentUpdate)
	}

	cmd, ok := result.Command.(nlp.UpdateCommand)
	if !ok {
		t.Fatalf("command type = %T, want UpdateCommand", result.Command)
	}
	if cmd.Target.Kind != nlp.TargetSelected {
		t.Fatalf("target kind = %v, want TargetSelected", cmd.Target.Kind)
	}
	if len(cmd.Ops) != 3 {
		t.Fatalf("ops len = %d, want 3", len(cmd.Ops))
	}
}

func TestParseFilterCommand(t *testing.T) {
	t.Parallel()

	result, err := nlp.Parse(`find state:now and project:work`, nlp.ParseOptions{})
	if err != nil {
		t.Fatalf("Parse(filter) error = %v", err)
	}
	if result.Intent != nlp.IntentFilter {
		t.Fatalf("Parse(filter) intent = %v, want %v", result.Intent, nlp.IntentFilter)
	}

	cmd, ok := result.Command.(nlp.FilterCommand)
	if !ok {
		t.Fatalf("command type = %T, want FilterCommand", result.Command)
	}
	binary, ok := cmd.Expr.(nlp.FilterBinary)
	if !ok {
		t.Fatalf("expr type = %T, want FilterBinary", cmd.Expr)
	}
	if binary.Op != nlp.FilterAnd {
		t.Fatalf("binary op = %v, want %v", binary.Op, nlp.FilterAnd)
	}
}

func TestParseInvalidUpdateTarget(t *testing.T) {
	t.Parallel()

	_, err := nlp.Parse(`set banana state:now`, nlp.ParseOptions{})
	if err == nil {
		t.Fatal("expected parse error for invalid update target")
	}
}
