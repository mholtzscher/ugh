package tasks

import (
	"context"

	"github.com/mholtzscher/ugh/cmd/cmdutil"
	"github.com/mholtzscher/ugh/cmd/meta"
	"github.com/mholtzscher/ugh/internal/output"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/urfave/cli/v3"
)

func newDoneCmd(d Deps) *cli.Command {
	return &cli.Command{
		Name:      "done",
		Aliases:   []string{"d"},
		Usage:     "Mark tasks as done",
		Category:  meta.TasksCategory.String(),
		ArgsUsage: "<id...>",
		Action: cmdutil.WithServiceAndWriteSync(d.NewService, d.MaybeSyncBeforeWrite, d.MaybeSyncAfterWrite, func(ctx context.Context, cmd *cli.Command, svc service.Service) error {
			ids, err := parseIDs(commandArgs(cmd))
			if err != nil {
				return err
			}

			count, err := svc.SetDone(ctx, ids, true)
			if err != nil {
				return err
			}
			writer := d.OutputWriter()
			return writer.WriteSummary(output.Summary{Action: "done", Count: count, IDs: ids})
		}),
	}
}
