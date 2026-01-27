package daemon

import (
	"errors"
	"fmt"

	"github.com/mholtzscher/ugh/internal/daemon/service"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the daemon service",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		mgr, err := getServiceManager()
		if err != nil {
			return fmt.Errorf("detect service manager: %w", err)
		}

		if err := mgr.Stop(); err != nil {
			if errors.Is(err, service.ErrNotRunning) {
				return errors.New("daemon is not running")
			}
			return fmt.Errorf("stop service: %w", err)
		}

		w := deps.OutputWriter()
		_, _ = fmt.Fprintln(w.Out, "Daemon stopped")
		return nil
	},
}
