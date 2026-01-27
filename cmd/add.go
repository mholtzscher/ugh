package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mholtzscher/ugh/internal/service"
	"github.com/spf13/cobra"
)

var addOpts struct {
	Priority  string
	Projects  []string
	Contexts  []string
	Meta      []string
	Done      bool
	Created   string
	Completed string
}

var addCmd = &cobra.Command{
	Use:     "add [todo.txt line]",
	Aliases: []string{"a"},
	Short:   "Add a task",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		line := strings.TrimSpace(strings.Join(args, " "))
		if line == "" {
			return errors.New("todo text required")
		}

		svc, err := newService(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = svc.Close() }()

		if err := maybeSyncBeforeWrite(ctx, svc); err != nil {
			return fmt.Errorf("sync pull: %w", err)
		}

		task, err := svc.CreateTask(ctx, service.CreateTaskRequest{
			Line:      line,
			Priority:  addOpts.Priority,
			Projects:  addOpts.Projects,
			Contexts:  addOpts.Contexts,
			Meta:      addOpts.Meta,
			Done:      addOpts.Done,
			Created:   addOpts.Created,
			Completed: addOpts.Completed,
		})
		if err != nil {
			return err
		}
		if err := maybeSyncAfterWrite(ctx, svc); err != nil {
			return fmt.Errorf("sync push: %w", err)
		}

		writer := outputWriter()
		return writer.WriteTask(task)
	},
}

func init() {
	addCmd.Flags().StringVarP(&addOpts.Priority, "priority", "p", "", "priority letter")
	addCmd.Flags().StringSliceVarP(&addOpts.Projects, "project", "P", nil, "project tag (repeatable)")
	addCmd.Flags().StringSliceVarP(&addOpts.Contexts, "context", "c", nil, "context tag (repeatable)")
	addCmd.Flags().StringSliceVarP(&addOpts.Meta, "meta", "m", nil, "metadata key:value (repeatable)")
	addCmd.Flags().BoolVarP(&addOpts.Done, "done", "x", false, "mark task done")
	addCmd.Flags().StringVar(&addOpts.Created, "created", "", "creation date (YYYY-MM-DD)")
	addCmd.Flags().StringVar(&addOpts.Completed, "completed", "", "completion date (YYYY-MM-DD)")
}
