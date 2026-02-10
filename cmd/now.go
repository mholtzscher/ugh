package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/mholtzscher/ugh/internal/flags"
	"github.com/mholtzscher/ugh/internal/service"
)

//nolint:gochecknoglobals // CLI command definitions are package-level by design.
var nowCmd = &cli.Command{
	Name:     "now",
	Aliases:  []string{"n"},
	Usage:    "List tasks you can act on now",
	Category: "Lists",
	Action: func(ctx context.Context, _ *cli.Command) error {
		filterExpr, err := buildListFilterExpr(listFilterOptions{State: flags.TaskStateNow})
		if err != nil {
			return err
		}

		svc, err := newService(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = svc.Close() }()

		tasks, err := svc.ListTasks(ctx, service.ListTasksRequest{
			TodoOnly: true,
			Filter:   filterExpr,
		})
		if err != nil {
			return err
		}

		writer := outputWriter()
		return writer.WriteTasks(tasks)
	},
}
