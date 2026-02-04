package cmd

import (
	"context"

	"github.com/mholtzscher/ugh/internal/flags"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/urfave/cli/v3"
)

var contextsCmd = &cli.Command{
	Name:    "contexts",
	Aliases: []string{"ctx"},
	Usage:   "List available context tags",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    flags.FlagAll,
			Aliases: []string{"a"},
			Usage:   "include completed tasks",
		},
		&cli.BoolFlag{
			Name:    flags.FlagDone,
			Aliases: []string{"x"},
			Usage:   "only completed tasks",
		},
		&cli.BoolFlag{
			Name:    flags.FlagTodo,
			Aliases: []string{"t"},
			Usage:   "only pending tasks",
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

		tags, err := svc.ListContexts(ctx, service.ListTagsRequest{
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
