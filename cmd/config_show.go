package cmd

import (
	"context"
	"encoding/json"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v3"

	"github.com/mholtzscher/ugh/internal/config"
)

//nolint:gochecknoglobals // CLI command definitions are package-level by design.
var configShowCmd = &cli.Command{
	Name:  "show",
	Usage: "Show configuration",
	Action: func(_ context.Context, _ *cli.Command) error {
		cfg := loadedConfig
		if cfg == nil {
			cfg = &config.Config{Version: config.DefaultVersion}
		}

		writer := outputWriter()
		if writer.JSON {
			enc := json.NewEncoder(writer.Out)
			return enc.Encode(configShowJSONPayload(cfg))
		}
		return toml.NewEncoder(writer.Out).Encode(cfg)
	},
}

func configShowJSONPayload(cfg *config.Config) map[string]any {
	return map[string]any{
		"version": cfg.Version,
		"db": map[string]any{
			"path":          cfg.DB.Path,
			"sync_url":      cfg.DB.SyncURL,
			"auth_token":    cfg.DB.AuthToken,
			"sync_on_write": cfg.DB.SyncOnWrite,
		},
		"daemon": map[string]any{
			"periodic_sync":      cfg.Daemon.PeriodicSync,
			"log_file":           cfg.Daemon.LogFile,
			"log_level":          cfg.Daemon.LogLevel,
			"sync_retry_max":     cfg.Daemon.SyncRetryMax,
			"sync_retry_backoff": cfg.Daemon.SyncRetryBackoff,
		},
	}
}
