package output

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pterm/pterm"

	"github.com/mholtzscher/ugh/internal/store"
)

type Summary struct {
	Action string  `json:"action"`
	Count  int64   `json:"count"`
	IDs    []int64 `json:"ids,omitempty"`
	File   string  `json:"file,omitempty"`
}

func writeHumanTask(out io.Writer, noColor bool, task *store.Task) error {
	rows := pterm.TableData{
		{"Field", "Value"},
		{"ID", strconv.FormatInt(task.ID, 10)},
		{"State", string(task.State)},
		{"Prev State", stateOrDash(task.PrevState)},
		{"Created", dayFromTimeOrDash(task.CreatedAt)},
		{"Updated", dayFromTimeOrDash(task.UpdatedAt)},
		{"Completed", dateTimeOrDash(task.CompletedAt)},
		{"Due", dateOrDash(task.DueOn)},
		{"Waiting For", emptyDash(task.WaitingFor)},
		{"Title", task.Title},
		{"Notes", emptyDash(task.Notes)},
		{"Projects", joinListOrDash(task.Projects)},
		{"Contexts", joinListOrDash(task.Contexts)},
		{"Meta", metaOrDash(task.Meta)},
	}
	return renderTable(out, noColor, rows)
}

func writeHumanList(out io.Writer, noColor bool, tasks []*store.Task) error {
	rows := pterm.TableData{{"ID", "State", "Due", "Task"}}
	for _, task := range tasks {
		rows = append(rows, []string{
			strconv.FormatInt(task.ID, 10),
			string(task.State),
			dateOrDash(task.DueOn),
			task.Title,
		})
	}
	return renderTable(out, noColor, rows)
}

func writeHumanSummary(out io.Writer, noColor bool, summary any) error {
	switch value := summary.(type) {
	case Summary:
		ids := "-"
		if len(value.IDs) > 0 {
			ids = joinIDs(value.IDs)
		}
		rows := pterm.TableData{{"Action", "Count", "IDs"}, {value.Action, strconv.FormatInt(value.Count, 10), ids}}
		return renderTable(out, noColor, rows)
	default:
		return writeRenderedLine(out, noColor, pterm.DefaultBasicText.Sprintln(fmt.Sprintf("%v", value)))
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
		parts = append(parts, strconv.FormatInt(id, 10))
	}
	return strings.Join(parts, ",")
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

func stateOrDash(value *store.State) string {
	if value == nil || *value == "" {
		return "-"
	}
	return string(*value)
}

func writeHumanTags(out io.Writer, noColor bool, tags []store.NameCount) error {
	return writeHumanTagsWithCounts(out, noColor, tags)
}

func writeHumanTagsWithCounts(out io.Writer, noColor bool, tags []store.NameCount) error {
	rows := pterm.TableData{{"Name", "Count"}}
	for _, tag := range tags {
		rows = append(rows, []string{tag.Name, strconv.FormatInt(tag.Count, 10)})
	}
	return renderTable(out, noColor, rows)
}

func writeHumanKeyValues(out io.Writer, noColor bool, rows []KeyValue) error {
	tableData := pterm.TableData{{"Key", "Value"}}
	for _, row := range rows {
		tableData = append(tableData, []string{row.Key, row.Value})
	}
	return renderTable(out, noColor, tableData)
}

func renderTable(out io.Writer, noColor bool, data pterm.TableData) error {
	table := pterm.DefaultTable.WithHasHeader().WithLeftAlignment().WithBoxed().WithData(data)
	rendered, err := table.Srender()
	if err != nil {
		return err
	}
	return writeRenderedLine(out, noColor, rendered)
}

func writeRenderedLine(out io.Writer, noColor bool, line string) error {
	if noColor {
		line = pterm.RemoveColorFromString(line)
	}
	_, err := fmt.Fprint(out, line)
	return err
}

const maxCommandDisplayLength = 50

func writeHumanHistory(out io.Writer, noColor bool, entries []*HistoryEntry) error {
	if len(entries) == 0 {
		return writeRenderedLine(out, noColor, pterm.DefaultBasicText.Sprintln("No history entries"))
	}

	rows := pterm.TableData{{"Time", "Status", "Intent", "Command"}}
	for _, e := range entries {
		status := "✓"
		if !e.Success {
			status = "✗"
		}
		intent := e.Intent
		if intent == "" {
			intent = "-"
		}
		timeStr := e.Time.Format("2006-01-02 15:04")
		cmd := e.Command
		if len(cmd) > maxCommandDisplayLength {
			cmd = cmd[:maxCommandDisplayLength-3] + "..."
		}
		rows = append(rows, []string{timeStr, status, intent, cmd})
	}
	return renderTable(out, noColor, rows)
}
