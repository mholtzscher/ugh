package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/mholtzscher/ugh/internal/config"

	"github.com/urfave/cli/v2"
)

type initResult struct {
	Action string `json:"action"`
	File   string `json:"file"`
	DBPath string `json:"dbPath"`
}

func initCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initialize configuration file",
		Action: func(c *cli.Context) error {
			cfgPath, err := configPathForWrite(c)
			if err != nil {
				return err
			}
			if _, err := os.Stat(cfgPath); err == nil {
				return fmt.Errorf("config file already exists: %s", cfgPath)
			} else if !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("stat config: %w", err)
			}

			dbPath, err := deps.DefaultDBPath()
			if err != nil {
				return err
			}
			cfg := config.Config{
				Version: config.DefaultVersion,
				DB:      config.DB{Path: dbPath},
				Daemon: config.Daemon{
					PeriodicSync:     "5m",
					LogLevel:         "info",
					SyncRetryMax:     3,
					SyncRetryBackoff: "1s",
				},
			}
			if err := config.Save(cfgPath, cfg); err != nil {
				return err
			}
			deps.SetConfig(&cfg)
			deps.SetConfigWasLoaded(true)

			writer := deps.OutputWriter(c)
			if writer.JSON {
				enc := json.NewEncoder(writer.Out)
				return enc.Encode(initResult{Action: "init", File: cfgPath, DBPath: dbPath})
			}
			_, err = fmt.Fprintf(writer.Out, "initialized %s\n", cfgPath)
			return err
		},
	}
}
