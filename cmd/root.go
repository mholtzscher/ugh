package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	configcmd "github.com/mholtzscher/ugh/cmd/config"
	daemoncmd "github.com/mholtzscher/ugh/cmd/daemon"
	"github.com/mholtzscher/ugh/internal/config"
	"github.com/mholtzscher/ugh/internal/flags"
	"github.com/mholtzscher/ugh/internal/output"
	"github.com/mholtzscher/ugh/internal/store"

	"github.com/urfave/cli/v3"
)

// Version is set at build time.
var Version = "0.1.1" // x-release-please-version

var (
	rootConfigPath  string
	rootDBPath      string
	rootJSON        bool
	rootNoColor     bool
	loadedConfig    *config.Config
	loadedConfigWas bool
)

var rootCmd = &cli.Command{
	Name:        "ugh",
	Usage:       "ugh is a GTD-first task CLI",
	Description: "ugh is a GTD-first task CLI with SQLite storage.",
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
}

func loadConfig(cmd *cli.Command) error {
	if cmd != nil {
		rootConfigPath = cmd.String(flags.FlagConfigPath)
		rootDBPath = cmd.String(flags.FlagDBPath)
		rootJSON = cmd.Bool(flags.FlagJSON)
		rootNoColor = cmd.Bool(flags.FlagNoColor)
	}

	configPath := rootConfigPath
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

func allowMissingConfig(cmd *cli.Command) bool {
	if cmd == nil {
		return false
	}
	if !cmd.HasName("set") && !cmd.HasName("init") {
		if cmd.Root() != cmd {
			return false
		}
		args := cmd.Args()
		if args.Len() < 2 || args.Get(0) != "config" {
			return false
		}
		return args.Get(1) == "set" || args.Get(1) == "init"
	}
	lineage := cmd.Lineage()
	if len(lineage) < 2 {
		return false
	}
	parent := lineage[len(lineage)-2]
	return parent != nil && parent.HasName("config")
}

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
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func init() {
	rootCmd.Commands = []*cli.Command{
		addCmd,
		inboxCmd,
		nextCmd,
		waitingCmd,
		somedayCmd,
		snoozedCmd,
		calendarCmd,
		listCmd,
		showCmd,
		editCmd,
		doneCmd,
		undoCmd,
		rmCmd,
		importCmd,
		exportCmd,
		projectsCmd,
		contextsCmd,
		syncCmd,
	}

	// Register subcommand packages
	configcmd.Register(rootCmd, configcmd.Deps{
		Config:             func() *config.Config { return loadedConfig },
		SetConfig:          func(c *config.Config) { loadedConfig = c },
		ConfigWasLoaded:    func() bool { return loadedConfigWas },
		SetConfigWasLoaded: func(b bool) { loadedConfigWas = b },
		OutputWriter:       outputWriter,
		ConfigPath:         func() string { return rootConfigPath },
		DefaultDBPath:      defaultDBPath,
	})

	daemoncmd.Register(rootCmd, daemoncmd.Deps{
		Config:       func() *config.Config { return loadedConfig },
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
	return output.NewWriter(rootJSON, rootNoColor)
}
