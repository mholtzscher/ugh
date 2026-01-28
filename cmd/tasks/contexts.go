package tasks

import (
	"github.com/mholtzscher/ugh/internal/service"

	"github.com/urfave/cli/v2"
)

func contextsCommand() *cli.Command {
	return &cli.Command{
		Name:    "contexts",
		Aliases: []string{"ctx"},
		Usage:   "List available context tags",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "Include completed tasks",
			},
			&cli.BoolFlag{
				Name:    "done",
				Aliases: []string{"x"},
				Usage:   "Only completed tasks",
			},
			&cli.BoolFlag{
				Name:    "todo",
				Aliases: []string{"t"},
				Usage:   "Only pending tasks",
			},
			&cli.BoolFlag{
				Name:  "counts",
				Usage: "Include counts",
			},
		},
		Action: func(c *cli.Context) error {
			ctx, cancel := deps.WithTimeout(c.Context)
			defer cancel()
			svc, err := deps.NewService(c)
			if err != nil {
				return err
			}
			defer func() { _ = svc.Close() }()

			tags, err := svc.ListContexts(ctx, service.ListTagsRequest{
				All:      c.Bool("all"),
				DoneOnly: c.Bool("done"),
				TodoOnly: c.Bool("todo"),
			})
			if err != nil {
				return err
			}

			writer := deps.OutputWriter(c)
			if c.Bool("counts") {
				return writer.WriteTagsWithCounts(tags)
			}
			return writer.WriteTags(tags)
		},
	}
}
