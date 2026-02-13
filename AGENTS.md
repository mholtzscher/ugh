# ugh AGENTS

Task CLI (SQLite/libSQL) with optional Turso sync.

**Stack**: Go 1.25+, urfave/cli/v3, goose (migrations), sqlc (SQLite), golangci-lint.

**Generated**: 2026-02-13
**Branch**: main
**Commit**: 1126f94

## Hard Rules

- Never commit code unless explicitly prompted by the user.
- After any change: run `just check`.

## Commands

```bash
just check              # generate fmt vet lint test tidy gomod2nix
just run -- <args>      # run local (urfave/cli flags need `--`)
just generate           # sqlc + go:generate (NLP strings)
just test               # go test -p 1 ./...
```

## Structure

```
./
├── main.go                     # entrypoint -> cmd.Execute
├── cmd/                        # CLI commands + wiring
├── internal/
│   ├── shell/                  # interactive REPL (preprocess/parse/compile/dispatch/render)
│   ├── nlp/                    # DSL lexer/parser + diagnostics
│   │   └── compile/            # AST -> service requests (“plan”)
│   ├── service/                # use-cases; thin layer over store
│   ├── store/                  # SQLite/libsql + migrations + sqlc + filter->SQL
│   ├── editor/                 # $EDITOR TOML flow + JSON schema
│   ├── output/                 # JSON/human/table/pipe output
│   ├── daemon/                 # periodic sync daemon (+ service-manager helpers)
│   ├── flags/                  # flag names + validators
│   ├── domain/                 # core task concepts (state/date/meta)
│   └── config/                 # config load/save + paths
├── db/queries/                 # sqlc query inputs
├── testdata/script/            # testscript E2E scripts
└── docs/                       # design notes (daemon doc is large)
```

## Where To Look

| Task | Location | Notes |
|------|----------|-------|
| Add/change CLI subcommand | `cmd/` | flat files; commands are package-level globals (`//nolint:gochecknoglobals`) |
| REPL behavior / pronouns / context | `internal/shell/executor.go` | preprocessing does naive `strings.ReplaceAll` |
| DSL grammar / parse errors | `internal/nlp/` | diagnostics + lexer/parser |
| DSL -> request rules | `internal/nlp/compile/plan.go` | emits `internal/service/*Request` |
| Task CRUD / done/undo / filters | `internal/store/store.go` | big; many `//nolint:*` hotspots |
| Filter expr -> SQL (JSON1) | `internal/store/filter_sql_builder.go` | json_each/json_array_length |
| Service “use cases” | `internal/service/` | bridges domain normalization <-> store |
| E2E CLI tests | `main_test.go` + `testdata/script/*.txt` | testscript harness; public CLI only |

## Conventions (Project-Specific)

- Lint is strict (`.golangci.yml`): goimports `local-prefixes: github.com/mholtzscher/ugh`, golines `max-len: 120`.
- depguard bans: `log` outside `**/main.go` (use `log/slog`); `math/rand` in non-test (use `math/rand/v2`);
  protobuf + uuid package bans (see `.golangci.yml`).
- `//nolint` must name linter + include rationale (`nolintlint`).
- Generated code (do not edit): `internal/store/sqlc/*`, `internal/nlp/*_string.go`.

## Gotchas

- Migration `internal/store/migrations/00009_reset_append_only_schema.sql` drops task tables (data loss).
- Repo may contain build artifacts (`./ugh`, `./result/`); ignore for code navigation.
- Docs mention an HTTP daemon API that does not exist (`docs/daemon-design.md`).

## Subdir AGENTS

- `cmd/AGENTS.md`
- `internal/shell/AGENTS.md`
- `internal/nlp/AGENTS.md`
- `internal/nlp/compile/AGENTS.md`
- `internal/service/AGENTS.md`
- `internal/store/AGENTS.md`
- `internal/editor/AGENTS.md`
- `internal/daemon/AGENTS.md`
- `testdata/script/AGENTS.md`
