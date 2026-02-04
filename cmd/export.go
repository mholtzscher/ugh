package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/mholtzscher/ugh/internal/flags"
	"github.com/mholtzscher/ugh/internal/output"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/mholtzscher/ugh/internal/todotxt"
	"github.com/urfave/cli/v3"
)

var exportCmd = &cli.Command{
	Name:      "export",
	Aliases:   []string{"ex"},
	Usage:     "Export tasks to todo.txt",
	ArgsUsage: "<path|->",
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
		&cli.StringFlag{
			Name:    flags.FlagProject,
			Aliases: []string{"P"},
			Usage:   "filter by project",
		},
		&cli.StringFlag{
			Name:    flags.FlagContext,
			Aliases: []string{"c"},
			Usage:   "filter by context",
		},
		&cli.StringFlag{
			Name:    flags.FlagPriority,
			Aliases: []string{"p"},
			Usage:   "filter by priority",
		},
		&cli.StringFlag{
			Name:    flags.FlagSearch,
			Aliases: []string{"s"},
			Usage:   "search text",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Args().Len() != 1 {
			return errors.New("export requires a file path or -")
		}
		path := cmd.Args().Get(0)
		if cmd.Bool(flags.FlagJSON) && path == "-" {
			return errors.New("--json cannot be used when exporting to stdout")
		}

		svc, err := newService(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = svc.Close() }()

		tasks, err := svc.ListTasks(ctx, service.ListTasksRequest{
			All:      cmd.Bool(flags.FlagAll),
			DoneOnly: cmd.Bool(flags.FlagDone),
			TodoOnly: cmd.Bool(flags.FlagTodo),
			Project:  cmd.String(flags.FlagProject),
			Context:  cmd.String(flags.FlagContext),
			Priority: cmd.String(flags.FlagPriority),
			Search:   cmd.String(flags.FlagSearch),
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

		out := outputWriter()
		return out.WriteSummary(output.ExportSummary{Action: "export", Count: int64(len(tasks)), File: path})
	},
}
