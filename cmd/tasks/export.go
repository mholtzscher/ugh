package tasks

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/mholtzscher/ugh/internal/output"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/mholtzscher/ugh/internal/todotxt"

	"github.com/urfave/cli/v2"
)

func exportCommand() *cli.Command {
	return &cli.Command{
		Name:      "export",
		Aliases:   []string{"ex"},
		Usage:     "Export tasks to todo.txt",
		ArgsUsage: "<path|->",
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
			if c.Args().Len() != 1 {
				return errors.New("export requires a file path or -")
			}
			ctx, cancel := deps.WithTimeout(c.Context)
			defer cancel()
			path := c.Args().First()
			if deps.FlagBool(c, "json") && path == "-" {
				return errors.New("--json cannot be used when exporting to stdout")
			}

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

			var file *os.File
			if path == "-" {
				file = os.Stdout
			} else {
				created, err := os.Create(path)
				if err != nil {
					return fmt.Errorf("create file: %w", err)
				}
				defer func() { _ = created.Close() }()
				file = created
			}

			writer := bufio.NewWriter(file)
			for _, task := range tasks {
				line := todotxt.Format(todotxt.Parsed{
					Done:           task.Done,
					Priority:       task.Priority,
					CompletionDate: task.CompletionDate,
					CreationDate:   task.CreationDate,
					Description:    task.Description,
					Projects:       task.Projects,
					Contexts:       task.Contexts,
					Meta:           task.Meta,
					Unknown:        task.Unknown,
				})
				if _, err := fmt.Fprintln(writer, line); err != nil {
					return err
				}
			}
			if err := writer.Flush(); err != nil {
				return err
			}

			out := deps.OutputWriter(c)
			return out.WriteSummary(output.ExportSummary{Action: "export", Count: int64(len(tasks)), File: path})
		},
	}
}
