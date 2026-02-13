# internal/shell AGENTS

REPL reads a line, resolves session context/pronouns, parses+compiles to a plan, executes via service, then renders.

## Where To Look
- `internal/shell/repl.go` REPL loop; built-ins; history record; display dispatch.
- `internal/shell/executor.go` preprocessInput (pronouns); injectContext (sticky #project/@context); compile+execute.
- `internal/shell/prompt.go` readline config; prompt (color/noColor); loads history via service.
- `internal/shell/prompt_editor.go` autocomplete + token styling (nlp.Lex).
- `internal/shell/display.go` ShowResult routing; Clear() ANSI sequence.
- `internal/output/*` Writer for line output + task version diffs.

## Conventions
- Pipeline order matters: preprocess -> parse -> diagnostics gate -> injectContext -> compile -> execute.
- SessionState mutations stay localized (LastTaskIDs/SelectedTaskID + sticky context fields).
- Keep noColor content-equivalent; only styling differs.
- ExecuteResult is executor<->display/history contract: always set Intent+Summary; set TaskIDs when meaningful.
- Prefer bounded time sources; avoid multiple time.Now() calls per execution.

## Anti-Patterns
- Naive pronoun/context rewriting via strings.ReplaceAll is not token-aware; boundary bugs.
- Silent error drops outside history recording.
- Executor printing to stdout/stderr; keep IO in display/output layers.
- New recursive Execute() calls; view already re-enters with a synthetic query.
- Display.Clear() in non-interactive flows (ANSI escapes pollute scripts).

## Notes
- Pronoun replacement requires surrounding spaces (" it ", " this "); punctuation/line edges do not match.
- "that" maps to second-to-last ID; if only one ID exists it falls back to last.
- "selected" is rewritten before parsing.
- Sticky context injection applies to Create/Update/Filter only; View/Context/Log bypass it.
- Script mode trims whitespace, skips blanks and '#'-prefixed lines; errors annotated with line number.
