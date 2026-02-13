package cmd

import (
	"context"
	"errors"

	"github.com/urfave/cli/v3"

	"github.com/mholtzscher/ugh/internal/flags"
	"github.com/mholtzscher/ugh/internal/service"
)

const listLimitUnset = -1

//nolint:gochecknoglobals // CLI command definitions are package-level by design.
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
				flags.OneOfCaseInsensitiveRule(flags.FieldState, flags.TaskStates()...),
			),
		},
		&cli.StringFlag{
			Name:    flags.FlagProject,
			Aliases: []string{"p"},
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
		&cli.StringFlag{
			Name:  flags.FlagWhere,
			Usage: "filter expression (e.g. \"state:now or project:work\")",
		},
		&cli.BoolFlag{
			Name:  flags.FlagRecent,
			Usage: "show most recently changed tasks",
		},
		&cli.IntFlag{
			Name:  flags.FlagLimit,
			Usage: "max tasks to show",
			Value: listLimitUnset,
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		limit := cmd.Int(flags.FlagLimit)
		if limit == 0 || limit < listLimitUnset {
			return errors.New("limit must be greater than 0")
		}
		if limit == listLimitUnset {
			limit = 0
		}

		filterExpr, err := buildListFilterExpr(listFilterOptions{
			Where:   cmd.String(flags.FlagWhere),
			State:   cmd.String(flags.FlagState),
			Project: cmd.String(flags.FlagProject),
			Context: cmd.String(flags.FlagContext),
			Search:  cmd.String(flags.FlagSearch),
		})
		if err != nil {
			return err
		}

		svc, err := newService(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = svc.Close() }()

		req := service.ListTasksRequest{
			All:      cmd.Bool(flags.FlagAll),
			DoneOnly: cmd.Bool(flags.FlagDone),
			TodoOnly: cmd.Bool(flags.FlagTodo),
			Filter:   filterExpr,
			Recent:   cmd.Bool(flags.FlagRecent),
			Limit:    int64(limit),
		}
		tasks, err := svc.ListTasks(ctx, req)
		if err != nil {
			return err
		}

		writer := outputWriter()
		return writer.WriteTasks(tasks)
	},
}
