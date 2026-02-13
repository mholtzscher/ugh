package cmd

import (
	"context"
	"errors"

	"github.com/urfave/cli/v3"

	"github.com/mholtzscher/ugh/internal/flags"
)

const defaultTaskLogLimit = 20

//nolint:gochecknoglobals // CLI command definitions are package-level by design.
var logCmd = &cli.Command{
	Name:      "log",
	Usage:     "Show version history for a task",
	Category:  "Tasks",
	ArgsUsage: "<id>",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  flags.FlagLimit,
			Usage: "max versions to show",
			Value: defaultTaskLogLimit,
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Args().Len() != 1 {
			return errors.New("log requires a task id")
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

		versions, err := svc.ListTaskVersions(ctx, ids[0], int64(cmd.Int(flags.FlagLimit)))
		if err != nil {
			return err
		}

		writer := outputWriter()
		return writer.WriteTaskVersionDiff(versions)
	},
}
