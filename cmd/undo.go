package cmd

import (
	"context"
	"fmt"

	"github.com/mholtzscher/ugh/internal/output"
	"github.com/urfave/cli/v3"
)

var undoCmd = &cli.Command{
	Name:      "undo",
	Aliases:   []string{"u"},
	Usage:     "Mark tasks as not done",
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

		if err := maybeSyncBeforeWrite(ctx, svc); err != nil {
			return fmt.Errorf("sync pull: %w", err)
		}

		count, err := svc.SetDone(ctx, ids, false)
		if err != nil {
			return err
		}
		if err := maybeSyncAfterWrite(ctx, svc); err != nil {
			return fmt.Errorf("sync push: %w", err)
		}
		writer := outputWriter()
		return writer.WriteSummary(output.Summary{Action: "undo", Count: count, IDs: ids})
	},
}
