package cmd

import (
	"context"

	"github.com/mholtzscher/ugh/internal/service"
	"github.com/urfave/cli/v3"
)

var nextCmd = &cli.Command{
	Name:     "next",
	Aliases:  []string{"n"},
	Usage:    "List next actions (available)",
	Category: "GTD Lists",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		svc, err := newService(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = svc.Close() }()

		today := todayUTC()
		tasks, err := svc.ListTasks(ctx, service.ListTasksRequest{
			TodoOnly:        true,
			Status:          "next",
			DeferOnOrBefore: today,
		})
		if err != nil {
			return err
		}

		writer := outputWriter()
		return writer.WriteTasks(tasks)
	},
}
