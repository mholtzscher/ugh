package config

import (
	"encoding/json"

	"github.com/mholtzscher/ugh/internal/config"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v2"
)

func showCommand() *cli.Command {
	return &cli.Command{
		Name:  "show",
		Usage: "Show configuration",
		Action: func(c *cli.Context) error {
			cfg := deps.Config()
			if cfg == nil {
				cfg = &config.Config{Version: config.DefaultVersion}
			}

			writer := deps.OutputWriter(c)
			if writer.JSON {
				enc := json.NewEncoder(writer.Out)
				return enc.Encode(cfg)
			}
			return toml.NewEncoder(writer.Out).Encode(cfg)
		},
	}
}
