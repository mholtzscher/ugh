package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/mholtzscher/ugh/internal/daemon/service"

	"github.com/urfave/cli/v3"
)

var daemonUninstallCmd = &cli.Command{
	Name:  "uninstall",
	Usage: "Uninstall the daemon system service",
	Description: `Uninstall the daemon system service.

Stops the service if running, then removes the service configuration.`,
	Action: func(ctx context.Context, cmd *cli.Command) error {
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

		w := outputWriter()
		_, _ = fmt.Fprintln(w.Out, "Service uninstalled from", servicePath)
		return nil
	},
}
