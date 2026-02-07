package daemon

import (
	"time"

	"github.com/mholtzscher/ugh/internal/config"
)

const (
	defaultPeriodicSyncMinutes = 5
	defaultSyncRetryMax        = 3
)

// Config holds parsed daemon configuration with defaults applied.
type Config struct {
	PeriodicSync     time.Duration
	LogFile          string
	LogLevel         string
	SyncRetryMax     int
	SyncRetryBackoff time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		PeriodicSync:     defaultPeriodicSyncMinutes * time.Minute,
		LogFile:          "",
		LogLevel:         "info",
		SyncRetryMax:     defaultSyncRetryMax,
		SyncRetryBackoff: 1 * time.Second,
	}
}

// ParseConfig creates a daemon Config from a config.Daemon, applying defaults.
func ParseConfig(d config.Daemon) Config {
	cfg := DefaultConfig()

	if d.PeriodicSync != "" {
		if dur, err := time.ParseDuration(d.PeriodicSync); err == nil {
			cfg.PeriodicSync = dur
		}
	}

	if d.LogFile != "" {
		cfg.LogFile = d.LogFile
	}

	if d.LogLevel != "" {
		cfg.LogLevel = d.LogLevel
	}

	if d.SyncRetryMax > 0 {
		cfg.SyncRetryMax = d.SyncRetryMax
	}

	if d.SyncRetryBackoff != "" {
		if dur, err := time.ParseDuration(d.SyncRetryBackoff); err == nil {
			cfg.SyncRetryBackoff = dur
		}
	}

	return cfg
}
