package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/mholtzscher/ugh/internal/config"
)

const configUnsetArgCount = 1

type configUnsetResult struct {
	Action string `json:"action"`
	Key    string `json:"key"`
	File   string `json:"file"`
}

//nolint:gochecknoglobals // CLI command definitions are package-level by design.
var configUnsetCmd = &cli.Command{
	Name:      "unset",
	Usage:     "Unset (remove) a configuration value",
	ArgsUsage: "<key>",
	Action: func(_ context.Context, cmd *cli.Command) error {
		if cmd.Args().Len() != configUnsetArgCount {
			return errors.New("unset requires a key")
		}
		cfg := loadedConfig
		if cfg == nil {
			return errors.New("no config loaded")
		}

		key := cmd.Args().Get(0)
		if err := unsetConfigValue(cfg, key); err != nil {
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

		writer := outputWriter()
		if writer.JSON {
			enc := json.NewEncoder(writer.Out)
			return enc.Encode(configUnsetResult{Action: "unset", Key: key, File: cfgPath})
		}
		return writer.WriteSuccess("unset " + key)
	},
}

func unsetConfigValue(cfg *config.Config, key string) error {
	switch key {
	case configKeyDBPath:
		cfg.DB.Path = ""
		return nil
	case configKeyDBSyncURL:
		cfg.DB.SyncURL = ""
		return nil
	case configKeyDBAuthToken:
		cfg.DB.AuthToken = ""
		return nil
	case configKeyDBSyncOnWrite:
		cfg.DB.SyncOnWrite = false
		return nil
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}
}
