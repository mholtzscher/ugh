package config

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mholtzscher/ugh/internal/config"

	"github.com/urfave/cli/v2"
)

type getResult struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func getCommand() *cli.Command {
	return &cli.Command{
		Name:      "get",
		Usage:     "Get a configuration value",
		ArgsUsage: "<key>",
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				return errors.New("get requires a key")
			}
			cfg := deps.Config()
			if cfg == nil {
				cfg = &config.Config{Version: config.DefaultVersion}
			}

			key := c.Args().First()
			value, err := getValue(cfg, key)
			if err != nil {
				return err
			}

			writer := deps.OutputWriter(c)
			if writer.JSON {
				enc := json.NewEncoder(writer.Out)
				return enc.Encode(getResult{Key: key, Value: value})
			}
			_, err = fmt.Fprintln(writer.Out, value)
			return err
		},
	}
}

func getValue(cfg *config.Config, key string) (string, error) {
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
		return fmt.Sprintf("%t", cfg.DB.SyncOnWrite), nil
	default:
		return "", fmt.Errorf("unknown config key: %s", key)
	}
}
