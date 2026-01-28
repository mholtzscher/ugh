package tasks

import (
	"errors"

	"github.com/urfave/cli/v2"
)

func showCommand() *cli.Command {
	return &cli.Command{
		Name:      "show",
		Aliases:   []string{"s"},
		Usage:     "Show a task",
		ArgsUsage: "<id>",
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				return errors.New("show requires a task id")
			}

			ctx, cancel := deps.WithTimeout(c.Context)
			defer cancel()
			ids, err := deps.ParseIDs(c.Args().Slice())
			if err != nil {
				return err
			}

			svc, err := deps.NewService(c)
			if err != nil {
				return err
			}
			defer func() { _ = svc.Close() }()

			task, err := svc.GetTask(ctx, ids[0])
			if err != nil {
				return err
			}

			writer := deps.OutputWriter(c)
			return writer.WriteTask(task)
		},
	}
}
