package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/mholtzscher/ugh/internal/daemon/service"
)

//nolint:gochecknoglobals // CLI command definitions are package-level by design.
var daemonRestartCmd = &cli.Command{
	Name:  "restart",
	Usage: "Restart the daemon service",
	Action: func(_ context.Context, _ *cli.Command) error {
		mgr, err := getServiceManager()
		if err != nil {
			return fmt.Errorf("detect service manager: %w", err)
		}

		// Stop first (ignore error if not running)
		_ = mgr.Stop()

		err = mgr.Start()
		if err != nil {
			if errors.Is(err, service.ErrNotInstalled) {
				return errors.New("service not installed - run 'ugh daemon install' first")
			}
			return fmt.Errorf("start service: %w", err)
		}

		w := outputWriter()
		_, _ = fmt.Fprintln(w.Out, "Daemon restarted")
		return nil
	},
}
