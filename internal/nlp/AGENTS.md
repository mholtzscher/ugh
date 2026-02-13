# internal/nlp AGENTS

Lexer + participle parser build typed DSL AST, emit diagnostics, then postprocess/normalize into validated commands.

## Where To Look
- `internal/nlp/lexer.go` stateful lexer rules; token priority/order matters.
- `internal/nlp/dsl_parser.go` participle build config (unions, token maps, case/unquote/elide).
- `internal/nlp/ast.go` grammar structs + interfaces (Command/Operation/FilterExpr), parser tags.
- `internal/nlp/dsl_parse.go` custom Parse() hooks for verbs/targets/values (PeekingLexer driven).
- `internal/nlp/dsl_postprocess.go` AST -> normalized command (validation, defaults, expr conversion).
- `internal/nlp/parser.go` public Parse entrypoint + mode gating; wraps parse errors into diagnostics.
- `internal/nlp/types.go` ParseResult + Diagnostic (Severity/Code/Hint) and intent/mode enums.
- `internal/nlp/lex.go` + `internal/nlp/dsl_symbols.go` debug/introspection helpers.
- `internal/nlp/*_test.go` grammar coverage + tokenization/parse behavior.

## Conventions
- Keep grammar (ast.go) declarative; put matching into Parse() methods in `internal/nlp/dsl_parse.go`.
- Lexer tokens: place specific patterns before catch-alls (SetField/AddField/RemoveField before Ident).
- Postprocess is the normalization/validation boundary (defaults, canonicalization, drop no-ops).
- Diagnostics: prefer stable Code strings; keep messages user-facing, hints actionable.
- Tag tokens are mapped to strip prefix (#/@) in `internal/nlp/dsl_parser.go`.

## Anti-Patterns
- Encoding semantics in lexer regexes beyond tokenization.
- Adding new Command/Operation types without updating participle.Union wiring.
- Returning raw participle errors without a Diagnostic.
- Editing generated stringer outputs (`types_string.go`, `ast_string.go`, `ast_field_string.go`).
- Letting postprocess accept ambiguous/empty structures and failing later.

## Notes
- In-progress quotes use a dedicated lexer "String" state; keep compatible with shell UX.
- Filter parsing builds a chain tree then converts to FilterExpr in postprocess.
- Token mapping trims tag prefixes; HashNumber keeps leading # for numeric target parsing.
- Regenerate enums via go:generate in `internal/nlp/ast.go` and `internal/nlp/types.go`.
