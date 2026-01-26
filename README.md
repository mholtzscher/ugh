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

## Output Formats

- **TTY**: Formatted table output (default)
- **JSON**: `--json` flag for machine-readable output
- **Pipe**: Plain todo.txt format when piped

## Global Flags

```
--db <path>     Custom database path (default: ~/.config/ugh/ugh.sqlite)
--json          Output as JSON
--no-color      Disable color output
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
