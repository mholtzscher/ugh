package cmd

import (
	"context"
	"errors"

	"github.com/urfave/cli/v3"
)

var showCmd = &cli.Command{
	Name:      "show",
	Aliases:   []string{"s"},
	Usage:     "Show a task",
	ArgsUsage: "<id>",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Args().Len() != 1 {
			return errors.New("show requires a task id")
		}
		ids, err := parseIDs(commandArgs(cmd))
		if err != nil {
			return err
		}

		svc, err := newService(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = svc.Close() }()

		task, err := svc.GetTask(ctx, ids[0])
		if err != nil {
			return err
		}

		writer := outputWriter()
		return writer.WriteTask(task)
	},
}
