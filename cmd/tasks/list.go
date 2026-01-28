package tasks

import (
	"github.com/mholtzscher/ugh/internal/service"

	"github.com/urfave/cli/v2"
)

func listCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls", "l"},
		Usage:   "List tasks",
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
			&cli.StringFlag{
				Name:    "project",
				Aliases: []string{"P"},
				Usage:   "Filter by project",
			},
			&cli.StringFlag{
				Name:    "context",
				Aliases: []string{"c"},
				Usage:   "Filter by context",
			},
			&cli.StringFlag{
				Name:    "priority",
				Aliases: []string{"p"},
				Usage:   "Filter by priority",
			},
			&cli.StringFlag{
				Name:    "search",
				Aliases: []string{"s"},
				Usage:   "Search text",
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

			tasks, err := svc.ListTasks(ctx, service.ListTasksRequest{
				All:      c.Bool("all"),
				DoneOnly: c.Bool("done"),
				TodoOnly: c.Bool("todo"),
				Project:  c.String("project"),
				Context:  c.String("context"),
				Priority: c.String("priority"),
				Search:   c.String("search"),
			})
			if err != nil {
				return err
			}

			writer := deps.OutputWriter(c)
			return writer.WriteTasks(tasks)
		},
	}
}
