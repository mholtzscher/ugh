package daemon

import (
	"context"
	"errors"
	"fmt"

	"github.com/mholtzscher/ugh/internal/daemon/service"

	"github.com/urfave/cli/v3"
)

var startCmd = &cli.Command{
	Name:  "start",
	Usage: "Start the daemon service",
	Description: `Start the daemon via the system service manager.

The service must be installed first with 'ugh daemon install'.`,
	Action: func(ctx context.Context, cmd *cli.Command) error {
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
