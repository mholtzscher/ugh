package daemon

import (
	"errors"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the daemon service",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Implement in Phase 1
		return errors.New("daemon stop: not implemented yet")
	},
}
