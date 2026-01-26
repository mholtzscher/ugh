package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mholtzscher/ugh/internal/config"
	"github.com/spf13/cobra"
)

type configGetResult struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := loadedConfig
		if cfg == nil {
			cfg = &config.Config{Version: config.DefaultVersion}
		}

		key := args[0]
		value, err := configGetValue(cfg, key)
		if err != nil {
			return err
		}

		writer := outputWriter()
		if writer.JSON {
			enc := json.NewEncoder(writer.Out)
			return enc.Encode(configGetResult{Key: key, Value: value})
		}
		_, err = fmt.Fprintln(writer.Out, value)
		return err
	},
}

func configGetValue(cfg *config.Config, key string) (string, error) {
	switch key {
	case "db.path":
		if cfg.DB.Path == "" {
			return "", errors.New("config key not set: db.path")
		}
		return cfg.DB.Path, nil
	default:
		return "", fmt.Errorf("unknown config key: %s", key)
	}
}
