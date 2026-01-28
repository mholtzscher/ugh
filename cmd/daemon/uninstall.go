package daemon

import (
	"errors"
	"fmt"

	"github.com/mholtzscher/ugh/internal/daemon/service"

	"github.com/urfave/cli/v2"
)

func uninstallCommand() *cli.Command {
	return &cli.Command{
		Name:  "uninstall",
		Usage: "Uninstall the daemon service",
		Action: func(c *cli.Context) error {
			mgr, err := getServiceManager()
			if err != nil {
				return fmt.Errorf("detect service manager: %w", err)
			}

			if err := mgr.Uninstall(); err != nil {
				if errors.Is(err, service.ErrNotInstalled) {
					return errors.New("service not installed")
				}
				return fmt.Errorf("uninstall service: %w", err)
			}

			w := deps.OutputWriter(c)
			_, _ = fmt.Fprintln(w.Out, "Service uninstalled")
			return nil
		},
	}
}
