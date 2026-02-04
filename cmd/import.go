package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mholtzscher/ugh/internal/flags"
	"github.com/mholtzscher/ugh/internal/output"
	"github.com/mholtzscher/ugh/internal/service"
	"github.com/urfave/cli/v3"
)

var importCmd = &cli.Command{
	Name:      "import",
	Aliases:   []string{"in"},
	Usage:     "Import tasks from a backup file",
	Category:  "Backup",
	ArgsUsage: "<path|->",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  flags.FlagFormat,
			Usage: "input format (jsonl|json)",
			Value: "jsonl",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
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

		svc, err := newService(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = svc.Close() }()

		if err := maybeSyncBeforeWrite(ctx, svc); err != nil {
			return fmt.Errorf("sync pull: %w", err)
		}

		format := strings.ToLower(strings.TrimSpace(cmd.String(flags.FlagFormat)))
		if format == "" {
			format = "jsonl"
		}
		switch format {
		case "jsonl", "json":
			// ok
		default:
			return fmt.Errorf("invalid format %q (expected jsonl|json)", format)
		}

		var added int64
		var skipped int64
		if format == "json" {
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
		if err := maybeSyncAfterWrite(ctx, svc); err != nil {
			return fmt.Errorf("sync push: %w", err)
		}

		writer := outputWriter()
		return writer.WriteSummary(output.ImportSummary{Action: "import", Added: added, Skipped: skipped, File: path})
	},
}

type backupTask struct {
	Title      string            `json:"title"`
	Notes      string            `json:"notes,omitempty"`
	Status     string            `json:"status,omitempty"`
	Priority   string            `json:"priority,omitempty"`
	Done       bool              `json:"done,omitempty"`
	DueOn      string            `json:"dueOn,omitempty"`
	DeferUntil string            `json:"deferUntil,omitempty"`
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
	return service.CreateTaskRequest{
		Title:      t.Title,
		Notes:      t.Notes,
		Status:     t.Status,
		Priority:   t.Priority,
		Projects:   t.Projects,
		Contexts:   t.Contexts,
		Meta:       meta,
		DueOn:      t.DueOn,
		DeferUntil: t.DeferUntil,
		WaitingFor: t.WaitingFor,
		Done:       t.Done,
	}
}
