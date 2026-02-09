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

	// Operations
	{Name: "AddOp", Pattern: `\+`},
	{Name: "RemoveOp", Pattern: `-`},
	{Name: "ClearOp", Pattern: `!`},

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

	// Identifiers and words (catch-all for regular words including alphanumeric)
	{Name: "Ident", Pattern: `[a-zA-Z0-9_-]+`},

	// Whitespace (elided)
	{Name: "Whitespace", Pattern: `\s+`},
})
