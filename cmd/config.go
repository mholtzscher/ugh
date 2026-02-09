package cmd

import (
	"github.com/mholtzscher/ugh/internal/config"

	"github.com/urfave/cli/v3"
)

// configCmd is the parent command for all config subcommands.
//
//nolint:gochecknoglobals // CLI command definitions are package-level by design.
var configCmd = &cli.Command{
	Name:     "config",
	Usage:    "Manage configuration",
	Category: "System",
	Commands: []*cli.Command{
		configShowCmd,
		configGetCmd,
		configSetCmd,
		configUnsetCmd,
	},
}

// configPathForWrite returns the path to write config to.
// Uses the user-specified path if set, otherwise the default.
func configPathForWrite() (string, error) {
	if rootConfigPath != "" {
		return rootConfigPath, nil
	}
	return config.DefaultPath()
}
