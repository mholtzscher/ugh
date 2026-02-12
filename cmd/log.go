package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/urfave/cli/v3"

	"github.com/mholtzscher/ugh/internal/flags"
	"github.com/mholtzscher/ugh/internal/output"
)

const defaultTaskLogLimit = 50

//nolint:gochecknoglobals // CLI command definitions are package-level by design.
var logCmd = &cli.Command{
	Name:      "log",
	Aliases:   []string{"events", "timeline"},
	Usage:     "Show per-task change history",
	Category:  "Tasks",
	ArgsUsage: "<id>",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    flags.FlagLimit,
			Aliases: []string{"n"},
			Usage:   "number of entries to show",
			Value:   defaultTaskLogLimit,
		},
		&cli.StringFlag{
			Name:  flags.FlagView,
			Usage: "render style (" + output.TaskEventViewsUsage + ")",
			Value: string(output.TaskEventViewDiff),
		},
		&cli.BoolFlag{
			Name:  flags.FlagVerbose,
			Usage: "show full values without truncation",
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

		task, err := svc.GetTask(ctx, ids[0])
		if err != nil {
			return fmt.Errorf("task %d not found", ids[0])
		}

		events, err := svc.ListTaskEvents(ctx, task.ID, int64(cmd.Int(flags.FlagLimit)))
		if err != nil {
			return err
		}

		view, err := output.ParseTaskEventView(cmd.String(flags.FlagView))
		if err != nil {
			return fmt.Errorf("parse %s: %w", flags.FlagView, err)
		}

		entries := make([]*output.TaskEventEntry, 0, len(events))
		for _, event := range events {
			entries = append(entries, &output.TaskEventEntry{
				ID:             event.ID,
				TaskID:         event.TaskID,
				Time:           time.Unix(event.Timestamp, 0).UTC(),
				Kind:           event.Kind,
				Summary:        event.Summary,
				ChangesJSON:    event.ChangesJSON,
				Origin:         event.Origin,
				ShellHistoryID: event.ShellHistoryID,
				ShellCommand:   event.ShellCommand,
			})
		}

		writer := outputWriter()
		return writer.WriteTaskEventsWithOptions(entries, output.TaskEventRenderOptions{
			View:    view,
			Verbose: cmd.Bool(flags.FlagVerbose),
		})
	},
}
