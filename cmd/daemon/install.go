package daemon

import (
	"errors"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the daemon as a system service",
	Long: `Install the daemon as a user-level system service.

On Linux (systemd):
  Creates ~/.config/systemd/user/ughd.service and enables it.

On macOS (launchd):
  Creates ~/Library/LaunchAgents/com.ugh.daemon.plist and loads it.

After installation, use 'ugh daemon start' to start the service.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Implement in Phase 1
		return errors.New("daemon install: not implemented yet")
	},
}
