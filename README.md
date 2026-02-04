# ugh

A todo.txt-inspired task CLI with SQLite storage.

## Installation

```bash
go install github.com/mholtzscher/ugh@latest
```

Or with Nix:

```bash
nix build
```

## Usage

```bash
# Add tasks (supports todo.txt format)
ugh add "Buy milk +groceries @errands"
ugh add "Call mom" -p A --project family

# List tasks
ugh list                    # pending tasks
ugh list --all              # include completed
ugh list --project groceries
ugh list --context errands
ugh list --priority A

# List available tags
ugh projects
ugh contexts

# Complete tasks
ugh done 1 2 3

# Undo completion
ugh undo 1

# Edit a task
ugh edit 1 --priority B --project work

# Show task details
ugh show 1

# Remove tasks
ugh rm 1 2

# Import/export todo.txt
ugh import todo.txt
ugh export - --all          # stdout
```

## Development

This CLI uses `github.com/urfave/cli/v3`. Flag names are centralized in
`internal/flags/constants.go` so commands can read values with `cmd.String()`,
`cmd.Bool()`, `cmd.Int()`, and friends.

Common tasks:

```bash
just build
just test
just lint
just fmt
```

## Output Formats

- **TTY**: Formatted table output (default)
- **JSON**: `--json` flag for machine-readable output
- **Pipe**: Plain todo.txt format when piped

## Configuration

ugh can be configured with a TOML file. The default location is:

- **Linux**: `~/.config/ugh/config.toml`
- **macOS**: `~/Library/Application Support/ugh/config.toml`
- **Windows**: `%AppData%\ugh\config.toml`

Use the `--config` flag to specify a custom config path.

### DB Path Resolution

The database path is resolved in this order:

1. `--db` flag (highest priority)
2. `db.path` in the config file
3. Default location: `~/.local/share/ugh/ugh.sqlite` (Linux)

Example config file:

```toml
version = 1

[db]
path = "~/.local/share/ugh/ugh.sqlite"
```

Paths support:
- Environment variable expansion (`$HOME`, `$USER`, etc.)
- Home directory expansion (`~/`)
- Relative paths (when set via config file, resolved relative to config file location)

### Sync Settings

When `db.sync_url` is set, you can enable automatic sync after writes:

```toml
[db]
sync_url = "libsql://example.turso.io"
auth_token = "..."
sync_on_write = true
```

## Global Flags

```
--config <path>  Path to config file
--db <path>      Custom database path (overrides config)
--json           Output as JSON
--no-color       Disable color output
```

## todo.txt Format

Tasks follow the [todo.txt](https://github.com/todotxt/todo.txt) format:

```
(A) 2024-01-15 Call mom +family @phone due:2024-01-20
x 2024-01-14 2024-01-10 Buy groceries +shopping
```

- `(A)` - Priority (A-Z)
- `+project` - Project tags
- `@context` - Context tags
- `key:value` - Metadata

## License

MIT
