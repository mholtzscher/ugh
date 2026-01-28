package output

import (
	"fmt"
	"io"
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

type SyncSummary struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

type SyncStatusSummary struct {
	Action          string `json:"action"`
	LastPullTime    int64  `json:"lastPullTime"`
	LastPushTime    int64  `json:"lastPushTime"`
	PendingChanges  int64  `json:"pendingChanges"`
	NetworkSent     int64  `json:"networkSentBytes"`
	NetworkReceived int64  `json:"networkReceivedBytes"`
	Revision        string `json:"revision"`
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
		{"Projects", formatProjects(task.Projects)},
		{"Contexts", formatContexts(task.Contexts)},
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
	table.Header("ID", "Status", "Priority", "Created", "Description", "Projects", "Contexts")
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
			task.Description,
			formatProjects(task.Projects),
			formatContexts(task.Contexts),
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
	case SyncSummary:
		table := tablewriter.NewWriter(out)
		table.Header("Action", "Message")
		if err := table.Append(value.Action, value.Message); err != nil {
			return err
		}
		return table.Render()
	case SyncStatusSummary:
		table := tablewriter.NewWriter(out)
		table.Header("Field", "Value")
		if err := table.Append("last_pull", formatSyncTime(value.LastPullTime)); err != nil {
			return err
		}
		if err := table.Append("last_push", formatSyncTime(value.LastPushTime)); err != nil {
			return err
		}
		if err := table.Append("pending_changes", fmt.Sprintf("%d", value.PendingChanges)); err != nil {
			return err
		}
		if err := table.Append("network_sent", fmt.Sprintf("%d bytes", value.NetworkSent)); err != nil {
			return err
		}
		if err := table.Append("network_received", fmt.Sprintf("%d bytes", value.NetworkReceived)); err != nil {
			return err
		}
		if err := table.Append("revision", value.Revision); err != nil {
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
	case SyncSummary:
		_, err := fmt.Fprintf(out, "%s: %s\n", value.Action, value.Message)
		return err
	case SyncStatusSummary:
		if _, err := fmt.Fprintf(out, "last_pull:\t%s\n", formatSyncTime(value.LastPullTime)); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(out, "last_push:\t%s\n", formatSyncTime(value.LastPushTime)); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(out, "pending_changes:\t%d\n", value.PendingChanges); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(out, "network_sent:\t%d bytes\n", value.NetworkSent); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(out, "network_received:\t%d bytes\n", value.NetworkReceived); err != nil {
			return err
		}
		_, err := fmt.Fprintf(out, "revision:\t%s\n", value.Revision)
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

func formatSyncTime(value int64) string {
	if value == 0 {
		return "never"
	}
	return fmt.Sprintf("%d", value)
}

func formatProjects(projects []string) string {
	if len(projects) == 0 {
		return "-"
	}
	parts := make([]string, 0, len(projects))
	for _, p := range projects {
		parts = append(parts, "+"+p)
	}
	return strings.Join(parts, " ")
}

func formatContexts(contexts []string) string {
	if len(contexts) == 0 {
		return "-"
	}
	parts := make([]string, 0, len(contexts))
	for _, c := range contexts {
		parts = append(parts, "@"+c)
	}
	return strings.Join(parts, " ")
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
