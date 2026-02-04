package daemon

import (
	"context"
	"errors"
	"fmt"

	"github.com/mholtzscher/ugh/internal/daemon/service"

	"github.com/urfave/cli/v3"
)

var restartCmd = &cli.Command{
	Name:  "restart",
	Usage: "Restart the daemon service",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		mgr, err := getServiceManager()
		if err != nil {
			return fmt.Errorf("detect service manager: %w", err)
		}

		// Stop first (ignore error if not running)
		_ = mgr.Stop()

		if err := mgr.Start(); err != nil {
			if errors.Is(err, service.ErrNotInstalled) {
				return errors.New("service not installed - run 'ugh daemon install' first")
			}
			return fmt.Errorf("start service: %w", err)
		}

		w := deps.OutputWriter()
		_, _ = fmt.Fprintln(w.Out, "Daemon restarted")
		return nil
	},
}
