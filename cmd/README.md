# Command Architecture

This CLI uses a registry-driven command tree.

## Overview

- `cmd/root.go` bootstraps global flags, loads config, and builds the tree.
- `cmd/registry` validates and assembles command specs into `*cli.Command` nodes.
- Feature modules register their specs via `Register(...)`:
  - `cmd/tasks`
  - `cmd/lists`
  - `cmd/backup`
  - `cmd/tags`
  - `cmd/sync`
  - `cmd/config`
  - `cmd/daemon`

## Runtime Services

Shared runtime behavior lives in `internal/runtime`:

- global option state (`--config`, `--db`, `--json`, `--no-color`)
- config loading/initialization
- db path resolution and store opening
- output writer creation
- sync-on-write helpers

`cmd/root.go` injects these runtime functions through each module's `Deps`.

## Registry Contract

Each module registers `registry.Spec` values:

- `ID` (`registry.ID`): unique command spec identifier
- `ParentID`: parent spec id for nested commands
- `Source`: module identifier for diagnostics
- `Build`: command constructor

The registry validates:

- duplicate IDs
- missing parents
- cycles
- empty command names
- duplicate/overlapping sibling tokens (name/aliases, case-insensitive)
- mixed child declaration styles (preconfigured children + registry children)

## Adding a New Command

1. Add command constructor(s) in the feature module.
2. Add typed `registry.ID` constant(s) in that module.
3. Register specs in `Register(...)` with a `Source` value.
4. If needed, add dependencies to that module's `Deps`.
5. Run:

```bash
just fmt
just vet
just lint
just test
```

## Guardrails

- `cmd/command_surface_test.go` snapshots command path/aliases/categories.
- `cmd/registry/registry_test.go` validates registry invariants.
