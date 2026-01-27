package config

import (
	"encoding/json"

	"github.com/mholtzscher/ugh/internal/config"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show configuration",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := deps.Config()
		if cfg == nil {
			cfg = &config.Config{Version: config.DefaultVersion}
		}

		writer := deps.OutputWriter()
		if writer.JSON {
			enc := json.NewEncoder(writer.Out)
			return enc.Encode(cfg)
		}
		return toml.NewEncoder(writer.Out).Encode(cfg)
	},
}
