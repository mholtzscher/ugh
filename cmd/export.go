package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/mholtzscher/ugh/internal/output"
	"github.com/mholtzscher/ugh/internal/store"
	"github.com/mholtzscher/ugh/internal/todotxt"
	"github.com/spf13/cobra"
)

var exportOpts struct {
	All      bool
	DoneOnly bool
	TodoOnly bool
	Project  string
	Context  string
	Priority string
	Search   string
}

var exportCmd = &cobra.Command{
	Use:     "export <path|->",
	Aliases: []string{"ex"},
	Short:   "Export tasks to todo.txt",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("export requires a file path or -")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		path := args[0]
		if rootOpts.JSON && path == "-" {
			return errors.New("--json cannot be used when exporting to stdout")
		}

		filters := store.Filters{
			All:      exportOpts.All,
			DoneOnly: exportOpts.DoneOnly,
			TodoOnly: exportOpts.TodoOnly,
			Project:  exportOpts.Project,
			Context:  exportOpts.Context,
			Priority: exportOpts.Priority,
			Search:   exportOpts.Search,
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

		var file *os.File
		if path == "-" {
			file = os.Stdout
		} else {
			created, err := os.Create(path)
			if err != nil {
				return fmt.Errorf("create file: %w", err)
			}
			defer func() { _ = created.Close() }()
			file = created
		}

		writer := bufio.NewWriter(file)
		for _, task := range tasks {
			line := todotxt.Format(todotxt.Parsed{
				Done:           task.Done,
				Priority:       task.Priority,
				CompletionDate: task.CompletionDate,
				CreationDate:   task.CreationDate,
				Description:    task.Description,
				Projects:       task.Projects,
				Contexts:       task.Contexts,
				Meta:           task.Meta,
				Unknown:        task.Unknown,
			})
			if _, err := fmt.Fprintln(writer, line); err != nil {
				return err
			}
		}
		if err := writer.Flush(); err != nil {
			return err
		}

		out := outputWriter()
		return out.WriteSummary(output.ExportSummary{Action: "export", Count: int64(len(tasks)), File: path})
	},
}

func init() {
	exportCmd.Flags().BoolVarP(&exportOpts.All, "all", "a", false, "include completed tasks")
	exportCmd.Flags().BoolVarP(&exportOpts.DoneOnly, "done", "x", false, "only completed tasks")
	exportCmd.Flags().BoolVarP(&exportOpts.TodoOnly, "todo", "t", false, "only pending tasks")
	exportCmd.Flags().StringVarP(&exportOpts.Project, "project", "P", "", "filter by project")
	exportCmd.Flags().StringVarP(&exportOpts.Context, "context", "c", "", "filter by context")
	exportCmd.Flags().StringVarP(&exportOpts.Priority, "priority", "p", "", "filter by priority")
	exportCmd.Flags().StringVarP(&exportOpts.Search, "search", "s", "", "search text")
}
