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
	Category:  "Tasks",
	ArgsUsage: "<title>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  flags.FlagState,
			Usage: "task state (inbox|now|waiting|later|done)",
			Value: "inbox",
		},
		&cli.StringFlag{
			Name:  flags.FlagNotes,
			Usage: "notes",
		},
		&cli.StringSliceFlag{
			Name:    flags.FlagProject,
			Aliases: []string{"P"},
			Usage:   "project (repeatable)",
		},
		&cli.StringSliceFlag{
			Name:    flags.FlagContext,
			Aliases: []string{"c"},
			Usage:   "context (repeatable)",
		},
		&cli.StringSliceFlag{
			Name:    flags.FlagMeta,
			Aliases: []string{"m"},
			Usage:   "metadata key:value (repeatable)",
		},
		&cli.StringFlag{
			Name:  flags.FlagDueOn,
			Usage: "due date (YYYY-MM-DD)",
		},
		&cli.StringFlag{
			Name:  flags.FlagWaitingFor,
			Usage: "waiting for (person/thing)",
		},
		&cli.BoolFlag{
			Name:    flags.FlagDone,
			Aliases: []string{"x"},
			Usage:   "mark task done",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		title := strings.TrimSpace(strings.Join(commandArgs(cmd), " "))
		if title == "" {
			return errors.New("title required")
		}

		svc, err := newService(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = svc.Close() }()

		if err := maybeSyncBeforeWrite(ctx, svc); err != nil {
			return fmt.Errorf("sync pull: %w", err)
		}

		state := cmd.String(flags.FlagState)
		if cmd.Bool(flags.FlagDone) {
			state = "done"
		}

		task, err := svc.CreateTask(ctx, service.CreateTaskRequest{
			Title:      title,
			Notes:      cmd.String(flags.FlagNotes),
			State:      state,
			Projects:   cmd.StringSlice(flags.FlagProject),
			Contexts:   cmd.StringSlice(flags.FlagContext),
			Meta:       cmd.StringSlice(flags.FlagMeta),
			DueOn:      cmd.String(flags.FlagDueOn),
			WaitingFor: cmd.String(flags.FlagWaitingFor),
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
