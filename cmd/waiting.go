package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/mholtzscher/ugh/internal/flags"
	"github.com/mholtzscher/ugh/internal/service"
)

//nolint:gochecknoglobals // CLI command definitions are package-level by design.
var waitingCmd = &cli.Command{
	Name:     "waiting",
	Aliases:  []string{"w"},
	Usage:    "List waiting-for items",
	Category: "Lists",
	Action: func(ctx context.Context, _ *cli.Command) error {
		svc, err := newService(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = svc.Close() }()

		tasks, err := svc.ListTasks(ctx, service.ListTasksRequest{
			TodoOnly: true,
			State:    flags.TaskStateWaiting,
		})
		if err != nil {
			return err
		}

		writer := outputWriter()
		return writer.WriteTasks(tasks)
	},
}
