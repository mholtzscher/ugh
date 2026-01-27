package cmd

import (
	"context"
	"fmt"

	"github.com/mholtzscher/ugh/internal/output"
	"github.com/spf13/cobra"
)

var undoCmd = &cobra.Command{
	Use:     "undo <id...>",
	Aliases: []string{"u"},
	Short:   "Mark tasks as not done",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		ids, err := parseIDs(args)
		if err != nil {
			return err
		}
		svc, err := newTaskService(ctx)
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
