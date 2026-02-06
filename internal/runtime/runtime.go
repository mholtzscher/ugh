package runtime

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/mholtzscher/ugh/internal/config"
	"github.com/mholtzscher/ugh/internal/output"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/mholtzscher/ugh/internal/store"
)

type State struct {
	ConfigPath string
	DBPath     string
	JSON       bool
	NoColor    bool
}

type Runtime struct {
	state           State
	loadedConfig    *config.Config
	loadedConfigWas bool
}

func New() *Runtime {
	return &Runtime{}
}

func (r *Runtime) SetGlobalOptions(configPath, dbPath string, jsonMode, noColor bool) {
	r.state.ConfigPath = configPath
	r.state.DBPath = dbPath
	r.state.JSON = jsonMode
	r.state.NoColor = noColor
}

func (r *Runtime) Config() *config.Config {
	return r.loadedConfig
}

func (r *Runtime) SetConfig(cfg *config.Config) {
	r.loadedConfig = cfg
}

func (r *Runtime) ConfigWasLoaded() bool {
	return r.loadedConfigWas
}

func (r *Runtime) SetConfigWasLoaded(wasLoaded bool) {
	r.loadedConfigWas = wasLoaded
}

func (r *Runtime) ConfigPath() string {
	return r.state.ConfigPath
}

func (r *Runtime) LoadConfig() error {
	result, err := config.Load(r.state.ConfigPath, true)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if !result.WasLoaded && r.shouldAutoInitConfig() {
		cfgPath := r.state.ConfigPath
		if cfgPath == "" {
			defaultPath, err := config.DefaultPath()
			if err != nil {
				return fmt.Errorf("get default config path: %w", err)
			}
			cfgPath = defaultPath
		}

		dbPath, err := DefaultDBPath()
		if err != nil {
			return err
		}

		result.Config = config.Config{
			Version: config.DefaultVersion,
			DB:      config.DB{Path: dbPath},
			Daemon: config.Daemon{
				PeriodicSync:     "5m",
				LogLevel:         "info",
				SyncRetryMax:     3,
				SyncRetryBackoff: "1s",
			},
		}

		if err := config.Save(cfgPath, result.Config); err != nil {
			return fmt.Errorf("save config: %w", err)
		}
		result.WasLoaded = true
	}

	r.loadedConfig = &result.Config
	r.loadedConfigWas = result.WasLoaded
	return nil
}

func (r *Runtime) shouldAutoInitConfig() bool {
	if r.state.ConfigPath != "" {
		return true
	}
	return r.state.DBPath == ""
}

func (r *Runtime) EffectiveDBPath() (string, error) {
	return effectiveDBPath(r.state.ConfigPath, r.state.DBPath, r.loadedConfig, r.loadedConfigWas)
}

func EffectiveDBPathForConfig(configPath string, cfg *config.Config) (string, error) {
	return effectiveDBPath(configPath, "", cfg, configPath != "")
}

func effectiveDBPath(configPath, overrideDBPath string, cfg *config.Config, cfgLoaded bool) (string, error) {
	if overrideDBPath != "" {
		return overrideDBPath, nil
	}
	if cfg != nil && cfg.DB.Path != "" {
		resolvedPath := configPath
		if resolvedPath == "" && cfgLoaded {
			defaultPath, err := config.DefaultPath()
			if err != nil {
				return "", fmt.Errorf("get default config path: %w", err)
			}
			resolvedPath = defaultPath
		}
		dbPath, err := config.ResolveDBPath(resolvedPath, cfg.DB.Path)
		if err != nil {
			return "", fmt.Errorf("resolve db path from config: %w", err)
		}
		return dbPath, nil
	}
	return DefaultDBPath()
}

func (r *Runtime) OpenStore(ctx context.Context) (*store.Store, error) {
	path, err := r.EffectiveDBPath()
	if err != nil {
		return nil, err
	}
	if err := PrepareDBPath(path); err != nil {
		return nil, err
	}

	opts := store.Options{Path: path}
	if r.loadedConfig != nil {
		opts.SyncURL = r.loadedConfig.DB.SyncURL
		opts.AuthToken = r.loadedConfig.DB.AuthToken
	}

	var st *store.Store
	maxRetries := 5
	backoff := 100 * time.Millisecond
	for i := range maxRetries {
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

func PrepareDBPath(path string) error {
	dbDir := filepath.Dir(path)
	if err := os.MkdirAll(dbDir, 0o755); err != nil {
		return fmt.Errorf("create db dir: %w", err)
	}
	cacheDir := filepath.Join(dbDir, ".cache")
	if os.Getenv("TURSO_GO_CACHE_DIR") == "" {
		_ = os.Setenv("TURSO_GO_CACHE_DIR", cacheDir)
	}
	return nil
}

func (r *Runtime) OutputWriter() output.Writer {
	return output.NewWriter(r.state.JSON, r.state.NoColor)
}

func (r *Runtime) NewService(ctx context.Context) (service.Service, error) {
	st, err := r.OpenStore(ctx)
	if err != nil {
		return nil, err
	}
	return service.NewTaskService(st), nil
}

func (r *Runtime) AutoSyncEnabled() bool {
	return r.loadedConfig != nil && r.loadedConfig.DB.SyncOnWrite && r.loadedConfig.DB.SyncURL != ""
}

func (r *Runtime) MaybeSyncBeforeWrite(ctx context.Context, svc service.Service) error {
	if !r.AutoSyncEnabled() {
		return nil
	}
	return svc.Sync(ctx)
}

func (r *Runtime) MaybeSyncAfterWrite(ctx context.Context, svc service.Service) error {
	if !r.AutoSyncEnabled() {
		return nil
	}
	return svc.Push(ctx)
}

func DefaultDBPath() (string, error) {
	dataDir, err := userDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dataDir, "ugh", "ugh.sqlite"), nil
}

func userDataDir() (string, error) {
	switch runtime.GOOS {
	case "darwin":
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

func isLockingError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "locked") || strings.Contains(errStr, "Locking")
}
