package daemon

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/mholtzscher/ugh/internal/flags"
	"github.com/urfave/cli/v3"
)

var logsCmd = &cli.Command{
	Name:  "logs",
	Usage: "Show daemon logs",
	Description: `Show logs from the daemon service.

On Linux: Uses journalctl to show systemd service logs.
On macOS: Tails the log file specified in config.`,
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    flags.FlagLines,
			Aliases: []string{"n"},
			Usage:   "number of lines to show",
			Value:   50,
		},
		&cli.BoolFlag{
			Name:  flags.FlagNoFollow,
			Usage: "don't follow log output",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		mgr, err := getServiceManager()
		if err != nil {
			return fmt.Errorf("detect service manager: %w", err)
		}

		w := deps.OutputWriter()

		// Create a context that can be cancelled
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		// Handle Ctrl+C gracefully
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)
		defer signal.Stop(sigCh)
		go func() {
			<-sigCh
			cancel()
		}()

		follow := !cmd.Bool(flags.FlagNoFollow)
		if err := mgr.TailLogs(ctx, follow, cmd.Int(flags.FlagLines), w.Out); err != nil {
			// Ignore context cancelled errors (user pressed Ctrl+C)
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("tail logs: %w", err)
		}

		return nil
	},
}
