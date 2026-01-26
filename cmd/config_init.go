package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/mholtzscher/ugh/internal/config"
	"github.com/spf13/cobra"
)

type configInitResult struct {
	Action string `json:"action"`
	File   string `json:"file"`
	DBPath string `json:"dbPath"`
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgPath, err := configPathForWrite()
		if err != nil {
			return err
		}
		if _, err := os.Stat(cfgPath); err == nil {
			return fmt.Errorf("config file already exists: %s", cfgPath)
		} else if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("stat config: %w", err)
		}

		dbPath, err := defaultDBPath()
		if err != nil {
			return err
		}
		cfg := config.Config{Version: config.DefaultVersion, DB: config.DB{Path: dbPath}}
		if err := config.Save(cfgPath, cfg); err != nil {
			return err
		}
		loadedConfig = &cfg
		loadedConfigWas = true

		writer := outputWriter()
		if writer.JSON {
			enc := json.NewEncoder(writer.Out)
			return enc.Encode(configInitResult{Action: "init", File: cfgPath, DBPath: dbPath})
		}
		_, err = fmt.Fprintf(writer.Out, "initialized %s\n", cfgPath)
		return err
	},
}
