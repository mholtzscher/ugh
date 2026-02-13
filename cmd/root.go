package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/pterm/pterm"
	"github.com/urfave/cli/v3"
	"golang.org/x/term"

	"github.com/mholtzscher/ugh/internal/config"
	"github.com/mholtzscher/ugh/internal/flags"
	"github.com/mholtzscher/ugh/internal/output"
	"github.com/mholtzscher/ugh/internal/store"
)

const (
	defaultDaemonSyncRetryMax = 3
	openStoreMaxRetries       = 5
	openStoreInitialBackoff   = 100 * time.Millisecond
)

// Version is set at build time.
//
//nolint:gochecknoglobals // Version must be mutable for ldflags injection at build time.
var Version = "0.2.0" // x-release-please-version

//nolint:gochecknoglobals // Root CLI flags and config cache are process-wide command state.
var (
	rootConfigPath  string
	rootDBPath      string
	rootJSON        bool
	rootNoColor     bool
	loadedConfig    *config.Config
	loadedConfigWas bool
)

//nolint:gochecknoglobals // Root command is package-level for CLI registration.
var rootCmd = &cli.Command{
	Name:        "ugh",
	Usage:       "ugh is a task CLI",
	Description: "ugh is a task CLI with SQLite storage.",
	Version:     Version,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  flags.FlagConfigPath,
			Usage: "path to config file",
		},
		&cli.StringFlag{
			Name:    flags.FlagDBPath,
			Aliases: []string{"d"},
			Usage:   "path to sqlite database (overrides config)",
		},
		&cli.BoolFlag{
			Name:    flags.FlagJSON,
			Aliases: []string{"j"},
			Usage:   "output json",
		},
		&cli.BoolFlag{
			Name:  flags.FlagNoColor,
			Usage: "disable color output",
		},
	},
	Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
		return ctx, loadConfig(cmd)
	},
	Action: func(_ context.Context, cmd *cli.Command) error {
		return cli.ShowRootCommandHelp(cmd)
	},
	Commands: []*cli.Command{
		addCmd,
		inboxCmd,
		nowCmd,
		waitingCmd,
		laterCmd,
		calendarCmd,
		listCmd,
		logCmd,
		showCmd,
		editCmd,
		doneCmd,
		undoCmd,
		rmCmd,
		projectsCmd,
		contextsCmd,
		syncCmd,
		configCmd,
		daemonCmd,
		shellCmd,
		historyCmd,
		seedCmd,
	},
}

//nolint:nestif // Config auto-initialization flow is explicit for readability.
func loadConfig(cmd *cli.Command) error {
	if cmd != nil {
		rootConfigPath = cmd.String(flags.FlagConfigPath)
		rootDBPath = cmd.String(flags.FlagDBPath)
		rootJSON = cmd.Bool(flags.FlagJSON)
		rootNoColor = cmd.Bool(flags.FlagNoColor)
	}

	if rootNoColor || os.Getenv("NO_COLOR") != "" {
		pterm.DisableColor()
	}

	configPath := rootConfigPath
	result, err := config.Load(configPath, true)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if !result.WasLoaded {
		if shouldAutoInitConfig() {
			cfgPath := configPath
			if cfgPath == "" {
				var defaultPath string
				defaultPath, err = config.DefaultPath()
				if err != nil {
					return fmt.Errorf("get default config path: %w", err)
				}
				cfgPath = defaultPath
			}

			var dbPath string
			dbPath, err = defaultDBPath()
			if err != nil {
				return err
			}

			result.Config = config.Config{
				Version: config.DefaultVersion,
				DB:      config.DB{Path: dbPath},
				Daemon: config.Daemon{
					PeriodicSync:     "5m",
					LogLevel:         "info",
					SyncRetryMax:     defaultDaemonSyncRetryMax,
					SyncRetryBackoff: "1s",
				},
			}

			err = config.Save(cfgPath, result.Config)
			if err != nil {
				return fmt.Errorf("save config: %w", err)
			}
			result.WasLoaded = true
		}
	}

	cacheLoadedConfig(result)
	return nil
}

func cacheLoadedConfig(result *config.LoadResult) {
	loadedConfig = &result.Config
	loadedConfigWas = result.WasLoaded
}

func shouldAutoInitConfig() bool {
	if rootConfigPath != "" {
		return true
	}
	return rootDBPath == ""
}

//nolint:nestif // Resolution order intentionally mirrors CLI precedence.
func effectiveDBPath() (string, error) {
	if rootDBPath != "" {
		return rootDBPath, nil
	}
	if loadedConfig != nil && loadedConfig.DB.Path != "" {
		var cfgPath string
		if rootConfigPath != "" {
			cfgPath = rootConfigPath
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
	if err := rootCmd.Run(context.Background(), os.Args); err != nil {
		writer := output.Writer{
			Out:  os.Stderr,
			JSON: false,
			TTY:  term.IsTerminal(int(os.Stderr.Fd())),
		}
		_ = writer.WriteError(err.Error())
		os.Exit(1)
	}
}

func openStore(ctx context.Context) (*store.Store, error) {
	path, err := effectiveDBPath()
	if err != nil {
		return nil, err
	}
	dbDir := filepath.Dir(path)
	err = os.MkdirAll(dbDir, 0o750)
	if err != nil {
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
	maxRetries := openStoreMaxRetries
	backoff := openStoreInitialBackoff
	for i := range maxRetries {
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
	return output.NewWriter(rootJSON)
}
