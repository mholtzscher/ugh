package cmd

import (
	"context"
	"errors"
	"strings"

	"github.com/mholtzscher/ugh/internal/service"
	"github.com/spf13/cobra"
)

var editOpts struct {
	Text string
}

var editCmd = &cobra.Command{
	Use:     "edit <id>",
	Aliases: []string{"e"},
	Short:   "Edit a task",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("edit requires a task id")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		ids, err := parseIDs(args)
		if err != nil {
			return err
		}

		line := strings.TrimSpace(editOpts.Text)
		if line == "" {
			return errors.New("todo text required")
		}

		svc, err := newTaskService(ctx)
		if err != nil {
			return err
		}
		defer svc.Close()

		saved, err := svc.UpdateTaskText(ctx, service.UpdateTaskRequest{
			ID:   ids[0],
			Text: line,
		})
		if err != nil {
			return err
		}
		writer := outputWriter()
		return writer.WriteTask(saved)
	},
}

func init() {
	editCmd.Flags().StringVarP(&editOpts.Text, "text", "t", "", "todo.txt formatted line")
}
