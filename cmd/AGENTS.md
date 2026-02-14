# cmd AGENTS

Urfave/cli command wiring: global flags + config load, output selection, sync wrappers, and store/service setup.

## Where To Look
- `cmd/root.go` root command, global flags, config bootstrap, `openStore`, `outputWriter`.
- `cmd/utils.go` shared helpers: args parsing, `newService`, `maybeSyncBeforeWrite/AfterWrite`.
- `cmd/sync.go` manual sync subcommands; JSON vs human output branching.
- `cmd/help_categories.go` command `Category` and help grouping.
- `cmd/daemon*.go` daemon command tree + install/start/stop/logs/status wiring.
- `cmd/config*.go` config subcommands; flag/value plumbing patterns.

## Conventions (cmd only)
- Commands are package-level `var` `*cli.Command` with `//nolint:gochecknoglobals` (CLI registry).
- Keep `Action` closures tiny; put logic in `runX(ctx, cmd)` helpers for testability/readability.
- Global flags live on root; access via cached `root*` vars (set in `Before: loadConfig`).
- Flag names come from `internal/flags`; avoid literal flag strings in cmd code.
- Output selection: call `outputWriter()`; if `writer.JSON` then `json.Encoder`, else `writer.Write*`.
- For human success/info messages, prefer `writer.WriteSuccess` / `writer.WriteInfo` over raw `WriteLine` when semantics are explicit.
- For mutating commands: `svc, _ := newService(ctx)` then `maybeSyncBeforeWrite` / write / `maybeSyncAfterWrite`.
- Prefer service-layer APIs (`internal/service`) over direct store calls unless command is store-specific (eg sync).

## Anti-Patterns
- Printing directly (`fmt.Println`, `pterm.*`) instead of `outputWriter()`; breaks `--json` and tests.
- Opening DB directly (`store.Open`) from subcommands; bypasses config resolution + lock backoff.
- Re-loading config per subcommand; root `Before` already establishes process-wide state.
- Forgetting `defer st.Close()` on manual store usage (sync/status style commands).
- Hardcoding output shape per command; keep JSON structs explicit and stable.
- Adding command-local flags that duplicate root flags (config/db/json/no-color); keep them global.

## Notes
- `loadConfig` may auto-initialize a default config (and default DB path) when caller did not force `--db`.
- `openStore` sets `TURSO_GO_CACHE_DIR` under the DB dir and retries on locking errors (daemon running).
- Errors from `Execute()` are rendered to stderr via `output.Writer` (not `--json`) to keep failure UX consistent.
