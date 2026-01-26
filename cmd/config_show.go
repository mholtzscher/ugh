package cmd

import (
	"encoding/json"

	"github.com/mholtzscher/ugh/internal/config"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show configuration",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := loadedConfig
		if cfg == nil {
			cfg = &config.Config{Version: config.DefaultVersion}
		}

		writer := outputWriter()
		if writer.JSON {
			enc := json.NewEncoder(writer.Out)
			return enc.Encode(cfg)
		}
		return toml.NewEncoder(writer.Out).Encode(cfg)
	},
}
