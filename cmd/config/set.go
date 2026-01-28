package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/mholtzscher/ugh/internal/config"

	"github.com/urfave/cli/v2"
)

type setResult struct {
	Action string `json:"action"`
	Key    string `json:"key"`
	Value  string `json:"value"`
	File   string `json:"file"`
}

func setCommand() *cli.Command {
	return &cli.Command{
		Name:      "set",
		Usage:     "Set a configuration value",
		ArgsUsage: "<key> <value>",
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 2 {
				return errors.New("set requires a key and value")
			}
			cfg := deps.Config()
			if cfg == nil {
				cfg = &config.Config{Version: config.DefaultVersion}
			}

			key := c.Args().Get(0)
			value := c.Args().Get(1)
			if err := setValue(cfg, key, value); err != nil {
				return err
			}

			cfgPath, err := configPathForWrite(c)
			if err != nil {
				return err
			}
			if err := config.Save(cfgPath, *cfg); err != nil {
				return err
			}
			deps.SetConfig(cfg)
			deps.SetConfigWasLoaded(true)

			writer := deps.OutputWriter(c)
			if writer.JSON {
				enc := json.NewEncoder(writer.Out)
				return enc.Encode(setResult{Action: "set", Key: key, Value: value, File: cfgPath})
			}
			_, err = fmt.Fprintf(writer.Out, "set %s\n", key)
			return err
		},
	}
}

func setValue(cfg *config.Config, key string, value string) error {
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
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}
}
