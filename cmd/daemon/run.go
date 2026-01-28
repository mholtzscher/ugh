package daemon

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/mholtzscher/ugh/internal/config"
	"github.com/mholtzscher/ugh/internal/daemon"
	"github.com/mholtzscher/ugh/internal/store"

	"github.com/urfave/cli/v2"
)

func runCommand() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "Run the daemon in foreground",
		Description: "Run the daemon server in the foreground.\n\n" +
			"This is primarily used by the system service manager (systemd/launchd).\n" +
			"For debugging, you can run it directly to see logs in the terminal.\n\n" +
			"The daemon provides background sync to Turso cloud on a periodic interval.\n" +
			"It only opens the database when syncing, avoiding lock contention with the CLI.",
		Action: func(c *cli.Context) error {
			ctx := c.Context

			cfg := getConfig()
			if cfg == nil {
				return fmt.Errorf("config not loaded")
			}

			daemonCfg := daemon.ParseConfig(cfg.Daemon)

			logger, logCloser, err := setupLogging(daemonCfg.LogFile, daemonCfg.LogLevel)
			if err != nil {
				return fmt.Errorf("setup logging: %w", err)
			}
			if logCloser != nil {
				defer func() { _ = logCloser.Close() }()
			}

			cfgPath := getConfigPath()
			dbPath, err := config.EffectiveDBPath(cfg, cfgPath)
			if err != nil {
				return fmt.Errorf("get db path: %w", err)
			}

			dbDir := filepath.Dir(dbPath)
			if err := os.MkdirAll(dbDir, 0o755); err != nil {
				return fmt.Errorf("create db dir: %w", err)
			}

			cacheDir := filepath.Join(dbDir, ".cache")
			if os.Getenv("TURSO_GO_CACHE_DIR") == "" {
				_ = os.Setenv("TURSO_GO_CACHE_DIR", cacheDir)
			}

			storeOpts := store.Options{
				Path:        dbPath,
				SyncURL:     cfg.DB.SyncURL,
				AuthToken:   cfg.DB.AuthToken,
				BusyTimeout: 5000,
			}

			d := daemon.New(storeOpts, daemonCfg, logger)

			return d.Run(ctx)
		},
	}
}

// setupLogging sets up the logger based on config.
func setupLogging(logFile, logLevel string) (*slog.Logger, io.Closer, error) {
	var out io.Writer = os.Stderr
	var closer io.Closer

	if logFile != "" {
		path := logFile
		if len(path) >= 2 && path[:2] == "~/" {
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, nil, err
			}
			path = filepath.Join(home, path[2:])
		}

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

	var handler slog.Handler
	if logFile != "" {
		handler = slog.NewJSONHandler(out, &slog.HandlerOptions{Level: level})
	} else {
		handler = slog.NewTextHandler(out, &slog.HandlerOptions{Level: level})
	}

	logger := slog.New(handler)
	return logger, closer, nil
}
