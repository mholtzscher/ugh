package daemon

import (
	"os"
	"path/filepath"

	"github.com/mholtzscher/ugh/internal/config"
	"github.com/mholtzscher/ugh/internal/daemon/service"
	"github.com/mholtzscher/ugh/internal/output"

	"github.com/urfave/cli/v2"
)

// Deps holds dependencies injected from the parent cmd package.
// This avoids circular imports between cmd and cmd/daemon.
type Deps struct {
	// Config returns the currently loaded config (may be nil).
	Config func() *config.Config
	// ConfigPath returns the resolved config path.
	ConfigPath func() string
	// OutputWriter returns the configured output writer.
	OutputWriter func(*cli.Context) output.Writer
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
	path := ""
	if deps.ConfigPath != nil {
		path = deps.ConfigPath()
	}
	if path == "" {
		defaultPath, err := config.DefaultPath()
		if err != nil {
			return ""
		}
		path = defaultPath
	}
	if _, err := os.Stat(path); err != nil {
		return ""
	}
	return path
}

// getConfig returns the loaded config, or nil if not loaded.
func getConfig() *config.Config {
	return deps.Config()
}

// Command is the parent command for all daemon subcommands.
func Command(d Deps) *cli.Command {
	deps = d
	return &cli.Command{
		Name:  "daemon",
		Usage: "Manage the background daemon",
		Description: "Manage the ugh background daemon for HTTP API and Turso sync.\n\n" +
			"The daemon provides:\n" +
			"  - HTTP API for external integrations (Raycast, scripts, etc.)\n" +
			"  - Background sync to Turso cloud (debounced, periodic)\n\n" +
			"Use 'ugh daemon install' to set up the system service, then\n" +
			"'ugh daemon start' to start it.",
		Subcommands: []*cli.Command{
			installCommand(),
			uninstallCommand(),
			startCommand(),
			stopCommand(),
			restartCommand(),
			statusCommand(),
			logsCommand(),
			runCommand(),
		},
	}
}
