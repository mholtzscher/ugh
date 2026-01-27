package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

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
	rootOpts        rootOptions
	loadedConfig    *config.Config
	loadedConfigWas bool
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
	if loadedConfig != nil && loadedConfig.DB.Path != "" {
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
		dbPath, err := config.ResolveDBPath(cfgPath, loadedConfig.DB.Path)
		if err != nil {
			return "", fmt.Errorf("resolve db path from config: %w", err)
		}
		return dbPath, nil
	}
	return defaultDBPath()
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
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(projectsCmd)
	rootCmd.AddCommand(contextsCmd)
	rootCmd.AddCommand(syncCmd)
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
	st, err := store.Open(ctx, opts)
	if err != nil {
		return nil, err
	}
	return st, nil
}

func defaultDBPath() (string, error) {
	dataDir, err := userDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dataDir, "ugh", "ugh.sqlite"), nil
}

func userDataDir() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		// macOS's "config" dir is ~/Library/Application Support.
		configDir, err := os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("user config dir: %w", err)
		}
		return configDir, nil
	case "windows":
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			return localAppData, nil
		}
		configDir, err := os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("user config dir: %w", err)
		}
		return configDir, nil
	default:
		if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {
			return xdgDataHome, nil
		}
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("home dir: %w", err)
		}
		return filepath.Join(homeDir, ".local", "share"), nil
	}
}

func outputWriter() output.Writer {
	return output.NewWriter(rootOpts.JSON, rootOpts.NoColor)
}
