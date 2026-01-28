package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
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

// Daemon holds daemon-specific configuration.
type Daemon struct {
	PeriodicSync     string `mapstructure:"periodic_sync" toml:"periodic_sync"`           // Background sync interval (default: "5m")
	LogFile          string `mapstructure:"log_file" toml:"log_file"`                     // Log file path (empty = stderr)
	LogLevel         string `mapstructure:"log_level" toml:"log_level"`                   // Log level: debug, info, warn, error (default: "info")
	SyncRetryMax     int    `mapstructure:"sync_retry_max" toml:"sync_retry_max"`         // Max sync retry attempts (default: 3)
	SyncRetryBackoff string `mapstructure:"sync_retry_backoff" toml:"sync_retry_backoff"` // Initial retry backoff (default: "1s")
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
	Viper     *viper.Viper // Expose Viper instance for flag binding and watching
}

var (
	ErrNotFound    = errors.New("config file not found")
	ErrInvalid     = errors.New("invalid config file")
	ErrMissingDB   = errors.New("db.path is required in config")
	ErrNotAbs      = errors.New("db.path must be an absolute path")
	ErrOutsideHome = errors.New("db.path outside home directory not allowed")
)

// DefaultPath returns the default config file path.
func DefaultPath() (string, error) {
	configDir, err := userConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "ugh", "config.toml"), nil
}

// NewViper creates a new Viper instance with defaults and env var binding.
func NewViper() *viper.Viper {
	v := viper.New()
	v.SetConfigType("toml")
	v.SetEnvPrefix(EnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Set defaults for all config keys to enable env var binding
	v.SetDefault("version", DefaultVersion)
	v.SetDefault("db.path", "")
	v.SetDefault("db.sync_url", "")
	v.SetDefault("db.auth_token", "")
	v.SetDefault("db.sync_on_write", false)
	v.SetDefault("daemon.periodic_sync", "5m")
	v.SetDefault("daemon.log_file", "")
	v.SetDefault("daemon.log_level", "info")
	v.SetDefault("daemon.sync_retry_max", 3)
	v.SetDefault("daemon.sync_retry_backoff", "1s")

	return v
}

// Load loads configuration from the given path (or default) using Viper.
// Environment variables with UGH_ prefix override file values.
// For example: UGH_DB_AUTH_TOKEN overrides db.auth_token.
func Load(path string, allowMissing bool) (*LoadResult, error) {
	v := NewViper()

	usedPath := path
	if usedPath == "" {
		defaultPath, err := DefaultPath()
		if err != nil {
			return nil, fmt.Errorf("get default config path: %w", err)
		}
		usedPath = defaultPath
	}

	v.SetConfigFile(usedPath)

	err := v.ReadInConfig()
	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		var pathError *os.PathError

		isNotFound := errors.As(err, &configFileNotFoundError) ||
			(errors.As(err, &pathError) && errors.Is(pathError.Err, os.ErrNotExist))

		if isNotFound {
			if allowMissing {
				// Return defaults (env vars still apply via AutomaticEnv)
				var cfg Config
				if unmarshalErr := v.Unmarshal(&cfg); unmarshalErr != nil {
					return nil, fmt.Errorf("unmarshal defaults: %w", unmarshalErr)
				}
				return &LoadResult{
					Config:    cfg,
					WasLoaded: false,
					Viper:     v,
				}, nil
			}
			return nil, ErrNotFound
		}

		// Parse error
		return nil, fmt.Errorf("%w: %v", ErrInvalid, err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalid, err)
	}

	if cfg.Version == 0 {
		cfg.Version = DefaultVersion
	}

	return &LoadResult{
		Config:    cfg,
		UsedPath:  usedPath,
		WasLoaded: true,
		Viper:     v,
	}, nil
}

// Save writes the config to the given path using atomic write.
func Save(path string, cfg Config) error {
	if path == "" {
		return errors.New("config path is empty")
	}

	configDir := filepath.Dir(path)
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	// Use Viper for consistent formatting
	v := NewViper()

	// Set all values from config struct using reflection
	v.Set("version", cfg.Version)
	v.Set("db.path", cfg.DB.Path)
	v.Set("db.sync_url", cfg.DB.SyncURL)
	v.Set("db.auth_token", cfg.DB.AuthToken)
	v.Set("db.sync_on_write", cfg.DB.SyncOnWrite)
	v.Set("daemon.periodic_sync", cfg.Daemon.PeriodicSync)
	v.Set("daemon.log_file", cfg.Daemon.LogFile)
	v.Set("daemon.log_level", cfg.Daemon.LogLevel)
	v.Set("daemon.sync_retry_max", cfg.Daemon.SyncRetryMax)
	v.Set("daemon.sync_retry_backoff", cfg.Daemon.SyncRetryBackoff)

	return writeConfigAtomically(v, path, configDir)
}

// writeConfigAtomically writes config to a temp file then renames it.
func writeConfigAtomically(v *viper.Viper, path, configDir string) error {
	file, err := os.CreateTemp(configDir, "config-*.toml")
	if err != nil {
		return fmt.Errorf("create temp config: %w", err)
	}
	defer func() { _ = os.Remove(file.Name()) }()

	v.SetConfigType("toml")
	if err := v.WriteConfigAs(file.Name()); err != nil {
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
		if errors.Is(err, os.ErrExist) {
			if removeErr := os.Remove(path); removeErr == nil {
				err = os.Rename(file.Name(), path)
			}
		}
		if err != nil {
			return fmt.Errorf("write config: %w", err)
		}
	}

	return nil
}

// Watch starts watching the config file for changes.
// The callback is called whenever the config file changes.
// Returns a function to stop watching.
func Watch(v *viper.Viper, callback func(Config)) func() {
	v.OnConfigChange(func(e fsnotify.Event) {
		var cfg Config
		if err := v.Unmarshal(&cfg); err != nil {
			return // Ignore invalid config changes
		}
		callback(cfg)
	})
	v.WatchConfig()

	return func() {
		v.OnConfigChange(nil)
	}
}

// ResolveDBPath resolves a database path, expanding ~ and environment variables.
// If the path is relative, it's resolved relative to the config file directory.
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

// UserDataDir returns the user data directory based on OS.
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

// DefaultDBPath returns the default database path.
func DefaultDBPath() (string, error) {
	dataDir, err := UserDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dataDir, "ugh", DefaultDBFilename), nil
}

// EffectiveDBPath returns the effective database path from config or default.
// It handles ~ expansion, environment variables, and relative paths.
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
