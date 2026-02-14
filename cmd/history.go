package cmd

import (
	"context"
	"errors"
	"time"

	"github.com/urfave/cli/v3"

	"github.com/mholtzscher/ugh/internal/flags"
	"github.com/mholtzscher/ugh/internal/output"
	"github.com/mholtzscher/ugh/internal/store"
)

const defaultHistoryLimit = 50

//nolint:gochecknoglobals // CLI command definitions are package-level by design.
var historyCmd = &cli.Command{
	Name:     "history",
	Aliases:  []string{"hist", "h"},
	Usage:    "View shell command history",
	Category: "System",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    flags.FlagLimit,
			Aliases: []string{"n"},
			Usage:   "number of entries to show",
			Value:   defaultHistoryLimit,
		},
		&cli.StringFlag{
			Name:  flags.FlagSearch,
			Usage: "search command text",
		},
		&cli.StringFlag{
			Name:  flags.FlagIntent,
			Usage: "filter by intent",
		},
		&cli.BoolFlag{
			Name:  flags.FlagSuccess,
			Usage: "only successful commands",
		},
		&cli.BoolFlag{
			Name:  flags.FlagFailed,
			Usage: "only failed commands",
		},
		&cli.BoolFlag{
			Name:    flags.FlagClear,
			Aliases: []string{"c"},
			Usage:   "clear all history",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		svc, err := newService(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = svc.Close() }()

		if cmd.Bool(flags.FlagClear) {
			if clearErr := svc.ClearShellHistory(ctx); clearErr != nil {
				return clearErr
			}
			writer := outputWriter()
			return writer.WriteInfo("History cleared")
		}

		limit := int64(cmd.Int(flags.FlagLimit))
		search := cmd.String(flags.FlagSearch)
		intent := cmd.String(flags.FlagIntent)

		if cmd.Bool(flags.FlagSuccess) && cmd.Bool(flags.FlagFailed) {
			return errors.New("cannot use both --success and --failed flags")
		}

		var success *bool
		if cmd.Bool(flags.FlagSuccess) {
			s := true
			success = &s
		} else if cmd.Bool(flags.FlagFailed) {
			s := false
			success = &s
		}

		var histories []*store.ShellHistory
		if search != "" || intent != "" || success != nil {
			histories, err = svc.SearchShellHistory(ctx, search, intent, success, limit)
		} else {
			histories, err = svc.ListShellHistory(ctx, limit)
		}
		if err != nil {
			return err
		}

		entries := make([]*output.HistoryEntry, len(histories))
		for i, h := range histories {
			entries[i] = &output.HistoryEntry{
				ID:      h.ID,
				Time:    time.Unix(h.Timestamp, 0),
				Command: h.Command,
				Success: h.Success,
				Summary: h.ResultSummary,
				Intent:  h.Intent,
			}
		}

		writer := outputWriter()
		return writer.WriteHistory(entries)
	},
}
