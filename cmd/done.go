package cmd

import (
	"context"
	"fmt"

	"github.com/mholtzscher/ugh/internal/output"
	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:     "done <id...>",
	Aliases: []string{"d"},
	Short:   "Mark tasks as done",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		ids, err := parseIDs(args)
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

		count, err := svc.SetDone(ctx, ids, true)
		if err != nil {
			return err
		}
		if err := maybeSyncAfterWrite(ctx, svc); err != nil {
			return fmt.Errorf("sync push: %w", err)
		}
		writer := outputWriter()
		return writer.WriteSummary(output.Summary{Action: "done", Count: count, IDs: ids})
	},
}
