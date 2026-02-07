package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/mholtzscher/ugh/internal/daemon/service"
)

//nolint:gochecknoglobals // CLI command definitions are package-level by design.
var daemonUninstallCmd = &cli.Command{
	Name:  "uninstall",
	Usage: "Uninstall the daemon system service",
	Description: `Uninstall the daemon system service.

Stops the service if running, then removes the service configuration.`,
	Action: func(_ context.Context, _ *cli.Command) error {
		mgr, err := getServiceManager()
		if err != nil {
			return fmt.Errorf("detect service manager: %w", err)
		}

		// Get the service path before uninstalling
		status, _ := mgr.Status()
		servicePath := status.ServicePath

		err = mgr.Uninstall()
		if err != nil {
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
