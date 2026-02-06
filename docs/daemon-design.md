# Daemon Design: HTTP API for Turso Sync

This document describes the design for a background daemon (`ughd`) that provides:
1. **Background sync** to Turso cloud (debounced, periodic)
2. **HTTP API** for external consumers (Raycast extension, scripts, etc.)

## Motivation

The CLI currently supports Turso sync via `sync_on_write`, but this adds network latency to every write operation. A daemon enables:
- **Instant CLI commands** - writes go to local SQLite only
- **Background sync** - daemon handles push/pull asynchronously
- **HTTP API** - enables rich integrations (Raycast, Alfred, scripts)

## Architecture

```
                    ┌──────────────────────────────────────────────────┐
                    │                   ughd daemon                    │
                    │                                                  │
                    │  ┌────────────┐  ┌────────────┐  ┌───────────┐  │
   Raycast ────────▶│  │   HTTP     │  │  Watcher   │  │   Sync    │  │
   Extension        │  │  Server    │  │ (fsnotify) │  │  Manager  │  │
   (fetch)          │  │ :9847      │  │            │  │           │  │
                    │  └─────┬──────┘  └─────┬──────┘  └─────┬─────┘  │
                    │        │               │               │        │
                    │        └───────────────┴───────────────┘        │
                    │                        │                        │
                    │                 ┌──────▼──────┐                  │
                    │                 │    Store    │                  │
                    │                 └──────┬──────┘                  │
                    └────────────────────────┼────────────────────────┘
                                             │
      ugh CLI ───────────────────────────────┤
      (direct SQLite access)                 │
                                             ▼
                                    ┌────────────────┐
                                    │  Local SQLite  │
                                    │   (libSQL)     │
                                    └────────┬───────┘
                                             │
                                             │ (background sync)
                                             ▼
                                    ┌────────────────┐
                                    │  Turso Cloud   │
                                    └────────────────┘
```

### Service Manager Integration

```
User runs:              ugh daemon start
                              │
                              ▼
                    ┌─────────────────┐
                    │  Detect OS and  │
                    │  service manager│
                    └────────┬────────┘
                             │
              ┌──────────────┴──────────────┐
              ▼                             ▼
    ┌─────────────────┐           ┌─────────────────┐
    │     Linux       │           │      macOS      │
    │  systemctl      │           │   launchctl     │
    │  --user start   │           │     start       │
    │  ughd.service   │           │ com.ugh.daemon  │
    └────────┬────────┘           └────────┬────────┘
             │                             │
             └──────────────┬──────────────┘
                            ▼
                   ┌─────────────────┐
                   │  Service runs:  │
                   │ ugh daemon run  │
                   └─────────────────┘
```

## Key Design Decisions

1. **CLI stays direct** - CLI commands always use local SQLite directly. The daemon is purely for HTTP API and background sync.

2. **System service management** - Daemon lifecycle managed via systemd (Linux) or launchd (macOS). The `ugh daemon` commands wrap the native service managers for a consistent UX.

3. **User services only** - Services run at user level (`systemctl --user`, `~/Library/LaunchAgents`). No sudo required, simpler permissions model.

4. **net/http stdlib** - No external HTTP framework dependencies. Standard library is sufficient for this use case.

5. **Localhost only** - HTTP server binds to `127.0.0.1` only for security.

6. **Daemon manages logging** - Logging handled by the daemon process itself via config, ensuring consistent behavior across platforms.

7. **Last-write-wins sync** - Turso/libSQL uses LWW conflict resolution. No merge or conflict markers - last push wins.

## Platform Support

| Platform | Status | Service Manager |
|----------|--------|-----------------|
| Linux (systemd) | Supported | `systemctl --user` |
| macOS | Supported | `launchctl` (LaunchAgents) |
| Linux (non-systemd) | Manual setup | Error with instructions |
| Windows | Unsupported | Future consideration |

## File Structure

```
internal/api/
├── server.go           # HTTP server setup, graceful shutdown
├── routes.go           # Route registration
├── middleware.go       # Logging, recovery, content-type, max-body
├── handlers_health.go  # GET /health
├── handlers_tasks.go   # GET/POST/PUT/PATCH/DELETE /tasks
├── handlers_sync.go    # GET/POST /sync/*
├── handlers_meta.go    # GET /projects, /contexts
└── response.go         # JSON response helpers, error formatting

internal/daemon/
├── daemon.go           # Daemon struct, Run(), Shutdown()
├── config.go           # DaemonConfig struct
├── watcher.go          # fsnotify DB file watcher
├── sync_manager.go     # Debounced + periodic sync logic
└── service/
    ├── manager.go      # ServiceManager interface, Detect()
    ├── systemd.go      # Linux systemd implementation
    └── launchd.go      # macOS launchd implementation

cmd/
├── root.go             # Root command, Execute(), global flags
├── add.go, list.go ... # Simple commands (no subcommands)
├── config/             # Config subcommand group
│   ├── config.go       # Parent command + Register()
│   ├── show.go         # ugh config show
│   ├── get.go          # ugh config get
│   └── set.go          # ugh config set
└── daemon/             # Daemon subcommand group
    ├── daemon.go       # Parent command + Register()
    ├── install.go      # ugh daemon install
    ├── uninstall.go    # ugh daemon uninstall
    ├── start.go        # ugh daemon start
    ├── stop.go         # ugh daemon stop
    ├── restart.go      # ugh daemon restart
    ├── status.go       # ugh daemon status
    ├── logs.go         # ugh daemon logs
    └── run.go          # ugh daemon run (foreground, called by service)
```

Subcommand packages expose a `Register(parent *cobra.Command)` function that the root
command calls to wire up the command tree.

The `api` package is independent of the daemon - it only depends on `store` and `service`.
The `daemon` package orchestrates the API server, file watcher, and sync manager.

## Configuration

Add `[daemon]` section to `~/.config/ugh/config.toml`:

```toml
[db]
path = "~/.local/share/ugh/ugh.sqlite"
sync_url = "libsql://your-db.turso.io"
auth_token = "your-token"
sync_on_write = false  # Disable when using daemon

[daemon]
listen = "127.0.0.1:9847"     # HTTP server address
sync_delay = "2s"             # Debounce: wait after last write before sync
periodic_sync = "5m"          # Background sync interval
log_file = ""                 # Empty = stderr, or path to log file
log_level = "info"            # debug, info, warn, error
shutdown_timeout = "30s"      # Max time for graceful shutdown
sync_retry_max = 3            # Max sync retry attempts
sync_retry_backoff = "1s"     # Initial retry backoff (doubles each retry)
```

## CLI Commands

### `ugh daemon install`

Generate and enable the system service.

```bash
ugh daemon install
```

**Behavior:**
1. Detect platform (Linux/macOS)
2. Generate service file with correct paths (binary, config)
3. Write to user service directory:
   - Linux: `~/.config/systemd/user/ughd.service`
   - macOS: `~/Library/LaunchAgents/com.ugh.daemon.plist`
4. Enable the service (systemctl enable / launchctl load)

**Output:**
```
Service installed at ~/.config/systemd/user/ughd.service
Run 'ugh daemon start' to start the daemon
```

### `ugh daemon uninstall`

Disable and remove the system service.

```bash
ugh daemon uninstall
```

**Behavior:**
1. Stop service if running
2. Disable service (systemctl disable / launchctl unload)
3. Remove service file

### `ugh daemon start`

Start the daemon via system service manager.

```bash
ugh daemon start
```

**Behavior:**
1. Check if service is installed
2. Start via `systemctl --user start ughd` or `launchctl start com.ugh.daemon`
3. Verify daemon is responding (poll `/health`)

### `ugh daemon stop`

Stop the running daemon.

```bash
ugh daemon stop
```

**Behavior:**
1. Stop via `systemctl --user stop ughd` or `launchctl stop com.ugh.daemon`

### `ugh daemon restart`

Restart the daemon.

```bash
ugh daemon restart
```

### `ugh daemon status`

Show daemon status.

```bash
ugh daemon status
```

**Output (human):**
```
Service:         installed
Status:          running
PID:             12345
Uptime:          2h 15m 30s
Listen:          127.0.0.1:9847
Last sync:       30s ago
Pending changes: 0
```

**Output (JSON):**
```json
{
  "installed": true,
  "running": true,
  "pid": 12345,
  "uptime_seconds": 8130,
  "listen": "127.0.0.1:9847",
  "last_sync_seconds_ago": 30,
  "pending_changes": 0
}
```

### `ugh daemon logs`

Tail daemon logs.

```bash
# Follow logs
ugh daemon logs

# Show last 100 lines without following
ugh daemon logs --lines 100 --no-follow
```

**Behavior:**
- Linux: `journalctl --user -u ughd -f`
- macOS: `tail -f <log_file>` (from config)

### `ugh daemon run`

Run the daemon server in foreground. This is what the service file executes.

```bash
ugh daemon run
```

**Behavior:**
1. Open store with sync configuration
2. Pull from remote to get latest state
3. Start HTTP server
4. Start file watcher (on DB and WAL files)
5. Start periodic sync ticker
6. Wait for shutdown signal (SIGTERM/SIGINT)
7. Graceful shutdown sequence

Users typically don't run this directly - it's called by the service manager.
Useful for debugging.

## HTTP API

Base URL: `http://127.0.0.1:9847`

All responses are JSON. Errors return appropriate HTTP status codes with:
```json
{
  "error": "error message"
}
```

### Health

#### `GET /health`

Health check endpoint.

**Response:**
```json
{
  "status": "ok",
  "uptime_seconds": 8130,
  "sync": {
    "configured": true,
    "last_sync_seconds_ago": 30,
    "last_error": null,
    "consecutive_failures": 0
  }
}
```

Status values:
- `"ok"` - Everything working
- `"degraded"` - Running but sync failing
- `"error"` - Critical error

### Tasks

#### `GET /tasks`

List tasks with optional filters.

**Query Parameters:**
| Param | Type | Description |
|-------|------|-------------|
| `state` | string | Filter by state (`inbox`, `now`, `waiting`, `later`, `done`) |
| `project` | string | Filter by project name |
| `context` | string | Filter by context name |
| `search` | string | Search in description, projects, contexts, meta |

**Response:**
```json
{
  "tasks": [
    {
      "id": 1,
      "state": "inbox",
      "completion_date": null,
      "creation_date": "2026-01-27",
      "description": "Buy milk",
      "projects": ["groceries"],
      "contexts": ["store"],
      "meta": {"due": "2026-01-30"},
      "created_at": "2026-01-27T10:00:00Z",
      "updated_at": "2026-01-27T10:00:00Z"
    }
  ]
}
```

#### `POST /tasks`

Create a new task.

**Request Body:**
```json
{
  "title": "Buy milk",
  "state": "now",
  "dueOn": "2026-01-30",
  "projects": ["groceries"],
  "contexts": ["store"],
  "meta": {"source": "voice"}
}
```

**Response:** `201 Created`
```json
{
  "task": { /* full task object */ }
}
```

#### `GET /tasks/:id`

Get a single task by ID.

**Response:**
```json
{
  "task": { /* full task object */ }
}
```

**Errors:**
- `404 Not Found` - Task doesn't exist

#### `PUT /tasks/:id`

Replace a task (full update).

**Request Body:**
```json
{
  "description": "Buy milk and eggs",
  "projects": ["groceries"],
  "contexts": ["store"],
  "meta": {"due": "2026-01-30"}
}
```

**Response:**
```json
{
  "task": { /* updated task object */ }
}
```

#### `PATCH /tasks/:id`

Partial update (only specified fields).

**Request Body:**
```json
{
}
```

**Response:**
```json
{
  "task": { /* updated task object */ }
}
```

#### `DELETE /tasks/:id`

Delete a task.

**Response:** `204 No Content`

### Bulk Operations

#### `POST /tasks/done`

Mark multiple tasks as done (`state=done`).

**Request Body:**
```json
{
  "ids": [1, 2, 3]
}
```

**Response:**
```json
{
  "updated": 3
}
```

#### `POST /tasks/undone`

Reopen multiple tasks (restores `prev_state`, clears completion timestamp).

**Request Body:**
```json
{
  "ids": [1, 2, 3]
}
```

**Response:**
```json
{
  "updated": 3
}
```

#### `DELETE /tasks/bulk`

Delete multiple tasks.

**Request Body:**
```json
{
  "ids": [1, 2, 3]
}
```

**Response:**
```json
{
  "deleted": 3
}
```

### Metadata

#### `GET /projects`

List all projects with task counts.

**Query Parameters:**
| Param | Type | Description |
|-------|------|-------------|
| `state` | string | Count only tasks in a given state (`inbox`, `now`, `waiting`, `later`, `done`) |

**Response:**
```json
{
  "projects": [
    {"name": "groceries", "count": 5},
    {"name": "work", "count": 12}
  ]
}
```

#### `GET /contexts`

List all contexts with task counts.

**Query Parameters:**
| Param | Type | Description |
|-------|------|-------------|
| `state` | string | Count only tasks in a given state (`inbox`, `now`, `waiting`, `later`, `done`) |

**Response:**
```json
{
  "contexts": [
    {"name": "store", "count": 3},
    {"name": "home", "count": 8}
  ]
}
```

### Sync

#### `GET /sync/status`

Get sync status.

**Response:**
```json
{
  "configured": true,
  "last_pull_time": 1706360000,
  "last_push_time": 1706359800,
  "pending_changes": 0,
  "network_sent_bytes": 1024,
  "network_received_bytes": 2048,
  "revision": "abc123"
}
```

#### `POST /sync`

Trigger full sync (pull + push).

**Response:**
```json
{
  "action": "sync",
  "message": "synced with remote"
}
```

#### `POST /sync/pull`

Pull changes from remote.

**Response:**
```json
{
  "action": "pull",
  "message": "pulled changes from remote"
}
```

#### `POST /sync/push`

Push local changes to remote.

**Response:**
```json
{
  "action": "push",
  "message": "pushed changes to remote"
}
```

## Implementation Details

### Service Manager Interface

```go
type ServiceStatus struct {
    Installed   bool
    Running     bool
    PID         int
    ServicePath string
}

type ServiceManager interface {
    Name() string
    Install(cfg InstallConfig) error
    Uninstall() error
    Start() error
    Stop() error
    Status() (ServiceStatus, error)
    LogPath() string
    TailLogs(ctx context.Context, follow bool, lines int, w io.Writer) error
}

type InstallConfig struct {
    BinaryPath string  // Absolute path to ugh binary
    ConfigPath string  // Path to config.toml
}
```

### Database Concurrency

Both CLI and daemon access the same SQLite database:

- **WAL mode** enabled for better concurrent access
- **busy_timeout=5000ms** for CLI, consider longer timeout for daemon (30s)
- **Sync mutex** in daemon prevents concurrent Pull/Push operations
- CLI writes go directly to DB; daemon detects via fsnotify

The fsnotify watcher monitors both the main `.sqlite` file and the `-wal` file to catch all writes.

### Sync Conflict Resolution

Turso/libSQL uses **Last-Write-Wins (LWW)** at the transaction level:
- Whoever pushes last "wins" for overlapping rows
- No merge, no conflict markers
- Frequent sync reduces the conflict window

**Expected behavior:** If Machine A and Machine B both edit the same task, the last one to push wins. The other machine's changes are silently overwritten on next pull.

### Graceful Shutdown Sequence

When receiving SIGTERM/SIGINT:

1. Stop file watcher (prevent new sync triggers)
2. Stop accepting new HTTP connections
3. Cancel pending debounce timer
4. Wait for in-flight HTTP requests (with timeout)
5. Wait for in-progress sync to complete
6. Attempt final push of any local changes
7. Close store (flushes WAL)

Timeout controlled by `daemon.shutdown_timeout` config.

### Error Recovery & Retry

Sync failures are retried with exponential backoff:

1. Initial failure: wait `sync_retry_backoff` (default 1s)
2. Second failure: wait 2s
3. Third failure: wait 4s
4. After `sync_retry_max` failures: stop retrying until next trigger

Failed sync status exposed via `/health` endpoint for monitoring.

### HTTP Server Configuration

```go
srv := &http.Server{
    Addr:           cfg.Listen,
    ReadTimeout:    5 * time.Second,
    WriteTimeout:   30 * time.Second,  // Longer for sync operations
    IdleTimeout:    60 * time.Second,
    MaxHeaderBytes: 1 << 16,  // 64KB
}
```

Request body size limited to 1MB via middleware.

## Implementation Phases

### Phase 1: Service Management

- [ ] `internal/daemon/service/manager.go` - Interface definition, Detect()
- [ ] `internal/daemon/service/systemd.go` - Linux implementation
- [ ] `internal/daemon/service/launchd.go` - macOS implementation
- [ ] `cmd/daemon/daemon.go` - Parent command + Register()
- [ ] `cmd/daemon/install.go` - Install command
- [ ] `cmd/daemon/uninstall.go` - Uninstall command
- [ ] `cmd/daemon/start.go` - Start command
- [ ] `cmd/daemon/stop.go` - Stop command
- [ ] `cmd/daemon/restart.go` - Restart command
- [ ] `cmd/daemon/status.go` - Status command
- [ ] `cmd/daemon/logs.go` - Logs command

### Phase 2: HTTP API Package

- [ ] `internal/api/server.go` - HTTP server with timeouts, graceful shutdown
- [ ] `internal/api/middleware.go` - Logging, recovery, content-type, max-body
- [ ] `internal/api/routes.go` - Route registration
- [ ] `internal/api/response.go` - JSON response helpers, error formatting
- [ ] `internal/api/handlers_health.go` - `GET /health`

### Phase 3: Task Handlers

- [ ] `internal/api/handlers_tasks.go`:
  - [ ] `GET /tasks` - List with filters
  - [ ] `POST /tasks` - Create
  - [ ] `GET /tasks/:id` - Get single
  - [ ] `PUT /tasks/:id` - Full update
  - [ ] `PATCH /tasks/:id` - Partial update
  - [ ] `DELETE /tasks/:id` - Delete
  - [ ] `POST /tasks/done` - Bulk done
  - [ ] `POST /tasks/undone` - Bulk undone
  - [ ] `DELETE /tasks/bulk` - Bulk delete

### Phase 4: Metadata & Sync Handlers

- [ ] `internal/api/handlers_meta.go`:
  - [ ] `GET /projects`
  - [ ] `GET /contexts`
- [ ] `internal/api/handlers_sync.go`:
  - [ ] `GET /sync/status`
  - [ ] `POST /sync`
  - [ ] `POST /sync/pull`
  - [ ] `POST /sync/push`

### Phase 5: Daemon Core + Background Sync

- [ ] `internal/daemon/config.go` - DaemonConfig parsing
- [ ] `internal/config/config.go` - Add `[daemon]` section
- [ ] `internal/daemon/daemon.go` - Daemon struct, Run(), Shutdown()
- [ ] `internal/daemon/watcher.go` - fsnotify for DB + WAL files
- [ ] `internal/daemon/sync_manager.go`:
  - [ ] Debounced sync after writes
  - [ ] Periodic sync ticker
  - [ ] Retry with exponential backoff
  - [ ] Sync on startup (pull)
- [ ] `cmd/daemon/run.go` - Foreground server command

### Phase 6: Polish & Documentation

- [ ] Structured logging with `log/slog` (JSON to file)
- [ ] Graceful shutdown sequence
- [ ] Tests for API handlers (`internal/api/*_test.go`)
- [ ] Tests for service managers (mock systemctl/launchctl)
- [ ] Tests for daemon (watcher, sync manager)
- [ ] Update README with daemon documentation
- [ ] Document manual setup for non-systemd Linux

## Security Considerations

1. **Localhost binding** - Server binds to `127.0.0.1` only, not accessible from network
2. **No authentication by default** - Acceptable for localhost-only service
3. **Optional auth token** - Future enhancement: `daemon.auth_token` config for paranoid users
4. **Service permissions** - User services run as the user, no elevated privileges

## Testing Strategy

### Unit Tests

- HTTP handler tests with `httptest` and mock store
- Sync manager tests with mock timers
- Service manager tests with mock exec

### Integration Tests

- Full daemon startup/shutdown with real DB
- HTTP API tests against running daemon
- Concurrent CLI + daemon access tests

### Platform Testing

- CI matrix: Linux (systemd), macOS
- Manual testing for service install/uninstall

## Future Enhancements

- **Unix socket option** - `daemon.socket = "/tmp/ughd.sock"` as alternative to TCP
- **WebSocket support** - Real-time task updates for live UIs
- **Authentication** - Bearer token auth for multi-user scenarios
- **Metrics endpoint** - Prometheus-compatible `/metrics`
- **Windows support** - Windows Service via NSSM or native API
- **Non-systemd Linux** - OpenRC, runit support
