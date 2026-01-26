package cmd

import (
	"context"

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
		svc, err := newTaskService(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = svc.Close() }()

		count, err := svc.SetDone(ctx, ids, true)
		if err != nil {
			return err
		}
		writer := outputWriter()
		return writer.WriteSummary(output.Summary{Action: "done", Count: count, IDs: ids})
	},
}
