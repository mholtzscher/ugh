package daemon

import (
	"errors"
	"fmt"

	"github.com/mholtzscher/ugh/internal/daemon/service"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the daemon service",
	Long: `Start the daemon via the system service manager.

The service must be installed first with 'ugh daemon install'.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		mgr, err := getServiceManager()
		if err != nil {
			return fmt.Errorf("detect service manager: %w", err)
		}

		if err := mgr.Start(); err != nil {
			if errors.Is(err, service.ErrNotInstalled) {
				return errors.New("service not installed - run 'ugh daemon install' first")
			}
			return fmt.Errorf("start service: %w", err)
		}

		w := deps.OutputWriter()
		_, _ = fmt.Fprintln(w.Out, "Daemon started")
		return nil
	},
}
