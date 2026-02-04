package config

import (
	"github.com/mholtzscher/ugh/internal/config"
	"github.com/mholtzscher/ugh/internal/output"

	"github.com/urfave/cli/v3"
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
	OutputWriter func() output.Writer
	// ConfigPath returns the user-specified config path (may be empty).
	ConfigPath func() string
	// DefaultDBPath returns the default database path.
	DefaultDBPath func() (string, error)
}

var deps Deps

// Cmd is the parent command for all config subcommands.
var Cmd = &cli.Command{
	Name:     "config",
	Usage:    "Manage configuration",
	Category: "System",
}

// Register adds the config command and its subcommands to the parent command.
// Must be called with valid Deps before the command tree is executed.
func Register(parent *cli.Command, d Deps) {
	deps = d
	Cmd.Commands = []*cli.Command{initCmd, showCmd, getCmd, setCmd}
	parent.Commands = append(parent.Commands, Cmd)
}

// configPathForWrite returns the path to write config to.
// Uses the user-specified path if set, otherwise the default.
func configPathForWrite() (string, error) {
	if path := deps.ConfigPath(); path != "" {
		return path, nil
	}
	return config.DefaultPath()
}
