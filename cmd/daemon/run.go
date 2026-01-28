package daemon

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/mholtzscher/ugh/internal/config"
	"github.com/mholtzscher/ugh/internal/daemon"
	"github.com/mholtzscher/ugh/internal/store"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the daemon in foreground",
	Long: `Run the daemon server in the foreground.

This is primarily used by the system service manager (systemd/launchd).
For debugging, you can run it directly to see logs in the terminal.

The daemon provides background sync to Turso cloud on a periodic interval.
It only opens the database when syncing, avoiding lock contention with the CLI.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// Get config
		cfg := getConfig()
		if cfg == nil {
			return fmt.Errorf("config not loaded")
		}

		// Parse daemon config
		daemonCfg := daemon.ParseConfig(cfg.Daemon)

		// Set up logging
		logger, logCloser, err := setupLogging(daemonCfg.LogFile, daemonCfg.LogLevel)
		if err != nil {
			return fmt.Errorf("setup logging: %w", err)
		}
		if logCloser != nil {
			defer func() { _ = logCloser.Close() }()
		}

		// Get database path
		dbPath, err := config.EffectiveDBPath(cfg, "")
		if err != nil {
			return fmt.Errorf("get db path: %w", err)
		}

		// Ensure DB directory exists
		dbDir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dbDir, 0o755); err != nil {
			return fmt.Errorf("create db dir: %w", err)
		}

		// Set cache dir for Turso
		cacheDir := filepath.Join(dbDir, ".cache")
		if os.Getenv("TURSO_GO_CACHE_DIR") == "" {
			_ = os.Setenv("TURSO_GO_CACHE_DIR", cacheDir)
		}

		// Build store options (daemon will open/close as needed)
		storeOpts := store.Options{
			Path:        dbPath,
			SyncURL:     cfg.DB.SyncURL,
			AuthToken:   cfg.DB.AuthToken,
			BusyTimeout: 5000, // 5 seconds - short since we open/close quickly
		}

		// Create daemon
		d := daemon.New(storeOpts, daemonCfg, logger)

		// Set up config watching for hot-reload
		if result := deps.ConfigResult(); result != nil && result.Viper != nil {
			stopWatching := config.Watch(result.Viper, func(newCfg config.Config) {
				logger.Info("config file changed, reloading daemon config")
				newDaemonCfg := daemon.ParseConfig(newCfg.Daemon)
				d.UpdateConfig(newDaemonCfg)
			})
			defer stopWatching()
		}

		return d.Run(ctx)
	},
}

// setupLogging sets up the logger based on config.
func setupLogging(logFile, logLevel string) (*slog.Logger, io.Closer, error) {
	// Determine output
	var out io.Writer = os.Stderr
	var closer io.Closer

	if logFile != "" {
		// Expand ~ to home dir
		path := logFile
		if len(path) >= 2 && path[:2] == "~/" {
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, nil, err
			}
			path = filepath.Join(home, path[2:])
		}

		// Create log directory if needed
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, nil, fmt.Errorf("create log dir: %w", err)
		}

		f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, nil, fmt.Errorf("open log file: %w", err)
		}
		out = f
		closer = f
	}

	// Determine log level
	var level slog.Level
	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Create handler
	var handler slog.Handler
	if logFile != "" {
		// Use JSON for file logging
		handler = slog.NewJSONHandler(out, &slog.HandlerOptions{Level: level})
	} else {
		// Use text for stderr
		handler = slog.NewTextHandler(out, &slog.HandlerOptions{Level: level})
	}

	logger := slog.New(handler)
	return logger, closer, nil
}
