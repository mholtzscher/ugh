package output

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/mholtzscher/ugh/internal/store"
	"golang.org/x/term"
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
	_, err := fmt.Fprintf(out, "ID: %d\nStatus: %s\nPriority: %s\nCreated: %s\nCompleted: %s\nText: %s\n", task.ID, status, emptyDash(task.Priority), createdDateOrDash(task), dateOrDash(task.CompletionDate), todoLine(task))
	return err
}

func writeHumanList(out io.Writer, tasks []*store.Task, noColor bool) error {
	const (
		idWidth       = 4
		statusWidth   = 6
		priorityWidth = 8
		createdWidth  = 10
		minTaskWidth  = 30
		defaultTask   = 60
	)
	padding := 2 * 5
	taskWidth := defaultTask
	if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && width > 0 {
		reserved := idWidth + statusWidth + priorityWidth + createdWidth + padding
		if width > reserved+minTaskWidth {
			taskWidth = width - reserved
		}
	}
	columns := []table.Column{
		{Title: "ID", Width: idWidth},
		{Title: "Status", Width: statusWidth},
		{Title: "Priority", Width: priorityWidth},
		{Title: "Created", Width: createdWidth},
		{Title: "Task", Width: taskWidth},
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
			createdDateOrDash(task),
			todoLine(task),
		})
	}
	height := len(rows) + 1
	if height < 1 {
		height = 1
	}
	model := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(height),
		table.WithFocused(false),
	)

	if noColor {
		model.SetStyles(table.DefaultStyles())
	} else {
		styles := table.DefaultStyles()
		header := styles.Header.BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).BorderBottom(true).Bold(true)
		styles.Header = header
		styles.Selected = styles.Cell
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
