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
	taskscmd "github.com/mholtzscher/ugh/cmd/tasks"
	"github.com/mholtzscher/ugh/internal/config"
	"github.com/mholtzscher/ugh/internal/output"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/mholtzscher/ugh/internal/store"

	"github.com/urfave/cli/v2"
)

var (
	loadedConfig     *config.Config
	loadedConfigWas  bool
	loadedConfigPath string
)

const defaultTimeout = 30 * time.Second

func Run() error {
	app := &cli.App{
		Name:                 "ugh",
		Usage:                "todo.txt-inspired task CLI",
		Description:          "A CLI task manager using todo.txt format with libSQL storage",
		EnableBashCompletion: true,
		Flags:                globalFlags(false),
		Before: func(c *cli.Context) error {
			return loadConfig(c)
		},
		Action: func(c *cli.Context) error {
			return cli.ShowAppHelp(c)
		},
		Commands: append(taskCommands(),
			withGlobalFlags(setCategory(configcmd.Command(configcmd.Deps{
				Config:             func() *config.Config { return loadedConfig },
				SetConfig:          func(c *config.Config) { loadedConfig = c },
				ConfigWasLoaded:    func() bool { return loadedConfigWas },
				SetConfigWasLoaded: func(b bool) { loadedConfigWas = b },
				OutputWriter:       outputWriter,
				ConfigPath:         func(c *cli.Context) string { return flagString(c, "config") },
				DefaultDBPath:      config.DefaultDBPath,
			}), "Admin")),
			withGlobalFlags(setCategory(daemoncmd.Command(daemoncmd.Deps{
				Config:       func() *config.Config { return loadedConfig },
				ConfigPath:   func() string { return loadedConfigPath },
				OutputWriter: outputWriter,
			}), "Admin")),
		),
	}

	return app.Run(normalizeArgs(app, os.Args))
}

func taskCommands() []*cli.Command {
	deps := taskscmd.Deps{
		WithTimeout:          withTimeout,
		NewService:           func(c *cli.Context) (service.Service, error) { return newService(c.Context, c) },
		ParseIDs:             parseIDs,
		OutputWriter:         outputWriter,
		MaybeSyncBeforeWrite: maybeSyncBeforeWrite,
		MaybeSyncAfterWrite:  maybeSyncAfterWrite,
		FlagString:           flagString,
		FlagBool:             flagBool,
		OpenStore:            openStore,
	}
	taskscmd.Init(deps)
	commands := taskscmd.Commands()
	for i, cmd := range commands {
		commands[i] = withGlobalFlags(cmd)
	}
	return commands
}

func setCategory(cmd *cli.Command, category string) *cli.Command {
	if cmd == nil {
		return cmd
	}
	cmd.Category = category
	return cmd
}

func loadConfig(c *cli.Context) error {
	configPath := flagString(c, "config")
	allowMissing := configPath == "" || allowMissingConfig(c)

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
	loadedConfigPath = result.UsedPath

	return nil
}

func allowMissingConfig(c *cli.Context) bool {
	if c == nil {
		return false
	}
	args := c.Args().Slice()
	return len(args) >= 2 && args[0] == "config" && (args[1] == "init" || args[1] == "set")
}

func effectiveDBPath(c *cli.Context) (string, error) {
	if dbPath := flagString(c, "db"); dbPath != "" {
		return dbPath, nil
	}

	var cfgPath string
	if configPath := flagString(c, "config"); configPath != "" {
		cfgPath = configPath
	} else if loadedConfigWas {
		defaultPath, err := config.DefaultPath()
		if err != nil {
			return "", fmt.Errorf("get default config path: %w", err)
		}
		cfgPath = defaultPath
	}

	return config.EffectiveDBPath(loadedConfig, cfgPath)
}

func openStore(ctx context.Context, c *cli.Context) (*store.Store, error) {
	path, err := effectiveDBPath(c)
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

	var st *store.Store
	maxRetries := 5
	backoff := 100 * time.Millisecond
	for i := 0; i < maxRetries; i++ {
		st, err = store.Open(ctx, opts)
		if err == nil {
			return st, nil
		}
		if !isLockingError(err) {
			return nil, err
		}
		if i < maxRetries-1 {
			time.Sleep(backoff)
			backoff *= 2
		}
	}
	return nil, fmt.Errorf("%w (is the daemon running? try 'ugh daemon stop' or use the HTTP API)", err)
}

func isLockingError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "locked") || strings.Contains(errStr, "Locking")
}

func outputWriter(c *cli.Context) output.Writer {
	return output.NewWriter(flagBool(c, "json"), flagBool(c, "no-color"))
}

func withTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, defaultTimeout)
}
