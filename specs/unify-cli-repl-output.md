# Unify CLI + REPL Human Output (REPL as Reference)

## Status

- Type: feature plan
- Effort: L
- Status: Ready for task breakdown

## Problem Definition

Human output diverges between non-interactive CLI commands and REPL commands.

- CLI uses table-centric formatting in `internal/output` for TTY and plain TSV for non-TTY.
- REPL often pre-renders string messages in `internal/shell/executor.go` and applies display styling in `internal/shell/display.go`.
- Error formatting differs between CLI and REPL, and parse diagnostics are not rendered consistently.

Goal: make human-facing output behavior the same between CLI and REPL, using current REPL style as the canonical format.

## Scope

In scope:

- Unify default human output (`--json` excluded) across:
  - Non-interactive CLI commands (`ugh <cmd>`)
  - REPL command execution (`ugh shell`)
- Keep TTY-aware behavior:
  - Styled output when writing to TTY and color enabled
  - Plain output when non-TTY or script mode
- Preserve existing JSON output contracts.
- Improve error/diagnostic consistency between CLI and REPL.

Out of scope:

- Auto-pager integration (`less`/`$PAGER`)
- Redesigning JSON schemas
- Changing command semantics or exit code contracts

## User Decisions Captured

- Target unification scope: default human output only.
- Reference style: current REPL output style.
- Compatibility tolerance: minor output breakage in default human mode is acceptable.
- TTY policy: style only when TTY; plain otherwise.
- Must-haves: consistent headers, task formatting, diff-like updates, truncation behavior, and diagnostics quality.
- `add`/`edit` output: one-liner only.
- `show` output: key-value lines (not table).
- Long output: no pager; keep truncation behavior.

## Discovery Summary

Areas explored:

- `internal/output/*` (writer and current human/plain/json paths)
- `internal/shell/*` (executor result formatting + display routing)
- `cmd/*` output usage patterns
- error/diagnostic flow in CLI + REPL + NLP parser

Key findings:

1. `internal/output.Writer` is already shared by CLI and parts of REPL, but REPL still bypasses it often.
2. REPL currently pre-renders list/detail text and styles in executor/display instead of returning typed payloads.
3. CLI output mostly routes through `internal/output`, but command behavior remains table-centric for TTY.
4. REPL script mode can still leak styled rendering through writer TTY assumptions in some paths.
5. Parse diagnostics exist (`nlp.Diagnostic`) but are often dropped because callers only surface `error` string.

## Constraints Inventory

- Keep `--json` output unchanged.
- Keep non-TTY plain output stable unless explicitly required by parity.
- Respect `--no-color` and `NO_COLOR` consistently in CLI and REPL.
- Avoid introducing untyped payload plumbing (`any`) in new result paths.
- Preserve strict lint/test expectations (`just check` required post-change).

## Decision

Adopt a **single human renderer in `internal/output`**, and make both CLI and REPL feed that renderer with structured data.

- Keep `internal/output` as the only output formatting authority.
- Refactor REPL executor/display contract to carry typed payloads (tasks/task/versions + level), not pre-rendered strings for task data.
- Preserve JSON + plain pathways as-is for machine/non-TTY safety.

## Target UX Specification

### 1) Task list output (TTY, human)

For list/filter-like results (CLI + REPL):

- Empty: `No tasks found`
- Non-empty:
  - Header line: `Found N task(s):`
  - Blank line
  - One line per task:
    - `  #<id> <title> [<state>]`
    - optional tags (`#project`, `@context`)
    - optional due date (`YYYY-MM-DD`)

NoColor variant prints same content, no ANSI styling.

### 2) Task detail output (`show`, TTY, human)

Use key-value lines, no table:

- `Task #<id>: <title>`
- Indented fields in fixed order:
  - State
  - Prev State
  - Due
  - Waiting For
  - Projects
  - Contexts
  - Meta
  - Created
  - Updated
  - Completed
  - Notes
- Missing values shown as `-`.

### 3) Create/update output

- Create: `Created task #<id>: <title>`
- Update: `Updated task #<id>: <title>`
- One line only.

### 4) Summary actions (`done`/`undo`/`rm`)

Compact one-line summaries:

- `done: <count> ids=#1,#2`
- `undo: <count> ids=#1,#2`
- `rm: <count> ids=#1,#2`

Severity styling:

- done/undo -> success
- rm -> warning

### 5) Diff/log output

- Task version diff stays field-diff style via `WriteTaskVersionDiff`.
- REPL must stop forcing `TTY: true` when rendering logs; writer mode should reflect actual shell mode/TTY.

### 6) Non-TTY + script behavior

- Non-TTY output remains plain and machine-friendly.
- REPL script modes (`--file`, `--stdin`) force plain output (no ANSI/box/table artifacts).

### 7) Error/diagnostics behavior

- CLI and REPL use one error renderer path.
- Parse diagnostics should include message and hint (when available).
- NoColor/non-TTY error line prefixed as `Error: ...`.

## Implementation Plan (Ordered)

### D1. Output writer unification in `internal/output` (M)

Files:

- `internal/output/human.go`
- `internal/output/output.go`

Work:

- Replace table-based task/list/summary human renderers with REPL-style compact renderers.
- Add key-value detail renderer for single task output.
- Keep JSON + plain branches intact.
- Add `NoColor` behavior to `output.Writer` and prefix printers.

Depends on: none

### D2. CLI command alignment to unified renderer (M)

Files:

- `cmd/add.go`
- `cmd/edit.go`
- `cmd/history.go`
- optional cleanup in `cmd/*` where direct print bypasses writer

Work:

- Ensure add/edit return one-liner style results.
- Replace direct stdout printing in history clear path with writer path.
- Keep JSON behavior unchanged.

Depends on: D1

### D3. REPL structured result contract + display routing (L)

Files:

- `internal/shell/types.go`
- `internal/shell/executor.go`
- `internal/shell/display.go`
- `internal/shell/repl.go`

Work:

- Extend `ExecuteResult` with typed payload fields:
  - `Tasks []*store.Task`
  - `Task *store.Task`
  - `Versions []*store.TaskVersion`
  - `Level ResultLevel` (enum)
- Stop embedding formatted task list output in `Message` for filter/list paths.
- Display should render payloads via `output.Writer` methods, then fallback to `Message` + level.
- Interactive REPL errors should render through the same writer/error path.
- Script mode should force non-TTY rendering.

Depends on: D1

### D4. Unified diagnostics/error rendering (M)

Files:

- `internal/output/output.go` (or new `internal/output/error.go`)
- `internal/nlp/*` (structured diagnostic error type)
- `cmd/root.go`
- `internal/shell/repl.go`
- `cmd/filter_expr.go` (if diagnostics support added there)

Work:

- Add structured diagnostic-aware error renderer (`WriteErr(error)`).
- Keep `WriteError(string)` for compatibility where needed.
- Migrate CLI root fatal path and REPL interactive error path to `WriteErr`.
- Preserve exit code semantics.

Depends on: D1, D3

## Acceptance Criteria

1. `ugh list` (TTY) output shape matches REPL filter/list shape (`Found N...` + compact task lines).
2. `ugh add ...` and REPL create both print one-line created message.
3. `ugh edit ...` and REPL update both print one-line updated message.
4. `ugh show <id>` prints key-value detail lines, no table.
5. `ugh log <id>` and REPL log render through same diff writer path without forced TTY.
6. REPL script modes never emit ANSI clear/style sequences in output.
7. Parse errors and diagnostics render consistently between CLI and REPL.
8. `--json` outputs are unchanged for existing commands.

## Verification Plan

Automated:

- Add/adjust script tests under `testdata/script/` for:
  - list/show/add/edit/done/undo/rm/log outputs (TTY-agnostic assertions where needed)
  - shell script mode plain rendering
  - parse diagnostic formatting
- Run: `just check`

Manual spot checks:

- Interactive `ugh shell` and non-interactive CLI command parity on same operations.
- `NO_COLOR=1` and `--no-color` parity checks.

## Risks + Mitigations

1. Risk: Existing users rely on table layout in TTY CLI.
   - Mitigation: Document release note as UX unification change; keep JSON/plain stable.

2. Risk: REPL refactor may regress history recording or command summaries.
   - Mitigation: Keep `Summary`/`Intent` contract unchanged; add script tests covering history semantics.

3. Risk: Mixed color control across pterm globals and writer-local choices.
   - Mitigation: Define single policy in root/shell init and pass explicit `NoColor` into writer.

## Non-Goals (Explicit)

- No new output format flags.
- No pager integration.
- No daemon/log stream format redesign.

## Rollout Notes

- Ship as one cohesive output-parity change.
- Include changelog note: “Default human CLI output now matches REPL output style.”

## Task Breakdown Readiness

Ready. Deliverables are ordered, scoped, and dependency-linked.

Open questions:

- None.
