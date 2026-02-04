package cmd

import (
	"context"

	"github.com/mholtzscher/ugh/internal/service"
	"github.com/urfave/cli/v3"
)

var waitingCmd = &cli.Command{
	Name:     "waiting",
	Aliases:  []string{"w"},
	Usage:    "List waiting-for items",
	Category: "GTD Lists",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		svc, err := newService(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = svc.Close() }()

		tasks, err := svc.ListTasks(ctx, service.ListTasksRequest{
			TodoOnly: true,
			Status:   "waiting",
		})
		if err != nil {
			return err
		}

		writer := outputWriter()
		return writer.WriteTasks(tasks)
	},
}
