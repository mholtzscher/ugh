//nolint:dupl // Context and project listing commands intentionally share structure.
package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/mholtzscher/ugh/internal/flags"
	"github.com/mholtzscher/ugh/internal/service"
)

//nolint:gochecknoglobals // CLI command definitions are package-level by design.
var projectsCmd = &cli.Command{
	Name:     "projects",
	Aliases:  []string{"proj"},
	Usage:    "List projects",
	Category: "Projects & Contexts",
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
		&cli.BoolFlag{
			Name:  flags.FlagCounts,
			Usage: "include counts",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		svc, err := newService(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = svc.Close() }()

		tags, err := svc.ListProjects(ctx, service.ListTagsRequest{
			All:      cmd.Bool(flags.FlagAll),
			DoneOnly: cmd.Bool(flags.FlagDone),
			TodoOnly: cmd.Bool(flags.FlagTodo),
		})
		if err != nil {
			return err
		}

		writer := outputWriter()
		if cmd.Bool(flags.FlagCounts) {
			return writer.WriteTagsWithCounts(tags)
		}
		return writer.WriteTags(tags)
	},
}
