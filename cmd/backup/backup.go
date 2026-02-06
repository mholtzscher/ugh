package backup

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/mholtzscher/ugh/cmd/cmdutil"
	"github.com/mholtzscher/ugh/cmd/meta"
	"github.com/mholtzscher/ugh/cmd/registry"
	"github.com/mholtzscher/ugh/internal/flags"
	"github.com/mholtzscher/ugh/internal/output"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/mholtzscher/ugh/internal/store"

	"github.com/urfave/cli/v3"
)

// Deps holds dependencies injected from the parent cmd package.
type Deps struct {
	NewService           func(context.Context) (service.Service, error)
	MaybeSyncBeforeWrite func(context.Context, service.Service) error
	MaybeSyncAfterWrite  func(context.Context, service.Service) error
	OutputWriter         func() output.Writer
}

const (
	importID registry.ID = "backup.import"
	exportID registry.ID = "backup.export"
)

// Register adds backup command specs to the registry.
func Register(r *registry.Registry, d Deps) error {
	return r.AddAll(
		registry.Spec{ID: importID, Source: "cmd/backup", Build: func() *cli.Command { return newImportCmd(d) }},
		registry.Spec{ID: exportID, Source: "cmd/backup", Build: func() *cli.Command { return newExportCmd(d) }},
	)
}

func newImportCmd(d Deps) *cli.Command {
	return &cli.Command{
		Name:      "import",
		Aliases:   []string{"in"},
		Usage:     "Import tasks from a backup file",
		Category:  meta.BackupCategory.String(),
		ArgsUsage: "<path|->",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  flags.FlagFormat,
				Usage: "input format (" + flags.FormatsUsage + ")",
				Value: flags.FormatJSONL,
				Action: flags.StringAction(
					flags.OneOfCaseInsensitiveRule(flags.FieldFormat, flags.Formats...),
				),
			},
		},
		Action: cmdutil.WithServiceAndWriteSync(d.NewService, d.MaybeSyncBeforeWrite, d.MaybeSyncAfterWrite, func(ctx context.Context, cmd *cli.Command, svc service.Service) error {
			if cmd.Args().Len() != 1 {
				return errors.New("import requires a file path or -")
			}
			path := cmd.Args().Get(0)
			var reader io.Reader
			if path == "-" {
				reader = os.Stdin
			} else {
				file, err := os.Open(path)
				if err != nil {
					return fmt.Errorf("open file: %w", err)
				}
				defer func() { _ = file.Close() }()
				reader = file
			}

			format := strings.ToLower(strings.TrimSpace(cmd.String(flags.FlagFormat)))
			if format == "" {
				format = flags.FormatJSONL
			}

			var added int64
			var skipped int64
			if format == flags.FormatJSON {
				data, err := io.ReadAll(reader)
				if err != nil {
					return fmt.Errorf("read import: %w", err)
				}
				var items []backupTask
				if err := json.Unmarshal(data, &items); err != nil {
					return fmt.Errorf("parse json: %w", err)
				}
				for _, item := range items {
					if strings.TrimSpace(item.Title) == "" {
						skipped++
						continue
					}
					if _, err := svc.CreateTask(ctx, item.toCreateRequest()); err != nil {
						return err
					}
					added++
				}
			} else {
				scanner := bufio.NewScanner(reader)
				buf := make([]byte, 0, 64*1024)
				scanner.Buffer(buf, 1024*1024)
				for scanner.Scan() {
					line := strings.TrimSpace(scanner.Text())
					if line == "" {
						skipped++
						continue
					}
					var item backupTask
					if err := json.Unmarshal([]byte(line), &item); err != nil {
						return fmt.Errorf("parse jsonl: %w", err)
					}
					if strings.TrimSpace(item.Title) == "" {
						skipped++
						continue
					}
					if _, err := svc.CreateTask(ctx, item.toCreateRequest()); err != nil {
						return err
					}
					added++
				}
				if err := scanner.Err(); err != nil {
					return fmt.Errorf("read import: %w", err)
				}
			}

			writer := d.OutputWriter()
			return writer.WriteSummary(output.ImportSummary{Action: "import", Added: added, Skipped: skipped, File: path})
		}),
	}
}

func newExportCmd(d Deps) *cli.Command {
	return &cli.Command{
		Name:      "export",
		Aliases:   []string{"ex"},
		Usage:     "Export tasks to a backup file",
		Category:  meta.BackupCategory.String(),
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
				Aliases: []string{"p"},
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
		Action: cmdutil.WithService(d.NewService, func(ctx context.Context, cmd *cli.Command, svc service.Service) error {
			if cmd.Args().Len() != 1 {
				return errors.New("export requires a file path or -")
			}
			path := cmd.Args().Get(0)
			format := strings.ToLower(strings.TrimSpace(cmd.String(flags.FlagFormat)))
			if format == "" {
				format = flags.FormatJSONL
			}

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
			writer := d.OutputWriter()
			return writer.WriteSummary(output.ExportSummary{Action: "export", Count: int64(len(tasks)), File: path})
		}),
	}
}

type backupTask struct {
	Title      string            `json:"title"`
	Notes      string            `json:"notes,omitempty"`
	State      string            `json:"state,omitempty"`
	DueOn      string            `json:"dueOn,omitempty"`
	WaitingFor string            `json:"waitingFor,omitempty"`
	Projects   []string          `json:"projects,omitempty"`
	Contexts   []string          `json:"contexts,omitempty"`
	Meta       map[string]string `json:"meta,omitempty"`
}

func (t backupTask) toCreateRequest() service.CreateTaskRequest {
	meta := make([]string, 0)
	for k, v := range t.Meta {
		meta = append(meta, k+":"+v)
	}
	state := strings.TrimSpace(t.State)
	return service.CreateTaskRequest{
		Title:      t.Title,
		Notes:      t.Notes,
		State:      state,
		Projects:   t.Projects,
		Contexts:   t.Contexts,
		Meta:       meta,
		DueOn:      t.DueOn,
		WaitingFor: t.WaitingFor,
	}
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
