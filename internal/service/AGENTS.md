# internal/service AGENTS

## OVERVIEW
Owns app use-cases: validate/normalize inputs, call store, shape outputs.

## WHERE TO LOOK
- `internal/service/interface.go`: service surface used by cmd/ and shell.
- `internal/service/requests.go`: request/response types; pointers = optional fields.
- `internal/service/parse.go`: normalization helpers (state, dates, meta flags).
- `internal/service/task_read.go`: list/get flows; filter handoff via `nlp.FilterExpr`.
- `internal/service/task_write.go`: create/update/done/delete; tag/meta merging rules.
- `internal/service/task_*_test.go`: behavior tests for parsing + store effects.

## CONVENTIONS
- Treat this as the boundary for user input: accept raw strings in requests, normalize immediately.
- Prefer domain helpers for normalization and error construction; keep errors actionable.
- Keep optionality explicit: pointers mean "set", `Clear*` bools mean "unset".
- Normalize collections before store calls: trim space, drop empty, de-dupe, stable order.
- Keep methods context-first; avoid hidden side effects (sync/push only when invoked).
- Do not expose sqlc/generated types; return store models only (e.g. `*store.Task`).

## ANTI-PATTERNS
- Parsing/validation in `cmd/` or `internal/shell/`; move it into service or domain.
- Passing raw state/date/meta through to store without normalization.
- Re-implementing filtering in Go; push filtering to store (and keep `nlp.FilterExpr` intact).
- Read paths that write (e.g. list triggers updates) unless the command name implies it.
- Mutating caller-owned slices/maps during normalization; copy first.

## NOTES
- `UpdateTaskRequest` is a patch; `FullUpdateTaskRequest` is a replace. Do not blur semantics.
- Meta flags come in as `k:v`; separator/format rules live in `internal/domain`.
- Date parsing is `YYYY-MM-DD` to a UTC day; nil means "not set".
- When adding a use-case: add request type, implement on `TaskService`, wire cmd/, add tests.
