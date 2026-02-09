package nlp

import (
	"github.com/alecthomas/participle/v2/lexer"
)

//nolint:gochecknoglobals // DSL lexer must be a package-level var for participle
var dslLexer = lexer.MustSimple([]lexer.SimpleRule{
	// Quoted strings first (highest priority)
	{Name: "Quoted", Pattern: `"(?:[^"\\]|\\.)*"`},

	// Tags - project (#) and context (@)
	{Name: "ProjectTag", Pattern: `#[a-zA-Z_][a-zA-Z0-9_-]*`},
	{Name: "ContextTag", Pattern: `@[a-zA-Z_][a-zA-Z0-9_-]*`},

	// Field setters with colon (op starters) - MUST come before Ident
	// These consume the field name and colon together
	{
		Name:    "SetField",
		Pattern: `\b(title|notes|due|waiting|state|project|projects|context|contexts|meta|id|text)\b\s*:`,
	},
	{
		Name:    "AddField",
		Pattern: `\+\s*\b(project|projects|context|contexts|meta)\b\s*:`,
	},
	{
		Name:    "RemoveField",
		Pattern: `-\s*\b(project|projects|context|contexts|meta)\b\s*:`,
	},
	{
		Name:    "ClearField",
		Pattern: `!\s*\b(notes|due|waiting|projects|contexts|meta)\b`,
	},

	// Clear op for non-field cases (just the ! symbol)
	{Name: "ClearOp", Pattern: `!`},

	// Add/Remove ops as standalone (for tag operations)
	{Name: "AddOp", Pattern: `\+`},
	{Name: "RemoveOp", Pattern: `-`},

	// Punctuation
	{Name: "Colon", Pattern: `:`},
	{Name: "LParen", Pattern: `\(`},
	{Name: "RParen", Pattern: `\)`},

	// Logical operators
	{Name: "AndOp", Pattern: `&&`},
	{Name: "OrOp", Pattern: `\|\|`},

	// Keywords
	{Name: "And", Pattern: `\band\b`},
	{Name: "Or", Pattern: `\bor\b`},
	{Name: "Not", Pattern: `\bnot\b`},

	// Verbs (commands)
	{Name: "Verb", Pattern: `\b(add|create|new|set|edit|update|find|show|list|filter)\b`},

	// Relative date keywords
	{Name: "RelativeDate", Pattern: `\b(today|tomorrow|next-week)\b`},

	// Target keywords
	{Name: "Target", Pattern: `\b(selected|it|this|that)\b`},

	// Identifiers and words (catch-all for regular words including alphanumeric)
	{Name: "Ident", Pattern: `[a-zA-Z0-9_-]+`},

	// Whitespace (elided)
	{Name: "Whitespace", Pattern: `\s+`},
})
