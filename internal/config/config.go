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
	DefaultVersion = 1
)

type DB struct {
	Path      string `toml:"path"`
	SyncURL   string `toml:"sync_url"`
	AuthToken string `toml:"auth_token"`
}

type Config struct {
	Version int `toml:"version"`
	DB      DB  `toml:"db"`
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

	if err := toml.NewEncoder(file).Encode(cfg); err != nil {
		_ = file.Close()
		return fmt.Errorf("encode config: %w", err)
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

func Load(path string, allowMissing bool) (*LoadResult, error) {
	if path == "" {
		defaultPath, err := DefaultPath()
		if err != nil {
			return nil, fmt.Errorf("get default config path: %w", err)
		}
		path = defaultPath
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if allowMissing {
				return &LoadResult{Config: Config{Version: DefaultVersion}, WasLoaded: false}, nil
			}
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if _, err := toml.Decode(string(data), &cfg); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalid, err)
	}

	if cfg.Version == 0 {
		cfg.Version = DefaultVersion
	}

	return &LoadResult{
		Config:    cfg,
		UsedPath:  path,
		WasLoaded: true,
	}, nil
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
	} else if path[:2] == "~/" {
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
