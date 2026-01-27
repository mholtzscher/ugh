package daemon

import (
	"errors"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the daemon system service",
	Long: `Uninstall the daemon system service.

Stops the service if running, then removes the service configuration.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Implement in Phase 1
		return errors.New("daemon uninstall: not implemented yet")
	},
}
