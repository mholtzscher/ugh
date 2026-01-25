package output

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/mholtzscher/ugh/internal/store"
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
	_, err := fmt.Fprintf(out, "ID: %d\nStatus: %s\nPriority: %s\nCreated: %s\nCompleted: %s\nText: %s\n", task.ID, status, emptyDash(task.Priority), dateOrDash(task.CreationDate), dateOrDash(task.CompletionDate), todoLine(task))
	return err
}

func writeHumanList(out io.Writer, tasks []*store.Task, noColor bool) error {
	columns := []table.Column{
		{Title: "ID", Width: 6},
		{Title: "St", Width: 3},
		{Title: "Pri", Width: 4},
		{Title: "Created", Width: 10},
		{Title: "Task", Width: 60},
	}
	rows := make([]table.Row, 0, len(tasks))
	for _, task := range tasks {
		status := "todo"
		if task.Done {
			status = "done"
		}
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", task.ID),
			status,
			emptyDash(task.Priority),
			dateOrDash(task.CreationDate),
			todoLine(task),
		})
	}
	model := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
	)

	if noColor {
		model.SetStyles(table.DefaultStyles())
	} else {
		styles := table.DefaultStyles()
		header := styles.Header.BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).BorderBottom(true).Bold(true)
		styles.Header = header
		styles.Selected = styles.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57"))
		model.SetStyles(styles)
	}
	_, err := fmt.Fprintln(out, model.View())
	return err
}

func writeHumanSummary(out io.Writer, summary any) error {
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

func joinIDs(ids []int64) string {
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, fmt.Sprintf("%d", id))
	}
	return strings.Join(parts, ",")
}
