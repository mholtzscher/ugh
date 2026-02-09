package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/urfave/cli/v3"

	"github.com/mholtzscher/ugh/internal/output"
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
			err = w.WriteWarning("Service not installed")
			if err != nil {
				return err
			}
			return w.WriteInfo("Run 'ugh daemon install' to set up the service")
		}

		rows := []output.KeyValue{
			{Key: "service", Value: "installed"},
			{Key: "path", Value: status.ServicePath},
		}
		if status.Running {
			rows = append(rows, output.KeyValue{Key: "status", Value: "running"})
			if status.PID > 0 {
				rows = append(rows, output.KeyValue{Key: "pid", Value: strconv.Itoa(status.PID)})
			}
		} else {
			rows = append(rows, output.KeyValue{Key: "status", Value: "stopped"})
		}

		return w.WriteKeyValues(rows)
	},
}
