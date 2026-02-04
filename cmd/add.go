package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mholtzscher/ugh/internal/flags"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/urfave/cli/v3"
)

var addCmd = &cli.Command{
	Name:      "add",
	Aliases:   []string{"a"},
	Usage:     "Add a task",
	ArgsUsage: "[todo.txt line]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    flags.FlagPriority,
			Aliases: []string{"p"},
			Usage:   "priority letter",
		},
		&cli.StringSliceFlag{
			Name:    flags.FlagProject,
			Aliases: []string{"P"},
			Usage:   "project tag (repeatable)",
		},
		&cli.StringSliceFlag{
			Name:    flags.FlagContext,
			Aliases: []string{"c"},
			Usage:   "context tag (repeatable)",
		},
		&cli.StringSliceFlag{
			Name:    flags.FlagMeta,
			Aliases: []string{"m"},
			Usage:   "metadata key:value (repeatable)",
		},
		&cli.BoolFlag{
			Name:    flags.FlagDone,
			Aliases: []string{"x"},
			Usage:   "mark task done",
		},
		&cli.StringFlag{
			Name:  flags.FlagCreated,
			Usage: "creation date (YYYY-MM-DD)",
		},
		&cli.StringFlag{
			Name:  flags.FlagCompleted,
			Usage: "completion date (YYYY-MM-DD)",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		line := strings.TrimSpace(strings.Join(commandArgs(cmd), " "))
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
			Priority:  cmd.String(flags.FlagPriority),
			Projects:  cmd.StringSlice(flags.FlagProject),
			Contexts:  cmd.StringSlice(flags.FlagContext),
			Meta:      cmd.StringSlice(flags.FlagMeta),
			Done:      cmd.Bool(flags.FlagDone),
			Created:   cmd.String(flags.FlagCreated),
			Completed: cmd.String(flags.FlagCompleted),
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
