package daemon

import (
	"os"
	"path/filepath"

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

var deps Deps

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
func getConfigPath() string {
	// Check if config was loaded
	cfg := deps.Config()
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
func getConfig() *config.Config {
	return deps.Config()
}

// Cmd is the parent command for all daemon subcommands.
var Cmd = &cli.Command{
	Name:  "daemon",
	Usage: "Manage the background daemon",
	Description: `Manage the ugh background daemon for HTTP API and Turso sync.

The daemon provides:
  - HTTP API for external integrations (Raycast, scripts, etc.)
  - Background sync to Turso cloud (debounced, periodic)

Use 'ugh daemon install' to set up the system service, then
'ugh daemon start' to start it.`,
}

// Register adds the daemon command and its subcommands to the parent command.
// Must be called with valid Deps before the command tree is executed.
func Register(parent *cli.Command, d Deps) {
	deps = d
	Cmd.Commands = []*cli.Command{
		installCmd,
		uninstallCmd,
		startCmd,
		stopCmd,
		restartCmd,
		statusCmd,
		logsCmd,
		runCmd,
	}
	parent.Commands = append(parent.Commands, Cmd)
}
