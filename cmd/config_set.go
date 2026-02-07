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

const configSetArgCount = 2

type configSetResult struct {
	Action string `json:"action"`
	Key    string `json:"key"`
	Value  string `json:"value"`
	File   string `json:"file"`
}

//nolint:gochecknoglobals // CLI command definitions are package-level by design.
var configSetCmd = &cli.Command{
	Name:      "set",
	Usage:     "Set a configuration value",
	ArgsUsage: "<key> <value>",
	Action: func(_ context.Context, cmd *cli.Command) error {
		if cmd.Args().Len() != configSetArgCount {
			return errors.New("set requires a key and value")
		}
		cfg := loadedConfig
		if cfg == nil {
			cfg = &config.Config{Version: config.DefaultVersion, UI: config.UI{Theme: config.DefaultUITheme}}
		}

		key := cmd.Args().Get(0)
		value := cmd.Args().Get(1)
		if err := setConfigValue(cfg, key, value); err != nil {
			return err
		}

		cfgPath, err := configPathForWrite()
		if err != nil {
			return err
		}
		err = config.Save(cfgPath, *cfg)
		if err != nil {
			return err
		}
		loadedConfig = cfg
		loadedConfigWas = true

		writer := outputWriter()
		if writer.JSON {
			enc := json.NewEncoder(writer.Out)
			return enc.Encode(configSetResult{Action: "set", Key: key, Value: value, File: cfgPath})
		}
		_, err = fmt.Fprintf(writer.Out, "set %s\n", key)
		return err
	},
}

func setConfigValue(cfg *config.Config, key string, value string) error {
	switch key {
	case "db.path":
		cfg.DB.Path = value
		return nil
	case "db.sync_url":
		cfg.DB.SyncURL = value
		return nil
	case "db.auth_token":
		cfg.DB.AuthToken = value
		return nil
	case "db.sync_on_write":
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid boolean for db.sync_on_write: %w", err)
		}
		cfg.DB.SyncOnWrite = parsed
		return nil
	case "ui.theme":
		cfg.UI.Theme = value
		return nil
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}
}
