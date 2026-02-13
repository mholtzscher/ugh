# Recent Tasks (List + REPL/DSL) — Implementation Spec

Status: Ready for task breakdown
Effort: M
Approved by: user
Date: 2026-02-13

## Problem

Users want a fast way to see tasks they just created or changed.

## Goals

- CLI supports a recent view under `list`.
- REPL/DSL supports `recent` / `recent:N` with examples using `list` (but `find` remains valid since verbs are synonyms).
- Default shows 20, and excludes done tasks by default.
- No behavior change to default `ugh list` unless flags/keyword used.

## Non-Goals

- Time windows (`--since/--until`, `recent:24h`).
- New top-level command `ugh recent`.
- General sorting framework (`--sort=...`).

## Definitions

- "Recent": sort by task `updated_at` descending (covers both creates + edits; new tasks start with `updated_at == created_at`).
- Default done handling stays as-is: done tasks excluded unless explicitly included via existing flags/filters.

## UX / Interface

### CLI

- New flags on `ugh list`:
  - `--recent`: enable recent ordering; if no explicit limit, default limit to 20.
  - `--limit N`: maximum rows to return.

Rules:

- If `--recent` set and `--limit` not provided: effective limit = 20.
- If `--limit` provided: must be `> 0`.
- `--recent` does not change filtering; it changes ordering + default limit.
- Existing `--all/--done/--todo` semantics unchanged.

Examples:

- `ugh list --recent`
- `ugh list --recent --limit 5`
- `ugh list --recent --where "project:work"`

### REPL/DSL

Add reserved modifier keyword:

- `recent` -> enable recent mode with default limit 20.
- `recent:N` -> enable recent mode with limit N (N > 0).

Prefer docs/examples with `list` verb:

- `list recent:10 and project:work`

But also accept `find recent:10 ...` since `find/list/show/filter` are synonyms in the DSL.

Text search for the literal word "recent" must use `text:recent` (because bare `recent` is reserved).

Invalid:

- `not recent` (modifier cannot be negated)
- `recent:0`, `recent:-1`, `recent:abc`

## Technical Approach

### Data sources

- Use `tasks_current.updated_at` for recency.
- No schema changes.

### Representation

Introduce a new NLP predicate kind `PredRecent`.

- Parse sources:
  - Bare token `recent`
  - Token `recent:<digits>` (treated as a single filter text token today; convert in postprocess)

### Request/Options plumbing

Extend request/options types to carry recent/limit without encoding them into SQL filter predicates.

- `internal/service.ListTasksRequest`:
  - add `Recent bool`
  - add `Limit int64` (0 means unset)

- `internal/store.ListTasksByExprOptions`:
  - add `Recent bool`
  - add `Limit int64` (0 means no LIMIT)

### Modifier extraction (service layer)

Implement a helper that removes `PredRecent` nodes from the filter expression and returns:

- `exprWithoutRecent nlp.FilterExpr` (nil allowed)
- `recentEnabled bool`
- `recentLimit int64` (0 if none)

Semantics when removing from boolean trees:

- `A and recent` -> `A`
- `A or recent` -> `A`
- `recent and recent:10` -> error if conflicting limits
- `not recent` -> error

Rationale: `recent` is a list modifier, not a filtering predicate. Treating it as a boolean literal would make `recent or <expr>` match everything.

Defaulting:

- `recent := req.Recent || extractedRecent`
- `limit := req.Limit`
- if `limit == 0` and extractedLimit > 0: use extractedLimit
- if `limit == 0` and recent: set limit to 20

Done-exclusion defaulting uses the stripped expression.

### Store query changes

In `Store.ListTasksByExpr`:

- If `opts.Recent`:
  - `ORDER BY t.updated_at DESC, t.version_id DESC`
- Else:
  - keep existing ordering (done last, due grouping, then `updated_at DESC`)
- If `opts.Limit > 0`: apply `LIMIT`.

## Acceptance Criteria

- [ ] `ugh list` output/order unchanged when `--recent` not set and no DSL recent modifier used.
- [ ] `ugh list --recent` returns at most 20 tasks, ordered by `updated_at DESC` (ties broken by `version_id DESC`).
- [ ] `ugh list --recent --limit 5` returns at most 5 tasks.
- [ ] `--where "recent:10 and project:work"` works and applies recent ordering + limit.
- [ ] REPL accepts `list recent:10` (documented) and also accepts `find recent:10`.
- [ ] `text:recent` performs a normal text search for the word "recent".
- [ ] `not recent` fails with a clear, stable error.
- [ ] Default done exclusion preserved (recent alone still hides done).

## Test Plan

Unit

- `internal/nlp/parser_test.go`: parse `recent` and `recent:10` -> `PredRecent`.
- `internal/nlp/compile/plan_test.go`: validate `recent` limit parsing (reject <=0 / non-int).
- `internal/service/*_test.go`: modifier stripping + boolean simplification + defaulting.

E2E (testscript)

- Add `testdata/script/recent_list.txt`:
  - create tasks, mutate one, assert ordering and limiting with `ugh list --recent --limit 1`
  - run REPL script with `list recent:1` and assert only one row.

## Files Likely Touched

- `cmd/list.go`
- `internal/flags/constants.go`
- `internal/service/requests.go`
- `internal/service/task_read.go`
- `internal/store/types.go`
- `internal/store/store.go`
- `internal/nlp/ast.go`
- `internal/nlp/dsl_postprocess.go`
- `internal/nlp/compile/plan.go`
- tests: `internal/nlp/*_test.go`, `internal/service/*_test.go`, `testdata/script/recent_list.txt`
