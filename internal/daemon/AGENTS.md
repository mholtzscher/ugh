# internal/daemon AGENTS

Periodic sync daemon plus user-level service-manager install/start/stop helpers (systemd/launchd).

## Where To Look
- `internal/daemon/daemon.go` sync loop, retries/backoff, SIGINT/SIGTERM shutdown, opens DB only during sync.
- `internal/daemon/config.go` defaults + parsing of TOML strings into time.Duration.
- `internal/daemon/service/manager.go` Manager interface, Status/InstallConfig structs, Detect() platform selection.
- `internal/daemon/service/systemd.go` user unit generation, systemctl/journalctl wrappers, status PID lookup.
- `internal/daemon/service/launchd.go` LaunchAgents plist generation, log file paths, launchctl status checks.
- `cmd/daemon_run.go` wiring (store.Options, busy timeout, TURSO_GO_CACHE_DIR, slog setup).
- `internal/config/config.go` daemon keys + db sync_url/auth_token.
- `internal/store/*` store.Open + st.Sync(ctx) is the only daemon DB operation.

## Conventions
- Keep daemon DB access scoped: open -> Sync -> close; no long-lived handles.
- Treat daemon.Config as fully-parsed types; keep TOML parsing at the boundary (ParseConfig).
- Logging: slog text to stderr for foreground; JSON when writing to file; file mode 0600, dir mode 0750.
- Service installs are per-user: systemd under XDG_CONFIG_HOME (or ~/.config), launchd under ~/Library/LaunchAgents.
- Service manager errors: prefer exported sentinels (ErrNotInstalled, ErrNotRunning, ...).

## Anti-Patterns
- Holding the DB open across ticks (locks out interactive CLI).
- Unbounded retries or ignoring ctx cancellation.
- Shelling out via "sh -c" or interpolating user args into ExecStart/launchctl calls.
- Adding config without threading through install templates.
- Allowing PeriodicSync <= 0 to reach time.NewTicker.

## Notes
- Daemon exits early when db.sync_url is empty; it is a no-op without Turso.
- Retry loop is attempt 0..SyncRetryMax (total SyncRetryMax+1 tries); backoff doubles each retry.
- Shutdown does a final syncOnce (no retry loop).
- systemd logs via journald; launchd logs are file-based and tailed via `tail`.
