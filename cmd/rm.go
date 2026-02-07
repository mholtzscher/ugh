package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/mholtzscher/ugh/internal/output"
)

//nolint:gochecknoglobals // CLI command definitions are package-level by design.
var rmCmd = &cli.Command{
	Name:      "rm",
	Usage:     "Delete tasks",
	Category:  "Tasks",
	ArgsUsage: "<id...>",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		ids, err := parseIDs(commandArgs(cmd))
		if err != nil {
			return err
		}
		svc, err := newService(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = svc.Close() }()

		err = maybeSyncBeforeWrite(ctx, svc)
		if err != nil {
			return fmt.Errorf("sync pull: %w", err)
		}

		count, err := svc.DeleteTasks(ctx, ids)
		if err != nil {
			return err
		}
		err = maybeSyncAfterWrite(ctx, svc)
		if err != nil {
			return fmt.Errorf("sync push: %w", err)
		}
		writer := outputWriter()
		return writer.WriteSummary(output.Summary{Action: "rm", Count: count, IDs: ids})
	},
}
