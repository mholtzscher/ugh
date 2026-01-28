package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	configcmd "github.com/mholtzscher/ugh/cmd/config"
	daemoncmd "github.com/mholtzscher/ugh/cmd/daemon"
	"github.com/mholtzscher/ugh/internal/config"
	"github.com/mholtzscher/ugh/internal/output"
	"github.com/mholtzscher/ugh/internal/store"

	"github.com/spf13/cobra"
)

type rootOptions struct {
	ConfigPath string
	DBPath     string
	JSON       bool
	NoColor    bool
}

var (
	rootOpts           rootOptions
	loadedConfig       *config.Config
	loadedConfigWas    bool
	loadedConfigResult *config.LoadResult // Includes Viper instance for flag binding/watching
)

var rootCmd = &cobra.Command{
	Use:          "ugh",
	Short:        "ugh is a todo.txt-inspired task CLI",
	Long:         "ugh is a todo.txt-inspired task CLI with SQLite storage.",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return loadConfig(cmd)
	},
}

func loadConfig(cmd *cobra.Command) error {
	configPath := rootOpts.ConfigPath
	allowMissing := configPath == "" || allowMissingConfig(cmd)

	result, err := config.Load(configPath, allowMissing)
	if err != nil {
		if errors.Is(err, config.ErrNotFound) && allowMissing {
			loadedConfig = &config.Config{Version: config.DefaultVersion}
			loadedConfigWas = false
			return nil
		}
		if errors.Is(err, config.ErrNotFound) && !allowMissing {
			return fmt.Errorf("config file not found: %s", configPath)
		}
		return fmt.Errorf("load config: %w", err)
	}

	loadedConfig = &result.Config
	loadedConfigWas = result.WasLoaded
	loadedConfigResult = result

	return nil
}

func allowMissingConfig(cmd *cobra.Command) bool {
	if cmd == nil {
		return false
	}
	if cmd.Name() == "set" || cmd.Name() == "init" {
		parent := cmd.Parent()
		return parent != nil && parent.Name() == "config"
	}
	return false
}

func effectiveDBPath() (string, error) {
	if rootOpts.DBPath != "" {
		return rootOpts.DBPath, nil
	}

	var cfgPath string
	if rootOpts.ConfigPath != "" {
		cfgPath = rootOpts.ConfigPath
	} else if loadedConfigWas {
		defaultPath, err := config.DefaultPath()
		if err != nil {
			return "", fmt.Errorf("get default config path: %w", err)
		}
		cfgPath = defaultPath
	}

	return config.EffectiveDBPath(loadedConfig, cfgPath)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&rootOpts.ConfigPath, "config", "", "path to config file")

	rootCmd.PersistentFlags().StringVarP(&rootOpts.DBPath, "db", "d", "", "path to sqlite database (overrides config)")
	rootCmd.PersistentFlags().BoolVarP(&rootOpts.JSON, "json", "j", false, "output json")
	rootCmd.PersistentFlags().BoolVar(&rootOpts.NoColor, "no-color", false, "disable color output")

	rootCmd.PersistentFlags().SortFlags = false

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(doneCmd)
	rootCmd.AddCommand(undoCmd)
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(projectsCmd)
	rootCmd.AddCommand(contextsCmd)
	rootCmd.AddCommand(syncCmd)

	// Register subcommand packages
	configcmd.Register(rootCmd, configcmd.Deps{
		Config:             func() *config.Config { return loadedConfig },
		SetConfig:          func(c *config.Config) { loadedConfig = c },
		ConfigWasLoaded:    func() bool { return loadedConfigWas },
		SetConfigWasLoaded: func(b bool) { loadedConfigWas = b },
		OutputWriter:       outputWriter,
		ConfigPath:         func() string { return rootOpts.ConfigPath },
		DefaultDBPath:      defaultDBPath,
	})

	daemoncmd.Register(rootCmd, daemoncmd.Deps{
		Config:       func() *config.Config { return loadedConfig },
		ConfigResult: func() *config.LoadResult { return loadedConfigResult },
		OutputWriter: outputWriter,
	})
}

func openStore(ctx context.Context) (*store.Store, error) {
	path, err := effectiveDBPath()
	if err != nil {
		return nil, err
	}
	dbDir := filepath.Dir(path)
	if err := os.MkdirAll(dbDir, 0o755); err != nil {
		return nil, fmt.Errorf("create db dir: %w", err)
	}
	cacheDir := filepath.Join(dbDir, ".cache")
	if os.Getenv("TURSO_GO_CACHE_DIR") == "" {
		_ = os.Setenv("TURSO_GO_CACHE_DIR", cacheDir)
	}

	opts := store.Options{Path: path}
	if loadedConfig != nil {
		opts.SyncURL = loadedConfig.DB.SyncURL
		opts.AuthToken = loadedConfig.DB.AuthToken
	}

	// Retry with backoff if database is locked (e.g., daemon is running)
	var st *store.Store
	maxRetries := 5
	backoff := 100 * time.Millisecond
	for i := 0; i < maxRetries; i++ {
		st, err = store.Open(ctx, opts)
		if err == nil {
			return st, nil
		}
		// Check if it's a locking error
		if !isLockingError(err) {
			return nil, err
		}
		if i < maxRetries-1 {
			time.Sleep(backoff)
			backoff *= 2 // Exponential backoff
		}
	}
	return nil, fmt.Errorf("%w (is the daemon running? try 'ugh daemon stop' or use the HTTP API)", err)
}

// isLockingError checks if the error is a database locking error.
func isLockingError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "locked") || strings.Contains(errStr, "Locking")
}

func defaultDBPath() (string, error) {
	return config.DefaultDBPath()
}

func outputWriter() output.Writer {
	return output.NewWriter(rootOpts.JSON, rootOpts.NoColor)
}
