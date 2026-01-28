package daemon

import (
	"errors"
	"fmt"

	"github.com/mholtzscher/ugh/internal/daemon/service"

	"github.com/urfave/cli/v2"
)

func stopCommand() *cli.Command {
	return &cli.Command{
		Name:  "stop",
		Usage: "Stop the daemon service",
		Action: func(c *cli.Context) error {
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

			w := deps.OutputWriter(c)
			_, _ = fmt.Fprintln(w.Out, "Daemon stopped")
			return nil
		},
	}
}
