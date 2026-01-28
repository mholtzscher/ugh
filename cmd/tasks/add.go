package tasks

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mholtzscher/ugh/internal/service"

	"github.com/urfave/cli/v2"
)

func addCommand() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Aliases:   []string{"a"},
		Usage:     "Add a task",
		ArgsUsage: "[todo.txt line]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "priority",
				Aliases: []string{"p"},
				Usage:   "Priority letter",
			},
			&cli.StringSliceFlag{
				Name:    "project",
				Aliases: []string{"P"},
				Usage:   "Project tag (repeatable)",
			},
			&cli.StringSliceFlag{
				Name:    "context",
				Aliases: []string{"c"},
				Usage:   "Context tag (repeatable)",
			},
			&cli.StringSliceFlag{
				Name:    "meta",
				Aliases: []string{"m"},
				Usage:   "Metadata key:value (repeatable)",
			},
			&cli.BoolFlag{
				Name:    "done",
				Aliases: []string{"x"},
				Usage:   "Mark task done",
			},
			&cli.StringFlag{
				Name:  "created",
				Usage: "Creation date (YYYY-MM-DD)",
			},
			&cli.StringFlag{
				Name:  "completed",
				Usage: "Completion date (YYYY-MM-DD)",
			},
		},
		Action: func(c *cli.Context) error {
			ctx, cancel := deps.WithTimeout(c.Context)
			defer cancel()

			line := strings.TrimSpace(strings.Join(c.Args().Slice(), " "))
			if line == "" {
				return errors.New("todo text required")
			}

			svc, err := deps.NewService(c)
			if err != nil {
				return err
			}
			defer func() { _ = svc.Close() }()

			if err := deps.MaybeSyncBeforeWrite(ctx, svc); err != nil {
				return fmt.Errorf("sync pull: %w", err)
			}

			task, err := svc.CreateTask(ctx, service.CreateTaskRequest{
				Line:      line,
				Priority:  c.String("priority"),
				Projects:  c.StringSlice("project"),
				Contexts:  c.StringSlice("context"),
				Meta:      c.StringSlice("meta"),
				Done:      c.Bool("done"),
				Created:   c.String("created"),
				Completed: c.String("completed"),
			})
			if err != nil {
				return err
			}
			if err := deps.MaybeSyncAfterWrite(ctx, svc); err != nil {
				return fmt.Errorf("sync push: %w", err)
			}

			writer := deps.OutputWriter(c)
			return writer.WriteTask(task)
		},
	}
}
