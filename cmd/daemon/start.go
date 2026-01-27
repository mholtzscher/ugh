package daemon

import (
	"errors"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the daemon service",
	Long: `Start the daemon via the system service manager.

The service must be installed first with 'ugh daemon install'.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Implement in Phase 1
		return errors.New("daemon start: not implemented yet")
	},
}
