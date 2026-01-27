package daemon

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
)

var logsOpts struct {
	Lines    int
	NoFollow bool
}

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show daemon logs",
	Long: `Show logs from the daemon service.

On Linux: Uses journalctl to show systemd service logs.
On macOS: Tails the log file specified in config.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		mgr, err := getServiceManager()
		if err != nil {
			return fmt.Errorf("detect service manager: %w", err)
		}

		w := deps.OutputWriter()

		// Create a context that can be cancelled
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Handle Ctrl+C gracefully
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)
		defer signal.Stop(sigCh)
		go func() {
			<-sigCh
			cancel()
		}()

		follow := !logsOpts.NoFollow
		if err := mgr.TailLogs(ctx, follow, logsOpts.Lines, w.Out); err != nil {
			// Ignore context cancelled errors (user pressed Ctrl+C)
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("tail logs: %w", err)
		}

		return nil
	},
}

func init() {
	logsCmd.Flags().IntVarP(&logsOpts.Lines, "lines", "n", 50, "number of lines to show")
	logsCmd.Flags().BoolVar(&logsOpts.NoFollow, "no-follow", false, "don't follow log output")
}
