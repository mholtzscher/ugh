# AGENTS.md - Agentic Coding Guide for ugh

This document provides essential information for AI coding agents working in this repository.

## Project Overview

**ugh** is a todo.txt-inspired task CLI with libSQL storage (Turso), written in Go 1.25.

- **Module**: `github.com/mholtzscher/ugh`
- **CLI Framework**: spf13/cobra
- **Database**: libSQL (turso.tech/database/tursogo)
- **SQL Generation**: sqlc
- **Migrations**: goose (embedded)
- **Dev Environment**: Nix flake + just task runner

## Build, Test, and Lint Commands

### Using just (recommended)

```bash
just build              # Build binary (development)
just build-release      # Build binary (release, stripped)
just run <args>         # Run with arguments: just run add "Buy milk"

just test               # Run all tests
just test-verbose       # Run tests with verbose output
just test-run NAME      # Run specific test: just test-run TestScripts
just test-pkg PKG       # Run tests for package: just test-pkg internal/store

just fmt                # Format code (go fmt)
just vet                # Static analysis (go vet)
just lint               # Comprehensive linting (golangci-lint)
just check              # Run all checks: fmt, vet, lint, test

just tidy               # Tidy go modules
just sqlc               # Regenerate sqlc code after SQL changes
just gomod2nix          # Update gomod2nix.toml after dependency changes
```

### Using Go directly

```bash
go build                           # Build
go test ./...                      # All tests
go test -v ./...                   # Verbose tests
go test -run TestScripts ./...     # Specific test by name
go test ./internal/store           # Tests for specific package
go fmt ./...                       # Format
go vet ./...                       # Static analysis
golangci-lint run                  # Lint (requires golangci-lint)
```

### Using Nix

```bash
nix build                # Build the package
nix flake check          # Run all checks (go-test, go-lint)
nix develop              # Enter development shell
```

## Project Structure

```
ugh/
├── main.go                    # Entry point (calls cmd.Execute())
├── main_test.go               # Integration tests using testscript
├── cmd/                       # CLI commands (Cobra)
│   ├── root.go                # Root command, global flags, store init
│   ├── add.go, list.go, ...   # Subcommands (one file per command)
│   └── utils.go               # Utility functions
├── internal/
│   ├── store/                 # Database layer
│   │   ├── store.go           # Store struct, Open, CRUD operations
│   │   ├── types.go           # Task and Filters types
│   │   ├── migrations.go      # Embedded migrations filesystem
│   │   ├── migrations/        # SQL migration files (goose format)
│   │   └── sqlc/              # Generated sqlc code (DO NOT EDIT)
│   ├── todotxt/               # todo.txt format parsing/formatting
│   │   ├── parse.go, format.go, types.go
│   └── output/                # Output formatting (JSON, human, plain)
├── db/queries/tasks.sql       # SQL queries for sqlc
└── testdata/script/           # testscript integration tests
```

## Code Style Guidelines

### Imports

Group imports in this order with blank lines between groups:
1. Standard library
2. Internal packages (this module)
3. External packages

```go
import (
    "context"
    "fmt"

    "github.com/mholtzscher/ugh/internal/store"
    "github.com/mholtzscher/ugh/internal/todotxt"

    "github.com/spf13/cobra"
)
```

### Naming Conventions

- **Files**: lowercase, underscore for multi-word (`task_projects.go`)
- **Packages**: short, lowercase, no underscores (`store`, `todotxt`, `output`)
- **Types**: PascalCase (`Task`, `Filters`, `Writer`)
- **Functions/Methods**: PascalCase for exported, camelCase for unexported
- **Variables**: camelCase (`rootOpts`, `addCmd`, `storeTask`)
- **Constants**: PascalCase for exported, camelCase for unexported
- **Acronyms**: Use consistent casing (`ID`, `JSON`, `SQL`, `TTY`)

### Error Handling

- Wrap errors with context using `fmt.Errorf("action: %w", err)`
- Return errors early, avoid deep nesting
- Use `errors.New()` for simple error messages
- Commands return errors; root command prints to stderr

```go
if err := doSomething(); err != nil {
    return fmt.Errorf("do something: %w", err)
}
```

### Cobra Command Pattern

Each command follows this structure:

```go
var cmdOpts struct {
    // Command-specific options
}

var cmdCmd = &cobra.Command{
    Use:     "cmd [args]",
    Aliases: []string{"c"},
    Short:   "Brief description",
    RunE: func(cmd *cobra.Command, args []string) error {
        ctx := context.Background()
        // Implementation
    },
}

func init() {
    cmdCmd.Flags().StringVarP(&cmdOpts.Field, "field", "f", "", "description")
}
```

### Database Pattern

- Use `openStore(ctx)` to get a store instance
- Always `defer st.Close()` after opening
- All times stored as Unix timestamps (int64)
- Dates stored as strings in "2006-01-02" format
- Use `sql.NullString` for nullable string columns

### Types

- Use `*time.Time` for optional date fields
- Use `map[string]string` for metadata
- Use `[]string` for projects, contexts, unknown tokens
- Initialize maps before use: `if task.Meta == nil { task.Meta = map[string]string{} }`

### Testing

Integration tests use `rogpeppe/go-internal/testscript`:
- Test scripts in `testdata/script/*.txt`
- Each script builds the binary and runs CLI commands
- Use `exec` for commands, `stdout` for assertions

```
# Example test script
exec ugh --db $WORK/db.sqlite add Buy milk +groceries
exec ugh --db $WORK/db.sqlite list
stdout 'Buy milk'
```

### Formatting

- Run `just fmt` or `go fmt ./...` before committing
- No manual formatting required - use standard Go formatting
- Line length: no hard limit, but keep reasonable

### Output Modes

The CLI supports three output modes:
1. **JSON** (`--json`): Machine-readable JSON output
2. **Human/TTY**: Formatted tables for interactive use
3. **Plain/Pipe**: todo.txt format for scripting

Always use `outputWriter()` and its methods for consistent output.

### Code Generation

- sqlc generates code in `internal/store/sqlc/` - DO NOT EDIT these files
- After modifying `db/queries/tasks.sql`, run `just sqlc`
- The `//go:generate` directive in `store.go` can also be used

### Dependencies

- After adding, removing, or updating Go dependencies, run `just tidy`
- For Nix builds, also run `just gomod2nix` to update `gomod2nix.toml`
