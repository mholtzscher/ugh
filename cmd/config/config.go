package config

import (
	"github.com/mholtzscher/ugh/internal/config"
	"github.com/mholtzscher/ugh/internal/output"

	"github.com/urfave/cli/v2"
)

// Deps holds dependencies injected from the parent cmd package.
// This avoids circular imports between cmd and cmd/config.
type Deps struct {
	// Config returns the currently loaded config (may be nil).
	Config func() *config.Config
	// SetConfig updates the loaded config.
	SetConfig func(*config.Config)
	// ConfigWasLoaded returns true if config was loaded from a file.
	ConfigWasLoaded func() bool
	// SetConfigWasLoaded updates the loaded-from-file flag.
	SetConfigWasLoaded func(bool)
	// OutputWriter returns the configured output writer.
	OutputWriter func(*cli.Context) output.Writer
	// ConfigPath returns the user-specified config path (may be empty).
	ConfigPath func(*cli.Context) string
	// DefaultDBPath returns the default database path.
	DefaultDBPath func() (string, error)
}

var deps Deps

// Command is the parent command for all config subcommands.
func Command(d Deps) *cli.Command {
	deps = d
	return &cli.Command{
		Name:  "config",
		Usage: "Manage configuration",
		Subcommands: []*cli.Command{
			initCommand(),
			showCommand(),
			getCommand(),
			setCommand(),
		},
	}
}

// configPathForWrite returns the path to write config to.
// Uses the user-specified path if set, otherwise the default.
func configPathForWrite(c *cli.Context) (string, error) {
	if path := deps.ConfigPath(c); path != "" {
		return path, nil
	}
	return config.DefaultPath()
}
