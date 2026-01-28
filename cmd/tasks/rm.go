package tasks

import (
	"fmt"

	"github.com/mholtzscher/ugh/internal/output"

	"github.com/urfave/cli/v2"
)

func rmCommand() *cli.Command {
	return &cli.Command{
		Name:      "rm",
		Usage:     "Delete tasks",
		ArgsUsage: "<id...>",
		Action: func(c *cli.Context) error {
			ctx, cancel := deps.WithTimeout(c.Context)
			defer cancel()
			ids, err := deps.ParseIDs(c.Args().Slice())
			if err != nil {
				return err
			}
			svc, err := deps.NewService(c)
			if err != nil {
				return err
			}
			defer func() { _ = svc.Close() }()

			if err := deps.MaybeSyncBeforeWrite(ctx, svc); err != nil {
				return fmt.Errorf("sync pull: %w", err)
			}

			count, err := svc.DeleteTasks(ctx, ids)
			if err != nil {
				return err
			}
			if err := deps.MaybeSyncAfterWrite(ctx, svc); err != nil {
				return fmt.Errorf("sync push: %w", err)
			}
			writer := deps.OutputWriter(c)
			return writer.WriteSummary(output.Summary{Action: "rm", Count: count, IDs: ids})
		},
	}
}
