package daemon

import (
	"errors"

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
		// TODO: Implement in Phase 1
		return errors.New("daemon logs: not implemented yet")
	},
}

func init() {
	logsCmd.Flags().IntVarP(&logsOpts.Lines, "lines", "n", 50, "number of lines to show")
	logsCmd.Flags().BoolVar(&logsOpts.NoFollow, "no-follow", false, "don't follow log output")
}
