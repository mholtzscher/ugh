package cmd

import (
	"context"

	"github.com/mholtzscher/ugh/internal/output"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm <id...>",
	Short: "Delete tasks",
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
		defer svc.Close()

		count, err := svc.DeleteTasks(ctx, ids)
		if err != nil {
			return err
		}
		writer := outputWriter()
		return writer.WriteSummary(output.Summary{Action: "rm", Count: count, IDs: ids})
	},
}
