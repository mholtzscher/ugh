package daemon

import (
	"errors"
	"fmt"

	"github.com/mholtzscher/ugh/internal/daemon/service"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the daemon system service",
	Long: `Uninstall the daemon system service.

Stops the service if running, then removes the service configuration.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		mgr, err := getServiceManager()
		if err != nil {
			return fmt.Errorf("detect service manager: %w", err)
		}

		// Get the service path before uninstalling
		status, _ := mgr.Status()
		servicePath := status.ServicePath

		if err := mgr.Uninstall(); err != nil {
			if errors.Is(err, service.ErrNotInstalled) {
				return errors.New("service is not installed")
			}
			return fmt.Errorf("uninstall service: %w", err)
		}

		w := deps.OutputWriter()
		_, _ = fmt.Fprintln(w.Out, "Service uninstalled from", servicePath)
		return nil
	},
}
