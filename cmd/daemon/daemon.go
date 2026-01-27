package daemon

import (
	"github.com/mholtzscher/ugh/internal/config"
	"github.com/mholtzscher/ugh/internal/output"

	"github.com/spf13/cobra"
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

// Cmd is the parent command for all daemon subcommands.
var Cmd = &cobra.Command{
	Use:   "daemon",
	Short: "Manage the background daemon",
	Long: `Manage the ugh background daemon for HTTP API and Turso sync.

The daemon provides:
  - HTTP API for external integrations (Raycast, scripts, etc.)
  - Background sync to Turso cloud (debounced, periodic)

Use 'ugh daemon install' to set up the system service, then
'ugh daemon start' to start it.`,
}

// Register adds the daemon command and its subcommands to the parent command.
// Must be called with valid Deps before the command tree is executed.
func Register(parent *cobra.Command, d Deps) {
	deps = d
	Cmd.AddCommand(installCmd)
	Cmd.AddCommand(uninstallCmd)
	Cmd.AddCommand(startCmd)
	Cmd.AddCommand(stopCmd)
	Cmd.AddCommand(restartCmd)
	Cmd.AddCommand(statusCmd)
	Cmd.AddCommand(logsCmd)
	Cmd.AddCommand(runCmd)
	parent.AddCommand(Cmd)
}
