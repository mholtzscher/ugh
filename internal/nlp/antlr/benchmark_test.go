package antlr_test

import (
	"testing"

	"github.com/mholtzscher/ugh/internal/nlp"
	antlrparser "github.com/mholtzscher/ugh/internal/nlp/antlr"
)

// Benchmark inputs representing various DSL complexity levels.
var benchInputs = []struct { //nolint:gochecknoglobals // benchmark test data
	name  string
	input string
}{
	{
		name:  "simple_create",
		input: "add buy milk",
	},
	{
		name:  "create_with_tags",
		input: "add buy milk #groceries @store due:tomorrow",
	},
	{
		name:  "complex_create",
		input: `add buy organic whole milk due:tomorrow #groceries @store waiting:alex notes:organic preferred state:now`,
	},
	{
		name:  "simple_update",
		input: "set selected state:now",
	},
	{
		name:  "complex_update",
		input: "set #123 state:now +project:work -context:old !due title:new important task",
	},
	{
		name:  "simple_filter",
		input: "find state:now",
	},
	{
		name:  "complex_filter_and",
		input: "find state:now and project:work",
	},
	{
		name:  "complex_filter_boolean",
		input: "find (state:now or state:waiting) and project:work and not #personal",
	},
	{
		name:  "view_command",
		input: "view inbox",
	},
	{
		name:  "context_command",
		input: "context #work",
	},
}

// ─── ANTLR Parser Benchmarks ────────────────────────────────────────────────

func BenchmarkANTLR_Parse(b *testing.B) {
	for _, bi := range benchInputs {
		b.Run(bi.name, func(b *testing.B) {
			opts := nlp.ParseOptions{}
			b.ResetTimer()
			b.ReportAllocs()
			for b.Loop() {
				_, _ = antlrparser.Parse(bi.input, opts)
			}
		})
	}
}

// ─── Participle Parser Benchmarks ───────────────────────────────────────────

func BenchmarkParticiple_Parse(b *testing.B) {
	for _, bi := range benchInputs {
		b.Run(bi.name, func(b *testing.B) {
			opts := nlp.ParseOptions{}
			b.ResetTimer()
			b.ReportAllocs()
			for b.Loop() {
				_, _ = nlp.Parse(bi.input, opts)
			}
		})
	}
}

// ─── Lex Benchmarks ────────────────────────────────────────────────────────

func BenchmarkANTLR_Lex(b *testing.B) {
	input := `add buy milk due:tomorrow #groceries @store waiting:alex notes:organic preferred`
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_, _ = antlrparser.Lex(input)
	}
}

func BenchmarkParticiple_Lex(b *testing.B) {
	input := `add buy milk due:tomorrow #groceries @store waiting:alex notes:organic preferred`
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_, _ = nlp.Lex(input)
	}
}
