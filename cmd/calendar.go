package cmd

import (
	"context"

	"github.com/mholtzscher/ugh/internal/service"
	"github.com/urfave/cli/v3"
)

var calendarCmd = &cli.Command{
	Name:     "calendar",
	Aliases:  []string{"cal"},
	Usage:    "List items with due dates",
	Category: "Lists",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		svc, err := newService(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = svc.Close() }()

		tasks, err := svc.ListTasks(ctx, service.ListTasksRequest{
			TodoOnly: true,
			DueOnly:  true,
		})
		if err != nil {
			return err
		}

		writer := outputWriter()
		return writer.WriteTasks(tasks)
	},
}
