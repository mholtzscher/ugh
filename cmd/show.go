package cmd

import (
	"context"
	"errors"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:     "show <id>",
	Aliases: []string{"s"},
	Short:   "Show a task",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("show requires a task id")
		}
		return nil
	},
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
		defer func() { _ = st.Close() }()

		task, err := st.GetTask(ctx, ids[0])
		if err != nil {
			return err
		}

		writer := outputWriter()
		return writer.WriteTask(task)
	},
}
