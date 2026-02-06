package daemon

import (
	"os"
	"path/filepath"

	"github.com/mholtzscher/ugh/cmd/meta"
	"github.com/mholtzscher/ugh/cmd/registry"
	"github.com/mholtzscher/ugh/internal/config"
	"github.com/mholtzscher/ugh/internal/daemon/service"
	"github.com/mholtzscher/ugh/internal/output"

	"github.com/urfave/cli/v3"
)

// Deps holds dependencies injected from the parent cmd package.
// This avoids circular imports between cmd and cmd/daemon.
type Deps struct {
	// Config returns the currently loaded config (may be nil).
	Config func() *config.Config
	// OutputWriter returns the configured output writer.
	OutputWriter func() output.Writer
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

// getConfigPath returns the config path to use for the daemon.
// Returns the path if config was explicitly set or loaded from default location.
func getConfigPath(d Deps) string {
	// Check if config was loaded
	cfg := d.Config()
	if cfg == nil {
		return ""
	}

	// Return the default config path
	defaultPath, err := config.DefaultPath()
	if err != nil {
		return ""
	}

	// Check if the file exists
	if _, err := os.Stat(defaultPath); err != nil {
		return ""
	}

	return defaultPath
}

// getConfig returns the loaded config, or nil if not loaded.
func getConfig(d Deps) *config.Config {
	return d.Config()
}

const (
	daemonID          registry.ID = "daemon"
	daemonInstallID   registry.ID = "daemon.install"
	daemonUninstallID registry.ID = "daemon.uninstall"
	daemonStartID     registry.ID = "daemon.start"
	daemonStopID      registry.ID = "daemon.stop"
	daemonRestartID   registry.ID = "daemon.restart"
	daemonStatusID    registry.ID = "daemon.status"
	daemonLogsID      registry.ID = "daemon.logs"
	daemonRunID       registry.ID = "daemon.run"
)

// Register adds daemon command specs to the registry.
func Register(r *registry.Registry, d Deps) error {
	return r.AddAll(
		registry.Spec{ID: daemonID, Source: "cmd/daemon", Build: newDaemonCmd},
		registry.Spec{ID: daemonInstallID, ParentID: daemonID, Source: "cmd/daemon", Build: func() *cli.Command { return newInstallCmd(d) }},
		registry.Spec{ID: daemonUninstallID, ParentID: daemonID, Source: "cmd/daemon", Build: func() *cli.Command { return newUninstallCmd(d) }},
		registry.Spec{ID: daemonStartID, ParentID: daemonID, Source: "cmd/daemon", Build: func() *cli.Command { return newStartCmd(d) }},
		registry.Spec{ID: daemonStopID, ParentID: daemonID, Source: "cmd/daemon", Build: func() *cli.Command { return newStopCmd(d) }},
		registry.Spec{ID: daemonRestartID, ParentID: daemonID, Source: "cmd/daemon", Build: func() *cli.Command { return newRestartCmd(d) }},
		registry.Spec{ID: daemonStatusID, ParentID: daemonID, Source: "cmd/daemon", Build: func() *cli.Command { return newStatusCmd(d) }},
		registry.Spec{ID: daemonLogsID, ParentID: daemonID, Source: "cmd/daemon", Build: func() *cli.Command { return newLogsCmd(d) }},
		registry.Spec{ID: daemonRunID, ParentID: daemonID, Source: "cmd/daemon", Build: func() *cli.Command { return newRunCmd(d) }},
	)
}

func newDaemonCmd() *cli.Command {
	return &cli.Command{
		Name:     "daemon",
		Usage:    "Manage the background daemon",
		Category: meta.SystemCategory.String(),
		Description: `Manage the ugh background daemon for HTTP API and Turso sync.

The daemon provides:
  - HTTP API for external integrations (Raycast, scripts, etc.)
  - Background sync to Turso cloud (debounced, periodic)

Use 'ugh daemon install' to set up the system service, then
'ugh daemon start' to start it.`,
	}
}
