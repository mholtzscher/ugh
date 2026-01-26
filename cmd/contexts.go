package cmd

import (
	"context"

	"github.com/mholtzscher/ugh/internal/service"
	"github.com/spf13/cobra"
)

var contextsOpts struct {
	All      bool
	DoneOnly bool
	TodoOnly bool
	Counts   bool
}

var contextsCmd = &cobra.Command{
	Use:     "contexts",
	Aliases: []string{"ctx"},
	Short:   "List available context tags",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		svc, err := newTaskService(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = svc.Close() }()

		tags, err := svc.ListContexts(ctx, service.ListTagsRequest{
			All:      contextsOpts.All,
			DoneOnly: contextsOpts.DoneOnly,
			TodoOnly: contextsOpts.TodoOnly,
		})
		if err != nil {
			return err
		}

		writer := outputWriter()
		if contextsOpts.Counts {
			return writer.WriteTagsWithCounts(tags)
		}
		return writer.WriteTags(tags)
	},
}

func init() {
	contextsCmd.Flags().BoolVarP(&contextsOpts.All, "all", "a", false, "include completed tasks")
	contextsCmd.Flags().BoolVarP(&contextsOpts.DoneOnly, "done", "x", false, "only completed tasks")
	contextsCmd.Flags().BoolVarP(&contextsOpts.TodoOnly, "todo", "t", false, "only pending tasks")
	contextsCmd.Flags().BoolVarP(&contextsOpts.Counts, "counts", "", false, "include counts")
}
