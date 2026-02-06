package config

import (
	"github.com/mholtzscher/ugh/cmd/meta"
	"github.com/mholtzscher/ugh/cmd/registry"
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
	// SetConfigWasLoaded updates the loaded-from-file flag.
	SetConfigWasLoaded func(bool)
	// OutputWriter returns the configured output writer.
	OutputWriter func() output.Writer
	// ConfigPath returns the user-specified config path (may be empty).
	ConfigPath func() string
}

const (
	configID   registry.ID = "config"
	configShow registry.ID = "config.show"
	configGet  registry.ID = "config.get"
	configSet  registry.ID = "config.set"
)

// Register adds config command specs to the registry.
func Register(r *registry.Registry, d Deps) error {
	return r.AddAll(
		registry.Spec{ID: configID, Source: "cmd/config", Build: newConfigCmd},
		registry.Spec{ID: configShow, ParentID: configID, Source: "cmd/config", Build: func() *cli.Command { return newShowCmd(d) }},
		registry.Spec{ID: configGet, ParentID: configID, Source: "cmd/config", Build: func() *cli.Command { return newGetCmd(d) }},
		registry.Spec{ID: configSet, ParentID: configID, Source: "cmd/config", Build: func() *cli.Command { return newSetCmd(d) }},
	)
}

func newConfigCmd() *cli.Command {
	return &cli.Command{
		Name:     "config",
		Usage:    "Manage configuration",
		Category: meta.SystemCategory.String(),
	}
}

// configPathForWrite returns the path to write config to.
// Uses the user-specified path if set, otherwise the default.
func configPathForWrite(d Deps) (string, error) {
	if path := d.ConfigPath(); path != "" {
		return path, nil
	}
	return config.DefaultPath()
}
