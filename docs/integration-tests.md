# Integration Test Coverage Matrix

This document maps CLI behavior to integration scripts under `testdata/script/`.

## How To Run

```bash
go test -run TestScripts ./...
```

Full repository checks:

```bash
just check
```

## Coverage Matrix

### Core task flows

- `add`, `list` basics: `testdata/script/basic.txt`
- `show` and `edit` flags: `testdata/script/show_edit.txt`
- `done`, `undo`, `rm`: `testdata/script/done_undo_rm.txt`
- global error handling and argument validation: `testdata/script/errors.txt`
- JSON output surface: `testdata/script/json.txt`

### List semantics

- list filters (`--all|--done|--todo`, project/context/search/state/where): `testdata/script/list_filters.txt`
- built-in list commands and aliases (`inbox|now|waiting|later|calendar`): `testdata/script/builtin_lists.txt`
- deterministic ordering semantics for `list --all`: `testdata/script/builtin_lists.txt`

### Projects and contexts

- baseline project/context listing: `testdata/script/projects_contexts.txt`
- `--counts` with `--all|--done|--todo`: `testdata/script/projects_contexts_counts.txt`

### Config and path resolution

- config file loading and `--db` override behavior: `testdata/script/config.txt`
- config command set/get/show basics: `testdata/script/config_cmd.txt`
- config unset and not-set behavior: `testdata/script/config_unset.txt`
- default path auto-init without flags: `testdata/script/default_paths_no_flags.txt`

### Sync

- config-level sync settings persistence: `testdata/script/sync.txt`
- deterministic sync failure paths (offline-safe): `testdata/script/sync_errors.txt`

### Shell/REPL and history

- REPL filter semantics: `testdata/script/repl_filters.txt`
- REPL filter errors: `testdata/script/repl_filters_errors.txt`
- history command filtering and clear flow: `testdata/script/history_cmd.txt`

### Testscript policy

- testscript scenarios should stay end-to-end and exercise behavior through the public CLI only.
- avoid direct SQL assertions or custom testscript DB commands in `testdata/script/`.

## Out Of Scope

- daemon service lifecycle integration (`daemon install/start/stop/restart/logs`) is intentionally excluded from scripts because it depends on host service managers (`launchd`/`systemd`) and is not hermetic.

## Contributor Guidance

- when adding a new top-level command or flag, add or update a script in `testdata/script/` and update this matrix in the same PR.
- prefer hybrid assertions:
  - exact `cmp` for stable pipe output
  - targeted regex for variable values (timestamps, absolute paths)
- avoid order assertions where timestamps can tie; construct deterministic fixtures.

## Known Gaps / TODO

- editor-interactive flow (`ugh edit <id>` without field flags) is not covered end-to-end due non-interactive test environment constraints.
- DB schema/migration invariants are intentionally not asserted in testscript e2e coverage.
