package config

import (
	"context"
	"encoding/json"

	"github.com/mholtzscher/ugh/internal/config"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v3"
)

func newShowCmd(d Deps) *cli.Command {
	return &cli.Command{
		Name:  "show",
		Usage: "Show configuration",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			cfg := d.Config()
			if cfg == nil {
				cfg = &config.Config{Version: config.DefaultVersion}
			}

			writer := d.OutputWriter()
			if writer.JSON {
				enc := json.NewEncoder(writer.Out)
				return enc.Encode(cfg)
			}
			return toml.NewEncoder(writer.Out).Encode(cfg)
		},
	}
}
