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
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "List tasks",
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
		defer func() { _ = st.Close() }()

		tasks, err := st.ListTasks(ctx, filters)
		if err != nil {
			return err
		}

		writer := outputWriter()
		return writer.WriteTasks(tasks)
	},
}

func init() {
	listCmd.Flags().BoolVarP(&listOpts.All, "all", "a", false, "include completed tasks")
	listCmd.Flags().BoolVarP(&listOpts.DoneOnly, "done", "x", false, "only completed tasks")
	listCmd.Flags().BoolVarP(&listOpts.TodoOnly, "todo", "t", false, "only pending tasks")
	listCmd.Flags().StringVarP(&listOpts.Project, "project", "P", "", "filter by project")
	listCmd.Flags().StringVarP(&listOpts.Context, "context", "c", "", "filter by context")
	listCmd.Flags().StringVarP(&listOpts.Priority, "priority", "p", "", "filter by priority")
	listCmd.Flags().StringVarP(&listOpts.Search, "search", "s", "", "search text")
}
