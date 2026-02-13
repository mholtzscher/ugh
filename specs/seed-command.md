# Seed Command (Temp DB)

## Problem
Need fast, reproducible(ish) local datasets for debugging/perf demos: tasks across states with realistic edits + status churn, written to a uniquely-named temp SQLite DB by default so real data is never touched.

## Goals
- New hidden command: `ugh seed`
- Flags: `--seed`, `--count` (deterministic content + mutation choices)
- Default output DB: OS temp dir, unique filename (no clobber)
- Data realism: tasks spread across existing states (`inbox|now|waiting|later|done`) with multiple versions per task via edits + done/undo
- Output: print `db path` + `seed` (JSON under `--json`)

## Non-goals
- No sync/libsql remote usage (always local SQLite file)
- No promise of deterministic timestamps (created/updated/completed vary per run)
- No custom user-supplied distributions/percentages (keep v1 small)

## CLI
Command: `ugh seed`

Flags:
- `--seed <int64>` (default `1`)
- `--count <int>` (default `200`)
- `--churn <int>` (default `5`) number of post-create mutations per task
- `--out <path>` optional explicit db path (overrides temp)
- `--force` only meaningful with `--out`; overwrite existing file

Help visibility:
- Command is hidden from root help + completion via `cli.Command{Hidden: true}`.
- Still callable directly (and `ugh seed --help` works).

## DB Path Rules
- If `--out` empty: create unique path via `os.CreateTemp("", "ugh-seed-<seed>-*.sqlite")`, close file, use returned name.
- If `--out` set:
  - resolve to absolute path
  - ensure parent dir exists (`os.MkdirAll`)
  - if file exists and `--force` not set: error
  - if file exists and `--force` set: remove then proceed

Safety:
- `ugh seed` MUST NOT use `cmd.openStore()` / `cmd.newService()` (those resolve configured DB); it must open the DB at the computed temp/`--out` path.

## Seeding Behavior
### RNG
- Use `math/rand/v2` seeded by `--seed`.
- Same seed+count+churn => same titles/field choices/mutation sequence. (Timestamps differ.)

### Initial task generation
Use existing service layer to ensure invariants + history are correct:
- Open store: `store.Open(ctx, store.Options{Path: outPath})` (no SyncURL/AuthToken)
- Create task service: `service.NewTaskService(st)` (or equivalent constructor used by `cmd.newService`, but without `openStore`).

Generate fields from fixed in-code pools (ASCII):
- Titles: verb + object + optional qualifier (e.g. "Draft quarterly report")
- Notes: short phrases; sometimes empty
- Projects/Contexts: 0-2 each from small fixed lists
- Meta: 0-3 keys from fixed set (e.g. `prio`, `est`, `ticket`, `source`)
- Due date: some tasks get `YYYY-MM-DD` near a fixed anchor date (e.g. 2026-01-01) +/- N days so seed output is stable across calendar time
- Waiting-for: only for some `waiting` tasks (name pool)

Initial state distribution (v1 target, tweakable but stable):
- inbox 35%, now 20%, later 20%, waiting 15%, done 10%

### Churn / version creation
For each task, apply `--churn` mutations, each guaranteed to create a visible diff for `ugh log`.
Mutation types (weighted):
- Edit title (append/change token)
- Edit notes (set/clear)
- State hop among non-done states (inbox/now/waiting/later)
- Set/clear due
- Set/clear waiting-for
- Add/remove one project
- Add/remove one context
- Set/remove one meta key
- Done + optional undo (via `SetDone(true)` then `SetDone(false)`)

No-op prevention:
- If chosen mutation would not change anything, re-pick up to N tries, then fall back to a guaranteed change (e.g. title tweak).

## Output
Human (default): print key-values (stable, one per line):
- `db: <path>`
- `seed: <seed>`
- `count: <count>`
- `churn: <churn>`

JSON (`--json`):
```json
{"dbPath":"...","seed":1,"count":200,"churn":5}
```

## Errors
- Invalid flags: `count < 0`, `churn < 0` => error
- `--out` exists without `--force` => error
- Any DB open/migration error => error

## Acceptance Criteria
- `ugh seed --seed 42 --count 10` prints a temp DB path; file exists.
- `ugh --db <path> list --all` shows tasks across states.
- `ugh --db <path> log 1` shows multiple versions and at least one field change.
- `ugh help` does NOT list `seed`.
- `ugh --db real.sqlite seed` still creates a temp/`--out` DB and does not open/modify `real.sqlite`.

## Test Plan
Add testscript `testdata/script/seed_cmd.txt`:
- `exec ugh seed --out $WORK/db.sqlite --seed 1 --count 3 --churn 2 --force`
- `exec ugh --db $WORK/db.sqlite list --all`
- `exec ugh --db $WORK/db.sqlite log 1`
- `exec ugh help`
- `! stdout '\\bseed\\b'`

Run: `just check`

## Implementation Notes (Files)
- `cmd/seed.go` new hidden command + flags + seed logic
- `cmd/root.go` register `seedCmd`
- `internal/flags/constants.go` add `FlagSeed`, `FlagCount`, `FlagChurn`, `FlagOut`, `FlagForce`
- `testdata/script/seed_cmd.txt` new E2E test

## Effort
M (new command + script; no schema changes)
