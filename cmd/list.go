package cmd

import (
	"context"

	"github.com/mholtzscher/ugh/internal/store"
	"github.com/spf13/cobra"
)

var listOpts struct {
	All      bool
	DoneOnly bool
	TodoOnly bool
	Project  string
	Context  string
	Priority string
	Search   string
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		filters := store.Filters{
			All:      listOpts.All,
			DoneOnly: listOpts.DoneOnly,
			TodoOnly: listOpts.TodoOnly,
			Project:  listOpts.Project,
			Context:  listOpts.Context,
			Priority: listOpts.Priority,
			Search:   listOpts.Search,
		}
		if !filters.All && !filters.DoneOnly && !filters.TodoOnly {
			filters.TodoOnly = true
		}

		st, err := openStore(ctx)
		if err != nil {
			return err
		}
		defer st.Close()

		tasks, err := st.ListTasks(ctx, filters)
		if err != nil {
			return err
		}

		writer := outputWriter()
		return writer.WriteTasks(tasks)
	},
}

func init() {
	listCmd.Flags().BoolVar(&listOpts.All, "all", false, "include completed tasks")
	listCmd.Flags().BoolVar(&listOpts.DoneOnly, "done", false, "only completed tasks")
	listCmd.Flags().BoolVar(&listOpts.TodoOnly, "todo", false, "only pending tasks")
	listCmd.Flags().StringVar(&listOpts.Project, "project", "", "filter by project")
	listCmd.Flags().StringVar(&listOpts.Context, "context", "", "filter by context")
	listCmd.Flags().StringVar(&listOpts.Priority, "priority", "", "filter by priority")
	listCmd.Flags().StringVar(&listOpts.Search, "search", "", "search text")
}
