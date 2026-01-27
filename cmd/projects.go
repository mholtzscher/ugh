package cmd

import (
	"context"

	"github.com/mholtzscher/ugh/internal/service"
	"github.com/spf13/cobra"
)

var projectsOpts struct {
	All      bool
	DoneOnly bool
	TodoOnly bool
	Counts   bool
}

var projectsCmd = &cobra.Command{
	Use:     "projects",
	Aliases: []string{"proj"},
	Short:   "List available project tags",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		svc, err := newService(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = svc.Close() }()

		tags, err := svc.ListProjects(ctx, service.ListTagsRequest{
			All:      projectsOpts.All,
			DoneOnly: projectsOpts.DoneOnly,
			TodoOnly: projectsOpts.TodoOnly,
		})
		if err != nil {
			return err
		}

		writer := outputWriter()
		if projectsOpts.Counts {
			return writer.WriteTagsWithCounts(tags)
		}
		return writer.WriteTags(tags)
	},
}

func init() {
	projectsCmd.Flags().BoolVarP(&projectsOpts.All, "all", "a", false, "include completed tasks")
	projectsCmd.Flags().BoolVarP(&projectsOpts.DoneOnly, "done", "x", false, "only completed tasks")
	projectsCmd.Flags().BoolVarP(&projectsOpts.TodoOnly, "todo", "t", false, "only pending tasks")
	projectsCmd.Flags().BoolVarP(&projectsOpts.Counts, "counts", "", false, "include counts")
}
