package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v3"
)

// daemonStatusOutput represents the status information for JSON output.
type daemonStatusOutput struct {
	Installed   bool   `json:"installed"`
	Running     bool   `json:"running"`
	PID         int    `json:"pid,omitempty"`
	ServicePath string `json:"servicePath,omitempty"`
}

//nolint:gochecknoglobals // CLI command definitions are package-level by design.
var daemonStatusCmd = &cli.Command{
	Name:  "status",
	Usage: "Show daemon service status",
	Description: `Show the status of the daemon service.

Displays whether the service is installed, running, and if running,
shows uptime and sync status from the daemon's health endpoint.`,
	Action: func(_ context.Context, _ *cli.Command) error {
		mgr, err := getServiceManager()
		if err != nil {
			return fmt.Errorf("detect service manager: %w", err)
		}

		status, err := mgr.Status()
		if err != nil {
			return fmt.Errorf("get status: %w", err)
		}

		w := outputWriter()
		if w.JSON {
			output := daemonStatusOutput{
				Installed:   status.Installed,
				Running:     status.Running,
				PID:         status.PID,
				ServicePath: status.ServicePath,
			}
			enc := json.NewEncoder(w.Out)
			return enc.Encode(output)
		}

		// Human-readable output
		if !status.Installed {
			_, _ = fmt.Fprintln(w.Out, "Service:  not installed")
			_, _ = fmt.Fprintln(w.Out, "Run 'ugh daemon install' to set up the service")
			return nil
		}

		_, _ = fmt.Fprintln(w.Out, "Service:  installed")
		_, _ = fmt.Fprintln(w.Out, "Path:    ", status.ServicePath)
		if status.Running {
			_, _ = fmt.Fprintln(w.Out, "Status:   running")
			if status.PID > 0 {
				_, _ = fmt.Fprintln(w.Out, "PID:     ", status.PID)
			}
		} else {
			_, _ = fmt.Fprintln(w.Out, "Status:   stopped")
		}

		return nil
	},
}
