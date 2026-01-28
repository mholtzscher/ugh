package daemon

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/urfave/cli/v2"
)

func logsCommand() *cli.Command {
	return &cli.Command{
		Name:  "logs",
		Usage: "Show daemon logs",
		Description: "Show logs from the daemon service.\n\n" +
			"On Linux: Uses journalctl to show systemd service logs.\n" +
			"On macOS: Tails the log file specified in config.",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "lines",
				Aliases: []string{"n"},
				Value:   50,
				Usage:   "Number of lines to show",
			},
			&cli.BoolFlag{
				Name:  "no-follow",
				Usage: "Don't follow log output",
			},
		},
		Action: func(c *cli.Context) error {
			mgr, err := getServiceManager()
			if err != nil {
				return fmt.Errorf("detect service manager: %w", err)
			}

			w := deps.OutputWriter(c)

			ctx, cancel := context.WithCancel(c.Context)
			defer cancel()

			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, os.Interrupt)
			defer signal.Stop(sigCh)
			go func() {
				<-sigCh
				cancel()
			}()

			follow := !c.Bool("no-follow")
			if err := mgr.TailLogs(ctx, follow, c.Int("lines"), w.Out); err != nil {
				if ctx.Err() != nil {
					return nil
				}
				return fmt.Errorf("tail logs: %w", err)
			}

			return nil
		},
	}
}
