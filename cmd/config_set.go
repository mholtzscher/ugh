package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mholtzscher/ugh/internal/config"
	"github.com/spf13/cobra"
)

type configSetResult struct {
	Action string `json:"action"`
	Key    string `json:"key"`
	Value  string `json:"value"`
	File   string `json:"file"`
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := loadedConfig
		if cfg == nil {
			cfg = &config.Config{Version: config.DefaultVersion}
		}

		key := args[0]
		value := args[1]
		if err := configSetValue(cfg, key, value); err != nil {
			return err
		}

		cfgPath, err := configPathForWrite()
		if err != nil {
			return err
		}
		if err := config.Save(cfgPath, *cfg); err != nil {
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

func configSetValue(cfg *config.Config, key string, value string) error {
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
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}
}

func configPathForWrite() (string, error) {
	if rootOpts.ConfigPath != "" {
		return rootOpts.ConfigPath, nil
	}
	return config.DefaultPath()
}
