package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

const (
	DefaultVersion    = 1
	EnvPrefix         = "UGH"
	DefaultDBFilename = "ugh.sqlite"
)

type DB struct {
	Path        string `mapstructure:"path" toml:"path"`
	SyncURL     string `mapstructure:"sync_url" toml:"sync_url"`
	AuthToken   string `mapstructure:"auth_token" toml:"auth_token"`
	SyncOnWrite bool   `mapstructure:"sync_on_write" toml:"sync_on_write"`
}

type Daemon struct {
	PeriodicSync     string `mapstructure:"periodic_sync" toml:"periodic_sync"`
	LogFile          string `mapstructure:"log_file" toml:"log_file"`
	LogLevel         string `mapstructure:"log_level" toml:"log_level"`
	SyncRetryMax     int    `mapstructure:"sync_retry_max" toml:"sync_retry_max"`
	SyncRetryBackoff string `mapstructure:"sync_retry_backoff" toml:"sync_retry_backoff"`
}

type Config struct {
	Version int    `mapstructure:"version" toml:"version"`
	DB      DB     `mapstructure:"db" toml:"db"`
	Daemon  Daemon `mapstructure:"daemon" toml:"daemon"`
}

type LoadResult struct {
	Config    Config
	UsedPath  string
	WasLoaded bool
}

var (
	ErrNotFound    = errors.New("config file not found")
	ErrInvalid     = errors.New("invalid config file")
	ErrMissingDB   = errors.New("db.path is required in config")
	ErrNotAbs      = errors.New("db.path must be an absolute path")
	ErrOutsideHome = errors.New("db.path outside home directory not allowed")
)

func DefaultPath() (string, error) {
	configDir, err := userConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "ugh", "config.toml"), nil
}

func Load(path string, allowMissing bool) (*LoadResult, error) {
	usedPath := path
	if usedPath == "" {
		defaultPath, err := DefaultPath()
		if err != nil {
			return nil, fmt.Errorf("get default config path: %w", err)
		}
		usedPath = defaultPath
	}

	var cfg Config

	if _, err := os.Stat(usedPath); err != nil {
		if os.IsNotExist(err) {
			if allowMissing {
				cfg = defaultConfig()
				return &LoadResult{
					Config:    cfg,
					UsedPath:  usedPath,
					WasLoaded: false,
				}, nil
			}
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalid, err)
	}

	data, err := os.ReadFile(usedPath)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalid, err)
	}

	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalid, err)
	}

	applyDefaults(&cfg)

	return &LoadResult{
		Config:    cfg,
		UsedPath:  usedPath,
		WasLoaded: true,
	}, nil
}

func defaultConfig() Config {
	cfg := Config{}
	applyDefaults(&cfg)
	return cfg
}

func applyDefaults(cfg *Config) {
	if cfg.Version == 0 {
		cfg.Version = DefaultVersion
	}
	if cfg.Daemon.PeriodicSync == "" {
		cfg.Daemon.PeriodicSync = "5m"
	}
	if cfg.Daemon.LogLevel == "" {
		cfg.Daemon.LogLevel = "info"
	}
	if cfg.Daemon.SyncRetryMax == 0 {
		cfg.Daemon.SyncRetryMax = 3
	}
	if cfg.Daemon.SyncRetryBackoff == "" {
		cfg.Daemon.SyncRetryBackoff = "1s"
	}
}

func Save(path string, cfg Config) error {
	if path == "" {
		return errors.New("config path is empty")
	}

	configDir := filepath.Dir(path)
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	file, err := os.CreateTemp(configDir, "config-*.toml")
	if err != nil {
		return fmt.Errorf("create temp config: %w", err)
	}
	defer func() { _ = os.Remove(file.Name()) }()

	enc := toml.NewEncoder(file)
	if err := enc.Encode(cfg); err != nil {
		_ = file.Close()
		return fmt.Errorf("write config: %w", err)
	}

	if err := file.Chmod(0o644); err != nil {
		_ = file.Close()
		return fmt.Errorf("chmod config: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close temp config: %w", err)
	}

	if err := os.Rename(file.Name(), path); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}

func ResolveDBPath(cfgPath, dbPath string) (string, error) {
	path := dbPath
	if path == "" {
		return "", errors.New("db path is empty")
	}

	path = filepath.Clean(path)
	path = os.ExpandEnv(path)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}

	if path == "~" {
		path = homeDir
	} else if len(path) >= 2 && path[:2] == "~/" {
		path = filepath.Join(homeDir, path[2:])
	}

	if !filepath.IsAbs(path) {
		if cfgPath != "" {
			cfgDir := filepath.Dir(cfgPath)
			path = filepath.Join(cfgDir, path)
		} else {
			abs, err := filepath.Abs(path)
			if err != nil {
				return "", fmt.Errorf("resolve relative path: %w", err)
			}
			path = abs
		}
	}

	path = filepath.Clean(path)

	return filepath.Abs(path)
}

func UserDataDir() (string, error) {
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

func DefaultDBPath() (string, error) {
	dataDir, err := UserDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dataDir, "ugh", DefaultDBFilename), nil
}

func EffectiveDBPath(cfg *Config, cfgPath string) (string, error) {
	if cfg != nil && cfg.DB.Path != "" {
		return ResolveDBPath(cfgPath, cfg.DB.Path)
	}
	return DefaultDBPath()
}

func userConfigDir() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		configDir, err := os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("user config dir: %w", err)
		}
		return configDir, nil
	case "windows":
		if appData := os.Getenv("APPDATA"); appData != "" {
			return appData, nil
		}
		configDir, err := os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("user config dir: %w", err)
		}
		return configDir, nil
	default:
		if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
			return xdgConfigHome, nil
		}
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("home dir: %w", err)
		}
		return filepath.Join(homeDir, ".config"), nil
	}
}
