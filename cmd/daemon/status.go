package daemon

import (
	"errors"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show daemon service status",
	Long: `Show the status of the daemon service.

Displays whether the service is installed, running, and if running,
shows uptime and sync status from the daemon's health endpoint.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Implement in Phase 1
		return errors.New("daemon status: not implemented yet")
	},
}
