package output

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/mholtzscher/ugh/internal/store"
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
	table := tablewriter.NewWriter(out)
	table.Header("Field", "Value")
	rows := [][]string{
		{"ID", fmt.Sprintf("%d", task.ID)},
		{"Status", string(task.Status)},
		{"Done", fmt.Sprintf("%t", task.Done)},
		{"Priority", emptyDash(task.Priority)},
		{"Created", dayFromTimeOrDash(task.CreatedAt)},
		{"Updated", dayFromTimeOrDash(task.UpdatedAt)},
		{"Completed", dateTimeOrDash(task.CompletedAt)},
		{"Due", dateOrDash(task.DueOn)},
		{"Defer Until", dateOrDash(task.DeferUntil)},
		{"Waiting For", emptyDash(task.WaitingFor)},
		{"Title", task.Title},
		{"Notes", emptyDash(task.Notes)},
		{"Projects", joinListOrDash(task.Projects)},
		{"Contexts", joinListOrDash(task.Contexts)},
		{"Meta", metaOrDash(task.Meta)},
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
	table.Header("ID", "Status", "Due", "Defer", "Task")
	for _, task := range tasks {
		row := []string{
			fmt.Sprintf("%d", task.ID),
			string(task.Status),
			dateOrDash(task.DueOn),
			dateOrDash(task.DeferUntil),
			task.Title,
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

func dateTimeOrDash(value *time.Time) string {
	if value == nil {
		return "-"
	}
	return value.UTC().Format(time.RFC3339)
}

func dayFromTimeOrDash(value time.Time) string {
	if value.IsZero() {
		return "-"
	}
	return value.UTC().Format("2006-01-02")
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

func joinListOrDash(values []string) string {
	if len(values) == 0 {
		return "-"
	}
	return strings.Join(values, ", ")
}

func metaOrDash(meta map[string]string) string {
	if len(meta) == 0 {
		return "-"
	}
	keys := make([]string, 0, len(meta))
	for k := range meta {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+meta[k])
	}
	return strings.Join(parts, ", ")
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
