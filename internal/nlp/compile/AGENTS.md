# internal/nlp/compile AGENTS

Compile `internal/nlp` AST into typed `internal/service` requests ("plan"), incl filter + date normalization.

## Where To Look
- `internal/nlp/compile/plan.go` Build entry; command switch -> service requests.
- `internal/nlp/compile/plan.go` NormalizeFilterExpr + compileFilterExpr (recursive) + compilePredicate rules.
- `internal/nlp/compile/plan.go` normalizeDate (relative date parsing + canonical YYYY-MM-DD output).
- `internal/nlp/compile/plan.go` buildUpdateRequest (explicit target vs SelectedTaskID).
- `internal/nlp/compile/plan_test.go` filter/date/state/id behavior; add new cases here first.

## Conventions
- Keep compilation explicit per op/predicate kind; prefer straight switches.
- Canonicalize early: TrimSpace, ToLower (dates), domain.NormalizeState (state).
- Date outputs are always `domain.DateLayoutYYYYMMDD`; parse exact layout before relative parsing.
- Thread time via BuildOptions.Now for determinism; avoid time.Now() in helpers.
- Filter compilation returns a new expr tree (do not mutate shared nodes).
- Errors should be user-facing, stable strings when possible (tests assert on them).

## Anti-Patterns
- Let invalid/ambiguous inputs leak into service/store (empty predicate text, non-positive id, non-canonical dates).
- Adding filter semantics here that belong in store (SQL mapping, JSON1 details).
- Wildcards for predicate kinds that do not support it.
- New implicit defaults without going through target resolution.
- Hidden time dependence.

## Notes
- normalizeDate accepts exact `YYYY-MM-DD`, special `next-week`, else go-naturaldate relative to Now.
- Invalid relative dates wrap as `domain.InvalidDateFormatError(value)`.
- compilePredicate rules: state normalize+validate; due requires non-empty+normalized date; id must be >0; project/context/text non-empty.
- Meta/project/context list inputs are comma-split, trimmed, de-duped; empty entries dropped.
