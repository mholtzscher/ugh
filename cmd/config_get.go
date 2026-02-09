package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/urfave/cli/v3"

	"github.com/mholtzscher/ugh/internal/config"
)

type configGetResult struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

//nolint:gochecknoglobals // CLI command definitions are package-level by design.
var configGetCmd = &cli.Command{
	Name:      "get",
	Usage:     "Get a configuration value",
	ArgsUsage: "<key>",
	Action: func(_ context.Context, cmd *cli.Command) error {
		if cmd.Args().Len() != 1 {
			return errors.New("get requires a key")
		}
		cfg := loadedConfig
		if cfg == nil {
			cfg = &config.Config{Version: config.DefaultVersion, UI: config.UI{Theme: config.DefaultUITheme}}
		}

		key := cmd.Args().Get(0)
		value, err := getConfigValue(cfg, key)
		if err != nil {
			return err
		}

		writer := outputWriter()
		if writer.JSON {
			enc := json.NewEncoder(writer.Out)
			return enc.Encode(configGetResult{Key: key, Value: value})
		}
		return writer.WriteLine(value)
	},
}

func getConfigValue(cfg *config.Config, key string) (string, error) {
	switch key {
	case "db.path":
		if cfg.DB.Path == "" {
			return "", errors.New("config key not set: db.path")
		}
		return cfg.DB.Path, nil
	case "db.sync_url":
		if cfg.DB.SyncURL == "" {
			return "", errors.New("config key not set: db.sync_url")
		}
		return cfg.DB.SyncURL, nil
	case "db.auth_token":
		if cfg.DB.AuthToken == "" {
			return "", errors.New("config key not set: db.auth_token")
		}
		return cfg.DB.AuthToken, nil
	case "db.sync_on_write":
		return strconv.FormatBool(cfg.DB.SyncOnWrite), nil
	case "ui.theme":
		if cfg.UI.Theme == "" {
			return config.DefaultUITheme, nil
		}
		return cfg.UI.Theme, nil
	default:
		return "", fmt.Errorf("unknown config key: %s", key)
	}
}
