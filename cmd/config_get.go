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
			cfg = &config.Config{Version: config.DefaultVersion}
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
	case configKeyDBPath:
		if cfg.DB.Path == "" {
			return "", errors.New("config key not set: " + configKeyDBPath)
		}
		return cfg.DB.Path, nil
	case configKeyDBSyncURL:
		if cfg.DB.SyncURL == "" {
			return "", errors.New("config key not set: " + configKeyDBSyncURL)
		}
		return cfg.DB.SyncURL, nil
	case configKeyDBAuthToken:
		if cfg.DB.AuthToken == "" {
			return "", errors.New("config key not set: " + configKeyDBAuthToken)
		}
		return cfg.DB.AuthToken, nil
	case configKeyDBSyncOnWrite:
		return strconv.FormatBool(cfg.DB.SyncOnWrite), nil
	default:
		return "", fmt.Errorf("unknown config key: %s", key)
	}
}
