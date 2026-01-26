package output

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/mholtzscher/ugh/internal/store"
	"github.com/mholtzscher/ugh/internal/todotxt"
	"github.com/olekukonko/tablewriter"
)

type Summary struct {
	Action string  `json:"action"`
	Count  int64   `json:"count"`
	IDs    []int64 `json:"ids,omitempty"`
	File   string  `json:"file,omitempty"`
}

type ImportSummary struct {
	Action  string `json:"action"`
	Added   int64  `json:"added"`
	Skipped int64  `json:"skipped"`
	File    string `json:"file,omitempty"`
}

type ExportSummary struct {
	Action string `json:"action"`
	Count  int64  `json:"count"`
	File   string `json:"file,omitempty"`
}

func writeHumanTask(out io.Writer, task *store.Task) error {
	status := "todo"
	if task.Done {
		status = "done"
	}
	table := tablewriter.NewWriter(out)
	table.Header("Field", "Value")
	rows := [][]string{
		{"ID", fmt.Sprintf("%d", task.ID)},
		{"Status", status},
		{"Priority", emptyDash(task.Priority)},
		{"Created", createdDateOrDash(task)},
		{"Completed", dateOrDash(task.CompletionDate)},
		{"Description", task.Description},
	}
	for _, row := range rows {
		if err := appendRow(table, row); err != nil {
			return err
		}
	}
	return table.Render()
}

func writeHumanList(out io.Writer, tasks []*store.Task) error {
	table := tablewriter.NewWriter(out)
	table.Header("ID", "Status", "Priority", "Created", "Task")
	for _, task := range tasks {
		status := "todo"
		if task.Done {
			status = "done"
		}
		row := []string{
			fmt.Sprintf("%d", task.ID),
			status,
			emptyDash(task.Priority),
			createdDateOrDash(task),
			humanTaskText(task),
		}
		if err := appendRow(table, row); err != nil {
			return err
		}
	}
	return table.Render()
}

func writeHumanSummary(out io.Writer, summary any) error {
	switch value := summary.(type) {
	case Summary:
		table := tablewriter.NewWriter(out)
		table.Header("Action", "Count", "IDs")
		ids := "-"
		if len(value.IDs) > 0 {
			ids = joinIDs(value.IDs)
		}
		if err := table.Append(value.Action, fmt.Sprintf("%d", value.Count), ids); err != nil {
			return err
		}
		return table.Render()
	case ImportSummary:
		table := tablewriter.NewWriter(out)
		table.Header("Action", "Added", "Skipped", "File")
		file := emptyDash(value.File)
		if err := table.Append(value.Action, fmt.Sprintf("%d", value.Added), fmt.Sprintf("%d", value.Skipped), file); err != nil {
			return err
		}
		return table.Render()
	case ExportSummary:
		table := tablewriter.NewWriter(out)
		table.Header("Action", "Count", "File")
		file := emptyDash(value.File)
		if err := table.Append(value.Action, fmt.Sprintf("%d", value.Count), file); err != nil {
			return err
		}
		return table.Render()
	default:
		_, err := fmt.Fprintf(out, "%v\n", value)
		return err
	}
}

func writePlainSummary(out io.Writer, summary any) error {
	switch value := summary.(type) {
	case Summary:
		line := fmt.Sprintf("%s: %d", value.Action, value.Count)
		if value.File != "" {
			line += " (" + value.File + ")"
		}
		if len(value.IDs) > 0 {
			line += " ids=" + joinIDs(value.IDs)
		}
		_, err := fmt.Fprintln(out, line)
		return err
	case ImportSummary:
		line := fmt.Sprintf("%s: added %d, skipped %d", value.Action, value.Added, value.Skipped)
		if value.File != "" {
			line += " (" + value.File + ")"
		}
		_, err := fmt.Fprintln(out, line)
		return err
	case ExportSummary:
		line := fmt.Sprintf("%s: %d", value.Action, value.Count)
		if value.File != "" {
			line += " (" + value.File + ")"
		}
		_, err := fmt.Fprintln(out, line)
		return err
	default:
		_, err := fmt.Fprintf(out, "%v\n", value)
		return err
	}
}

func emptyDash(value string) string {
	if value == "" {
		return "-"
	}
	return value
}

func dateOrDash(value *time.Time) string {
	if value == nil {
		return "-"
	}
	return value.Format("2006-01-02")
}

func createdDateOrDash(task *store.Task) string {
	if task == nil {
		return "-"
	}
	if task.CreationDate != nil {
		return task.CreationDate.Format("2006-01-02")
	}
	if !task.CreatedAt.IsZero() {
		return task.CreatedAt.UTC().Format("2006-01-02")
	}
	return "-"
}

func joinIDs(ids []int64) string {
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, fmt.Sprintf("%d", id))
	}
	return strings.Join(parts, ",")
}

func appendRow(table *tablewriter.Table, row []string) error {
	values := make([]any, len(row))
	for i, val := range row {
		values[i] = val
	}
	return table.Append(values...)
}

func humanTaskText(task *store.Task) string {
	if task == nil {
		return ""
	}
	parsed := todotxt.Parsed{
		Description: task.Description,
		Projects:    task.Projects,
		Contexts:    task.Contexts,
		Meta:        task.Meta,
		Unknown:     task.Unknown,
	}
	return todotxt.Format(parsed)
}

func writeHumanTags(out io.Writer, tags []store.NameCount) error {
	table := tablewriter.NewWriter(out)
	table.Header("Name", "Count")
	for _, tag := range tags {
		if err := appendRow(table, []string{tag.Name, fmt.Sprintf("%d", tag.Count)}); err != nil {
			return err
		}
	}
	return table.Render()
}
