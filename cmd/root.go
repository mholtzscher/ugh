package cmd

import (
	"context"
	"fmt"
	"github.com/mholtzscher/ugh/internal/output"
	"github.com/mholtzscher/ugh/internal/store"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"time"
)

type rootOptions struct {
	DBPath  string
	JSON    bool
	NoColor bool
}

var rootCmd = &cobra.Command{
	Use:          "ugh",
	Short:        "ugh is a todo.txt-inspired task CLI",
	Long:         "ugh is a todo.txt-inspired task CLI with SQLite storage.",
	SilenceUsage: true,
}

var rootOpts rootOptions

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootOpts.DBPath, "db", "d", "", "path to sqlite database")
	rootCmd.PersistentFlags().BoolVarP(&rootOpts.JSON, "json", "j", false, "output json")
	rootCmd.PersistentFlags().BoolVar(&rootOpts.NoColor, "no-color", false, "disable color output")
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(doneCmd)
	rootCmd.AddCommand(undoCmd)
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(exportCmd)
}

func openStore(ctx context.Context) (*store.Store, error) {
	path := rootOpts.DBPath
	if path == "" {
		var err error
		path, err = defaultDBPath()
		if err != nil {
			return nil, err
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create db dir: %w", err)
	}
	return store.Open(ctx, path)
}

func defaultDBPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("config dir: %w", err)
	}
	return filepath.Join(configDir, "ugh", "ugh.sqlite"), nil
}

func outputWriter() output.Writer {
	return output.NewWriter(rootOpts.JSON, rootOpts.NoColor)
}

func nowDate() *time.Time {
	val := time.Now().UTC()
	return &val
}
