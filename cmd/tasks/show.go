package tasks

import (
	"context"
	"errors"

	"github.com/mholtzscher/ugh/cmd/cmdutil"
	"github.com/mholtzscher/ugh/cmd/meta"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/urfave/cli/v3"
)

func newShowCmd(d Deps) *cli.Command {
	return &cli.Command{
		Name:      "show",
		Aliases:   []string{"s"},
		Usage:     "Show a task",
		Category:  meta.TasksCategory.String(),
		ArgsUsage: "<id>",
		Action: cmdutil.WithService(d.NewService, func(ctx context.Context, cmd *cli.Command, svc service.Service) error {
			if cmd.Args().Len() != 1 {
				return errors.New("show requires a task id")
			}
			ids, err := parseIDs(commandArgs(cmd))
			if err != nil {
				return err
			}

			task, err := svc.GetTask(ctx, ids[0])
			if err != nil {
				return err
			}

			writer := d.OutputWriter()
			return writer.WriteTask(task)
		}),
	}
}
