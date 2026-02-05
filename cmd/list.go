package cmd

import (
	"context"

	"github.com/mholtzscher/ugh/internal/flags"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/urfave/cli/v3"
)

var listCmd = &cli.Command{
	Name:     "list",
	Aliases:  []string{"ls", "l"},
	Usage:    "List tasks",
	Category: "Tasks",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    flags.FlagAll,
			Aliases: []string{"a"},
			Usage:   "include completed tasks",
			Action: flags.BoolAction(
				flags.MutuallyExclusiveBoolFlagsRule(flags.FlagAll, flags.FlagDone, flags.FlagTodo),
			),
		},
		&cli.BoolFlag{
			Name:    flags.FlagDone,
			Aliases: []string{"x"},
			Usage:   "only completed tasks",
			Action: flags.BoolAction(
				flags.MutuallyExclusiveBoolFlagsRule(flags.FlagAll, flags.FlagDone, flags.FlagTodo),
			),
		},
		&cli.BoolFlag{
			Name:    flags.FlagTodo,
			Aliases: []string{"t"},
			Usage:   "only pending tasks",
			Action: flags.BoolAction(
				flags.MutuallyExclusiveBoolFlagsRule(flags.FlagAll, flags.FlagDone, flags.FlagTodo),
			),
		},
		&cli.StringFlag{
			Name:  flags.FlagState,
			Usage: "filter by state (" + flags.TaskStatesUsage + ")",
			Action: flags.StringAction(
				flags.OneOfCaseInsensitiveRule(flags.FieldState, flags.TaskStates...),
			),
		},
		&cli.StringFlag{
			Name:    flags.FlagProject,
			Aliases: []string{"P"},
			Usage:   "filter by project",
		},
		&cli.StringFlag{
			Name:    flags.FlagContext,
			Aliases: []string{"c"},
			Usage:   "filter by context",
		},
		&cli.StringFlag{
			Name:    flags.FlagSearch,
			Aliases: []string{"s"},
			Usage:   "search text",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		svc, err := newService(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = svc.Close() }()

		tasks, err := svc.ListTasks(ctx, service.ListTasksRequest{
			All:      cmd.Bool(flags.FlagAll),
			DoneOnly: cmd.Bool(flags.FlagDone),
			TodoOnly: cmd.Bool(flags.FlagTodo),
			State:    cmd.String(flags.FlagState),
			Project:  cmd.String(flags.FlagProject),
			Context:  cmd.String(flags.FlagContext),
			Search:   cmd.String(flags.FlagSearch),
		})
		if err != nil {
			return err
		}

		writer := outputWriter()
		return writer.WriteTasks(tasks)
	},
}
