package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/mholtzscher/ugh/internal/daemon/service"
)

//nolint:gochecknoglobals // CLI command definitions are package-level by design.
var daemonStopCmd = &cli.Command{
	Name:  "stop",
	Usage: "Stop the daemon service",
	Action: func(_ context.Context, _ *cli.Command) error {
		mgr, err := getServiceManager()
		if err != nil {
			return fmt.Errorf("detect service manager: %w", err)
		}

		err = mgr.Stop()
		if err != nil {
			if errors.Is(err, service.ErrNotRunning) {
				return errors.New("daemon is not running")
			}
			return fmt.Errorf("stop service: %w", err)
		}

		w := outputWriter()
		_, _ = fmt.Fprintln(w.Out, "Daemon stopped")
		return nil
	},
}
