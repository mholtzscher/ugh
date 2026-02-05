package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mholtzscher/ugh/internal/flags"
	"github.com/mholtzscher/ugh/internal/output"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/mholtzscher/ugh/internal/store"
	"github.com/urfave/cli/v3"
)

var exportCmd = &cli.Command{
	Name:      "export",
	Aliases:   []string{"ex"},
	Usage:     "Export tasks to a backup file",
	Category:  "Backup",
	ArgsUsage: "<path|->",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  flags.FlagFormat,
			Usage: "output format (" + flags.FormatsUsage + ")",
			Value: flags.FormatJSONL,
			Action: flags.StringAction(
				flags.OneOfCaseInsensitiveRule(flags.FieldFormat, flags.Formats...),
			),
		},
		&cli.BoolFlag{
			Name:    flags.FlagAll,
			Aliases: []string{"a"},
			Usage:   "include completed tasks",
			Action: flags.BoolAction(
				flags.MutuallyExclusiveBoolFlagsRule(flags.FlagAll, flags.FlagDone, flags.FlagTodo),
			),
		},
		&cli.BoolFlag{
			Name:    flags.FlagDone,
			Aliases: []string{"x"},
			Usage:   "only completed tasks",
			Action: flags.BoolAction(
				flags.MutuallyExclusiveBoolFlagsRule(flags.FlagAll, flags.FlagDone, flags.FlagTodo),
			),
		},
		&cli.BoolFlag{
			Name:    flags.FlagTodo,
			Aliases: []string{"t"},
			Usage:   "only pending tasks",
			Action: flags.BoolAction(
				flags.MutuallyExclusiveBoolFlagsRule(flags.FlagAll, flags.FlagDone, flags.FlagTodo),
			),
		},
		&cli.StringFlag{
			Name:  flags.FlagState,
			Usage: "filter by state (" + flags.TaskStatesUsage + ")",
			Action: flags.StringAction(
				flags.OneOfCaseInsensitiveRule(flags.FieldState, flags.TaskStates...),
			),
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
		format := strings.ToLower(strings.TrimSpace(cmd.String(flags.FlagFormat)))
		if format == "" {
			format = flags.FormatJSONL
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
			State:    cmd.String(flags.FlagState),
			Project:  cmd.String(flags.FlagProject),
			Context:  cmd.String(flags.FlagContext),
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

		if format == flags.FormatJSON {
			enc := json.NewEncoder(file)
			payload := make([]any, 0, len(tasks))
			for _, task := range tasks {
				payload = append(payload, outputTask(task))
			}
			if err := enc.Encode(payload); err != nil {
				return err
			}
		} else {
			writer := bufio.NewWriter(file)
			enc := json.NewEncoder(writer)
			for _, task := range tasks {
				if err := enc.Encode(outputTask(task)); err != nil {
					return err
				}
			}
			if err := writer.Flush(); err != nil {
				return err
			}
		}

		// Avoid mixing data with a summary when writing to stdout.
		if path == "-" {
			return nil
		}
		out := outputWriter()
		return out.WriteSummary(output.ExportSummary{Action: "export", Count: int64(len(tasks)), File: path})
	},
}

func outputTask(task *store.Task) any {
	if task == nil {
		return nil
	}
	meta := task.Meta
	if meta == nil {
		meta = map[string]string{}
	}
	projects := task.Projects
	if projects == nil {
		projects = []string{}
	}
	contexts := task.Contexts
	if contexts == nil {
		contexts = []string{}
	}
	return map[string]any{
		"id":          task.ID,
		"state":       string(task.State),
		"title":       task.Title,
		"notes":       task.Notes,
		"dueOn":       formatDay(task.DueOn),
		"waitingFor":  task.WaitingFor,
		"projects":    projects,
		"contexts":    contexts,
		"meta":        meta,
		"createdAt":   task.CreatedAt.UTC().Format(time.RFC3339),
		"updatedAt":   task.UpdatedAt.UTC().Format(time.RFC3339),
		"completedAt": formatTime(task.CompletedAt),
	}
}

func formatDay(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.UTC().Format("2006-01-02")
}

func formatTime(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.UTC().Format(time.RFC3339)
}
