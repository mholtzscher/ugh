package nlp

import (
	"github.com/alecthomas/participle/v2/lexer"
)

//nolint:gochecknoglobals // DSL lexer must be a package-level var for participle
var dslLexer = lexer.MustStateful(lexer.Rules{
	"Root": {
		// Quoted strings first (highest priority)
		{Name: "Quoted", Pattern: `"(?:[^"\\]|\\.)*"`},

		// Numeric hash IDs (must come before ProjectTag)
		{Name: "HashNumber", Pattern: `#[0-9]+`},

		// Tags - project (#) and context (@)
		{Name: "ProjectTag", Pattern: `#[a-zA-Z_][a-zA-Z0-9_-]*`},
		{Name: "ContextTag", Pattern: `@[a-zA-Z_][a-zA-Z0-9_-]*`},

		// Tag prefixes for interactive completion/highlighting
		{Name: "ProjectTagPrefix", Pattern: `#`},
		{Name: "ContextTagPrefix", Pattern: `@`},

		// Field setters with colon (op starters) - MUST come before Ident
		// These consume the field name and colon together
		{
			Name:    "SetField",
			Pattern: `\b(title|notes|due|waiting|waiting-for|waiting_for|state|project|projects|context|contexts|meta|id|text)\b\s*:`,
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
			Pattern: `!\s*\b(notes|due|waiting|waiting-for|waiting_for|projects|contexts|meta)\b`,
		},

		// Clear op for non-field cases (just the ! symbol)
		{Name: "ClearOp", Pattern: `!`},

		// Add/Remove ops as standalone (for tag operations)
		{Name: "AddOp", Pattern: `\+`},
		{Name: "RemoveOp", Pattern: `-`},

		// Punctuation
		{Name: "Colon", Pattern: `:`},
		{Name: "Comma", Pattern: `,`},
		{Name: "LParen", Pattern: `\(`},
		{Name: "RParen", Pattern: `\)`},

		// Logical operators
		{Name: "AndOp", Pattern: `&&`},
		{Name: "OrOp", Pattern: `\|\|`},

		// In-progress quoted string support for interactive shell.
		{Name: "QuoteStart", Pattern: `"`, Action: lexer.Push("String")},

		// Identifiers and words (catch-all for regular words including alphanumeric)
		{Name: "Ident", Pattern: `[a-zA-Z0-9_-]+`},

		// Whitespace (elided)
		{Name: "Whitespace", Pattern: `\s+`},
	},
	"String": {
		{Name: "QuoteEnd", Pattern: `"`, Action: lexer.Pop()},
		{Name: "StringEscape", Pattern: `\\.`},
		{Name: "StringText", Pattern: `[^"\\]+`},
		{Name: "StringBackslash", Pattern: `\\`},
	},
})
