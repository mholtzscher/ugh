package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/mholtzscher/ugh/internal/daemon/service"
)

//nolint:gochecknoglobals // CLI command definitions are package-level by design.
var daemonStartCmd = &cli.Command{
	Name:  "start",
	Usage: "Start the daemon service",
	Description: `Start the daemon via the system service manager.

The service must be installed first with 'ugh daemon install'.`,
	Action: func(_ context.Context, _ *cli.Command) error {
		mgr, err := getServiceManager()
		if err != nil {
			return fmt.Errorf("detect service manager: %w", err)
		}

		err = mgr.Start()
		if err != nil {
			if errors.Is(err, service.ErrNotInstalled) {
				return errors.New("service not installed - run 'ugh daemon install' first")
			}
			return fmt.Errorf("start service: %w", err)
		}

		w := outputWriter()
		_, _ = fmt.Fprintln(w.Out, "Daemon started")
		return nil
	},
}
