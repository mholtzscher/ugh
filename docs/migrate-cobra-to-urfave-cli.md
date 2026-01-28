# Migration: Cobra + Viper -> urfave/cli

This repo currently uses Cobra for CLI structure and Viper for config/env handling. This document outlines the steps to migrate the CLI to `github.com/urfave/cli/v2` while keeping:

- Flag overrides (highest priority)
- Environment variable overrides
- Config file defaults

and dropping:

- Live reload / config watching (daemon will require restart to pick up config changes)

## Decisions (Locked)

- Config watching/live reload: remove entirely
- Config path validation constraints: remove (config can live anywhere)
- Help output formatting: accept urfave/cli formatting
- Short flags: review and simplify (not necessarily 1:1 compatible)
- Daemon: keep `ugh daemon <subcommand>` structure

## Target Behavior

### Config precedence

For any setting backed by a global flag/env/config key (ex: DB path):

1. CLI flag (ex: `--db`)
2. Env var (ex: `UGH_DB`)
3. Config file (ex: `db.path`)
4. Default

Note: in urfave/cli, env vars are associated with flags. If a value is needed from config file, load it in `Before` and apply to a config struct; if a flag is set, it should override.

### Config skip for specific commands

We skip config loading for commands that create or modify config when the config may not exist yet.

Approach: command metadata.

- Commands that should skip config load (at minimum): `ugh config init`, `ugh config set`
- Add: `Metadata: map[string]interface{}{"skipConfig": true}` on those commands
- Root `Before` hook checks `c.Command.Metadata["skipConfig"]`

## Step-by-Step Migration

### 1) Add urfave/cli dependency and remove Cobra/Viper

1. Add dependency:

   - `go get github.com/urfave/cli/v2`

2. Remove Cobra and Viper dependencies:

   - Delete imports (code changes will come in later phases)
   - Run `go mod tidy`

3. If you use Nix in this repo:

   - Run `just gomod2nix`

### 2) Replace Cobra root with a cli.App

Files:

- `main.go` (replace `cmd.Execute()` usage)
- `cmd/root.go` (split into helpers; Cobra root command goes away)

Create an `*cli.App` with:

- `Flags` holding all current root flags
- `Before` hook that loads config (unless `skipConfig` metadata is set)
- `Commands` listing all command builders

Global flags to preserve:

- `--config` (alias currently exists; decide whether to keep `-c`)
- `--db` / `-d`
- `--json` / `-j`
- `--no-color`

Suggested app skeleton (note: `-c` for config conflicts with `-c` for context on some commands; consider dropping it):

```go
app := &cli.App{
    Name:                 "ugh",
    Usage:                "todo.txt-inspired task CLI",
    Description:          "A CLI task manager using todo.txt format with libSQL storage",
    Version:              version, // set this from build info
    EnableBashCompletion: true,
    Flags: []cli.Flag{
        &cli.StringFlag{
            Name:     "config",
            EnvVars:  []string{"UGH_CONFIG"},
            Category: "Global",
        },
        &cli.StringFlag{
            Name:     "db",
            Aliases:  []string{"d"},
            EnvVars:  []string{"UGH_DB"},
            Category: "Global",
        },
        &cli.BoolFlag{
            Name:     "json",
            Aliases:  []string{"j"},
            EnvVars:  []string{"UGH_JSON"},
            Category: "Global",
        },
        &cli.BoolFlag{
            Name:     "no-color",
            EnvVars:  []string{"UGH_NO_COLOR"},
            Category: "Global",
        },
    },
    Before: func(c *cli.Context) error {
        if skip, ok := c.Command.Metadata["skipConfig"].(bool); ok && skip {
            return nil
        }
        return loadConfig(c)
    },
    Commands: []*cli.Command{
        AddCommand(),
        ListCommand(),
        // ...
        ConfigCommand(),
        DaemonCommand(),
    },
}
```

**Notes on App fields:**
- `Usage`: Short description shown in help (like Cobra's `Short`)
- `Description`: Detailed help text shown with `--help` (like Cobra's `Long`)
- `Version`: Enables `--version` flag automatically
- `EnableBashCompletion`: Enables bash completion support
- `Category`: Groups flags in help output for better organization

### 3) Implement config load + shared helpers

Goal: get rid of leaking Viper instances and keep a simple shared config object accessible to commands.

Files:

- `cmd/root.go` (or a new file like `cmd/app.go` / `cmd/appctx.go`)

Implement:

- `loadConfig(c *cli.Context) error`
- `getConfig(c *cli.Context) *config.Config`
- `outputWriterFromContext(c *cli.Context) output.Writer` (or keep `outputWriter()` but make it read global flags from `c`)

Store values in metadata:

- `c.App.Metadata["config"] = &cfg` (pointer to loaded config)
- `c.App.Metadata["configWasLoaded"] = result.WasLoaded`
- `c.App.Metadata["configPath"] = result.UsedPath` (optional, for relative path resolution)

Use `c.IsSet("db")` to decide whether `--db` should override.

### 4) Simplify internal config implementation (optional but recommended)

Files:

- `internal/config/config.go`

Changes:

- Remove `LoadResult.Viper *viper.Viper`
- Remove `Watch()` entirely
- Replace Viper-based parsing with a small TOML parser (or keep Viper internally but do not expose it)
- Remove home-directory restriction / validation logic

Note: You can keep using Viper internally for config loading even after migrating CLI to urfave/cli. The key change is not exposing the Viper instance in `LoadResult`.

Implementation guidance (simple + reliable):

- Keep `Load(path, allowMissing)` returning `(Config + WasLoaded + UsedPath)`
- Parse TOML (file only)
- Apply defaults in code
- Apply env vars either:
  - via urfave/cli flags (preferred: env handled by CLI automatically), or
  - via explicit `os.Getenv` mapping (only if you have config values not represented by flags)

**Note:** urfave/cli automatically binds `EnvVars` to flags. When a flag is accessed via `c.String("flag")`, it will already contain the env var value if no flag was explicitly set. This means you typically don't need explicit `os.Getenv` calls for values that have corresponding flags.

### 5) Convert each Cobra command to a cli.Command builder

Pattern change:

- From: `var fooCmd = &cobra.Command{...}` + `func init() { fooCmd.Flags()... }`
- To: `func FooCommand() *cli.Command { return &cli.Command{...} }`

General mapping:

- Cobra `Use`: urfave/cli `Name` + `ArgsUsage` (optional)
- Cobra `Short`: urfave/cli `Usage`
- Cobra `Long`: urfave/cli `Description`
- Cobra `Aliases`: urfave/cli `Aliases`
- Cobra `RunE`: urfave/cli `Action`
- Cobra per-command flags: urfave/cli `Flags: []cli.Flag{...}`

**Flag binding:**
- Use `Destination: &opts.Field` for simple binding
- Use `c.String("flag")` / `c.Bool("flag")` when you need to distinguish "unset" vs default

**Action return values:**
- Return `nil` for success
- Return `error` to print to stderr and exit with non-zero status (automatic)}

Accessing global flags in commands: App-level flags are inherited by all commands, so `c.String("db")` works in any command Action.

Commands to migrate:

- `cmd/add.go`
- `cmd/list.go`
- `cmd/show.go`
- `cmd/edit.go`
- `cmd/done.go`
- `cmd/undo.go`
- `cmd/rm.go`
- `cmd/import.go`
- `cmd/export.go`
- `cmd/projects.go`
- `cmd/contexts.go`
- `cmd/sync.go` (plus subcommands)
- `cmd/config/config.go` (plus init/show/get/set)
- `cmd/daemon/daemon.go` (plus install/uninstall/start/stop/restart/status/logs/run)

**Subcommand pattern** (for daemon, config, sync):

```go
func DaemonCommand() *cli.Command {
    return &cli.Command{
        Name:        "daemon",
        Usage:       "Manage the background daemon",
        Description: "Install, start, stop, and manage the ugh daemon",
        Subcommands: []*cli.Command{
            DaemonInstallCommand(),
            DaemonStartCommand(),
            DaemonStopCommand(),
            // ...
        },
    }
}
```

### 6) Apply skipConfig metadata to config init/set

Files:

- `cmd/config/init.go`
- `cmd/config/set.go`

Add:

```go
Metadata: map[string]interface{}{
    "skipConfig": true,
},
```

### 7) Review and simplify short flags

Create a single table of global + command flags and decide what to keep.

Suggested keep (high value, low conflict):

- Global: `--db/-d`, `--json/-j`
- Common command: `add -a` alias for command name is fine (command alias), `list -l` alias for command name is fine

Suggested review carefully:

- `--config/-c` may conflict with `--context/-c` on some commands (currently `add` uses `-c` for context)

**Flag display options:**

Use `HideDefault: true` on flags for cleaner help output when defaults are empty or not meaningful:

```go
&cli.StringFlag{
    Name:        "config",
    Usage:       "Path to config file",
    EnvVars:     []string{"UGH_CONFIG"},
    HideDefault: true,
}
```

If you drop `-c` for config, ensure docs/test scripts are updated.

### 8) Update tests and docs

Files:

- `testdata/script/*.txt`
- `README.md` and any docs referencing Cobra help output / flags

What to verify:

- Any script that asserts help output will need updating
- Any script that assumes a particular error string from Cobra may need updating
- Confirm the precedence behavior:
  - `--db` overrides config `db.path`
  - `UGH_DB` overrides config `db.path` when `--db` is not provided

### 9) Build, lint, test

Run:

- `just fmt`
- `just vet`
- `just test`

If Nix is used in CI:

- `nix flake check`

## Suggested Order of Work (Practical)

1. Add `cli.App` in `main.go` while still calling existing Cobra command builders (temporary bridge is possible but usually not worth it)
2. Convert one command end-to-end (ex: `add`) to establish the pattern
3. Convert all commands mechanically
4. Delete Cobra-specific code paths
5. Simplify `internal/config` (or do it early if it blocks removing Viper)
6. Update tests

## Checklist

- [ ] `main.go` builds a `cli.App`
- [ ] Global flags exist and work across all commands
- [ ] Config is loaded once per invocation (`Before` hook)
- [ ] `ugh config init` and `ugh config set` work without an existing config file
- [ ] No config watching remains
- [ ] Cobra and Viper dependencies removed
- [ ] Tests updated and passing
