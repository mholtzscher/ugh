package cmd

import (
	"os"
	"path/filepath"

	"github.com/mholtzscher/ugh/internal/config"
	"github.com/mholtzscher/ugh/internal/daemon/service"

	"github.com/urfave/cli/v3"
)

// daemonCmd is the parent command for all daemon subcommands.
var daemonCmd = &cli.Command{
	Name:     "daemon",
	Usage:    "Manage the background daemon",
	Category: "System",
	Description: `Manage the ugh background daemon for HTTP API and Turso sync.

The daemon provides:
  - HTTP API for external integrations (Raycast, scripts, etc.)
  - Background sync to Turso cloud (debounced, periodic)

Use 'ugh daemon install' to set up the system service, then
'ugh daemon start' to start it.`,
	Commands: []*cli.Command{
		daemonInstallCmd,
		daemonUninstallCmd,
		daemonStartCmd,
		daemonStopCmd,
		daemonRestartCmd,
		daemonStatusCmd,
		daemonLogsCmd,
		daemonRunCmd,
	},
}

// getServiceManager returns the appropriate service manager for the current platform.
func getServiceManager() (service.Manager, error) {
	return service.Detect()
}

// getBinaryPath returns the absolute path to the current executable.
func getBinaryPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(exe)
}

// getDaemonConfigPath returns the config path to use for the daemon.
// Returns the path if config was explicitly set or loaded from default location.
func getDaemonConfigPath() string {
	cfg := loadedConfig
	if cfg == nil {
		return ""
	}

	defaultPath, err := config.DefaultPath()
	if err != nil {
		return ""
	}

	if _, err := os.Stat(defaultPath); err != nil {
		return ""
	}

	return defaultPath
}
