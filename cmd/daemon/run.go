package daemon

import (
	"errors"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the daemon in foreground",
	Long: `Run the daemon server in the foreground.

This is primarily used by the system service manager (systemd/launchd).
For debugging, you can run it directly to see logs in the terminal.

The daemon provides:
  - HTTP API on the configured listen address (default: 127.0.0.1:9847)
  - Background sync to Turso cloud
  - File watcher for detecting CLI writes`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Implement in Phase 5
		return errors.New("daemon run: not implemented yet")
	},
}
