package service

import (
	"context"
	"errors"
	"io"
	"runtime"
)

// ErrUnsupportedPlatform indicates the current platform is not supported.
var ErrUnsupportedPlatform = errors.New("unsupported platform")

// ErrNotInstalled indicates the service is not installed.
var ErrNotInstalled = errors.New("service not installed")

// ErrAlreadyInstalled indicates the service is already installed.
var ErrAlreadyInstalled = errors.New("service already installed")

// ErrNotRunning indicates the service is not running.
var ErrNotRunning = errors.New("service not running")

// Status represents the current state of the service.
type Status struct {
	Installed   bool
	Running     bool
	PID         int
	ServicePath string
}

// InstallConfig contains configuration for service installation.
type InstallConfig struct {
	BinaryPath string // Absolute path to ugh binary
	ConfigPath string // Path to config.toml (may be empty)
	Listen     string // HTTP listen address (default: 127.0.0.1:9847)
}

// Manager defines the interface for managing the daemon system service.
type Manager interface {
	// Name returns a human-readable name for this service manager (e.g., "systemd", "launchd").
	Name() string

	// Install generates and installs the service configuration.
	// Returns ErrAlreadyInstalled if the service is already installed.
	Install(cfg InstallConfig) error

	// Uninstall stops the service if running and removes the service configuration.
	// Returns ErrNotInstalled if the service is not installed.
	Uninstall() error

	// Start starts the service.
	// Returns ErrNotInstalled if the service is not installed.
	Start() error

	// Stop stops the service.
	// Returns ErrNotRunning if the service is not running.
	Stop() error

	// Status returns the current status of the service.
	Status() (Status, error)

	// LogPath returns the path to the log file (if applicable).
	// May return empty string if logs are managed by the service manager.
	LogPath() string

	// TailLogs streams logs from the service.
	// If follow is true, continues streaming until ctx is cancelled.
	// lines specifies how many historical lines to show.
	TailLogs(ctx context.Context, follow bool, lines int, w io.Writer) error
}

// Detect returns the appropriate ServiceManager for the current platform.
// Returns ErrUnsupportedPlatform if no suitable manager is available.
func Detect() (Manager, error) {
	switch runtime.GOOS {
	case "linux":
		return detectLinux()
	case "darwin":
		return NewLaunchd(), nil
	default:
		return nil, ErrUnsupportedPlatform
	}
}

// detectLinux checks for systemd on Linux.
func detectLinux() (Manager, error) {
	// Check if systemd is available by looking for systemctl
	sm := NewSystemd()
	if sm.isAvailable() {
		return sm, nil
	}
	return nil, errors.New("systemd not available - manual service setup required")
}
