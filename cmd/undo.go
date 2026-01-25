package cmd

import (
	"context"

	"github.com/mholtzscher/ugh/internal/output"
	"github.com/spf13/cobra"
)

var undoCmd = &cobra.Command{
	Use:   "undo <id...>",
	Short: "Mark tasks as not done",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		ids, err := parseIDs(args)
		if err != nil {
			return err
		}
		st, err := openStore(ctx)
		if err != nil {
			return err
		}
		defer st.Close()

		count, err := st.SetDone(ctx, ids, false)
		if err != nil {
			return err
		}
		writer := outputWriter()
		return writer.WriteSummary(output.Summary{Action: "undo", Count: count, IDs: ids})
	},
}
