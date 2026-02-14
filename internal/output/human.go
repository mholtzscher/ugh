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

func writeHumanTask(out io.Writer, formatter *TimeFormatter, task *store.Task) error {
	if task == nil {
		return nil
	}

	header := "Task " +
		pterm.ThemeDefault.PrimaryStyle.Sprint("#"+strconv.FormatInt(task.ID, 10)) +
		": " + task.Title

	rows := []KeyValue{
		{Key: "State", Value: formatDetailState(task.State)},
		{Key: "Prev State", Value: formatDetailPrevState(task.PrevState)},
		{Key: "Due", Value: formatDetailDate(formatter, task.DueOn, pterm.ThemeDefault.WarningMessageStyle)},
		{Key: "Waiting For", Value: emptyDash(task.WaitingFor)},
		{Key: "Projects", Value: formatDetailList(task.Projects, pterm.ThemeDefault.PrimaryStyle)},
		{Key: "Contexts", Value: formatDetailList(task.Contexts, pterm.ThemeDefault.SuccessMessageStyle)},
		{Key: "Meta", Value: metaOrDash(task.Meta)},
		{Key: "Created", Value: formatTimeOrDash(formatter, task.CreatedAt)},
		{Key: "Updated", Value: formatTimeOrDash(formatter, task.UpdatedAt)},
		{Key: "Completed", Value: formatTimePtrOrDash(formatter, task.CompletedAt)},
		{Key: "Notes", Value: emptyDash(task.Notes)},
	}

	var builder strings.Builder
	builder.WriteString(header)
	builder.WriteByte('\n')
	for _, row := range rows {
		builder.WriteString("  ")
		builder.WriteString(pterm.ThemeDefault.SecondaryStyle.Sprint(row.Key))
		builder.WriteString(": ")
		builder.WriteString(row.Value)
		builder.WriteByte('\n')
	}

	_, err := fmt.Fprint(out, builder.String())
	return err
}

func writeHumanList(out io.Writer, formatter *TimeFormatter, tasks []*store.Task) error {
	if len(tasks) == 0 {
		_, err := fmt.Fprintln(out, "No tasks found")
		return err
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Found %d task(s):\n", len(tasks)))
	for _, task := range tasks {
		builder.WriteString(formatTaskLine(formatter, task))
		builder.WriteByte('\n')
	}

	_, err := fmt.Fprint(out, builder.String())
	return err
}

func writeHumanSummary(out io.Writer, summary any) error {
	switch value := summary.(type) {
	case Summary:
		ids := ""
		if len(value.IDs) > 0 {
			ids = " ids=" + joinSummaryIDs(value.IDs)
		}
		line := fmt.Sprintf("%s: %d%s", value.Action, value.Count, ids)
		if value.File != "" {
			line += " (" + value.File + ")"
		}
		switch value.Action {
		case "done", "undo":
			line = pterm.ThemeDefault.SuccessMessageStyle.Sprint(line)
		case "rm":
			line = pterm.ThemeDefault.WarningMessageStyle.Sprint(line)
		}
		_, err := fmt.Fprintln(out, line)
		return err
	default:
		return writeRenderedLine(out, pterm.DefaultBasicText.Sprintln(fmt.Sprintf("%v", value)))
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
			line += " ids=" + joinSummaryIDs(value.IDs)
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

func formatTimeOrDash(formatter *TimeFormatter, value time.Time) string {
	if value.IsZero() {
		return "-"
	}
	return formatter.Format(value)
}

func formatTimePtrOrDash(formatter *TimeFormatter, value *time.Time) string {
	if value == nil {
		return "-"
	}
	return formatter.Format(*value)
}

func joinSummaryIDs(ids []int64) string {
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, "#"+strconv.FormatInt(id, 10))
	}
	return strings.Join(parts, ",")
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
	return writeHumanTagsWithCounts(out, tags)
}

func writeHumanTagsWithCounts(out io.Writer, tags []store.NameCount) error {
	rows := pterm.TableData{{"Name", "Count"}}
	for _, tag := range tags {
		rows = append(rows, []string{tag.Name, strconv.FormatInt(tag.Count, 10)})
	}
	return renderTable(out, rows)
}

func writeHumanKeyValues(out io.Writer, rows []KeyValue) error {
	tableData := pterm.TableData{{"Key", "Value"}}
	for _, row := range rows {
		tableData = append(tableData, []string{row.Key, row.Value})
	}
	return renderTable(out, tableData)
}

func writeHumanInfoBlock(out io.Writer, title string, rows []KeyValue) error {
	header := pterm.ThemeDefault.HighlightStyle.Sprint(title)
	if len(rows) == 0 {
		_, err := fmt.Fprintln(out, header)
		return err
	}
	data := make(pterm.TableData, 0, len(rows))
	for _, row := range rows {
		data = append(data, []string{row.Key, row.Value})
	}
	table := pterm.DefaultTable.WithData(data)
	rendered, err := table.Srender()
	if err != nil {
		return err
	}
	return writeRenderedLine(out, header+"\n"+rendered)
}

func writePlainInfoBlock(out io.Writer, title string, rows []KeyValue) error {
	var builder strings.Builder
	builder.WriteString(title)
	builder.WriteByte('\n')
	for _, row := range rows {
		builder.WriteString("  ")
		builder.WriteString(row.Key)
		builder.WriteString(": ")
		builder.WriteString(row.Value)
		builder.WriteByte('\n')
	}
	_, err := fmt.Fprint(out, builder.String())
	return err
}

func renderTable(out io.Writer, data pterm.TableData) error {
	table := pterm.DefaultTable.WithHasHeader().WithLeftAlignment().WithBoxed().WithData(data)
	rendered, err := table.Srender()
	if err != nil {
		return err
	}
	return writeRenderedLine(out, rendered)
}

func writeRenderedLine(out io.Writer, line string) error {
	_, err := fmt.Fprint(out, line)
	return err
}

func formatTaskLine(formatter *TimeFormatter, task *store.Task) string {
	if task == nil {
		return ""
	}

	state := task.State
	if state == "" {
		state = "inbox"
	}

	idStr := formatTaskID(task.ID)
	stateStr := formatTaskState(string(state))
	tags := formatTaskTags(task.Projects, task.Contexts)
	dueStr := formatTaskDueDate(formatter, task.DueOn)

	line := fmt.Sprintf("  %s %s %s", idStr, task.Title, stateStr)
	if tags != "" {
		line += " " + tags
	}
	if dueStr != "" {
		line += " " + dueStr
	}
	return line
}

func formatTaskID(id int64) string {
	return pterm.ThemeDefault.PrimaryStyle.Sprint("#" + strconv.FormatInt(id, 10))
}

func formatTaskState(state string) string {
	return pterm.ThemeDefault.SecondaryStyle.Sprint("[" + state + "]")
}

func formatTaskTags(projects, contexts []string) string {
	tags := make([]string, 0, len(projects)+len(contexts))
	for _, project := range projects {
		tags = append(tags, pterm.ThemeDefault.PrimaryStyle.Sprint("#"+project))
	}
	for _, context := range contexts {
		tags = append(tags, pterm.ThemeDefault.SuccessMessageStyle.Sprint("@"+context))
	}
	return strings.Join(tags, " ")
}

func formatTaskDueDate(formatter *TimeFormatter, dueOn *time.Time) string {
	if dueOn == nil {
		return ""
	}
	date := formatDateWithFormatter(formatter, dueOn)
	return pterm.ThemeDefault.WarningMessageStyle.Sprint(date)
}

func formatDetailState(value store.State) string {
	if value == "" {
		return "-"
	}
	return pterm.ThemeDefault.SecondaryStyle.Sprint(string(value))
}

func formatDetailPrevState(value *store.State) string {
	if value == nil || *value == "" {
		return "-"
	}
	return pterm.ThemeDefault.SecondaryStyle.Sprint(string(*value))
}

func formatDetailDate(formatter *TimeFormatter, value *time.Time, style pterm.Style) string {
	if value == nil {
		return "-"
	}
	date := formatDateWithFormatter(formatter, value)
	return style.Sprint(date)
}

func formatDetailList(values []string, style pterm.Style) string {
	if len(values) == 0 {
		return "-"
	}
	styled := make([]string, len(values))
	for i, value := range values {
		styled[i] = style.Sprint(value)
	}
	return strings.Join(styled, ", ")
}

const maxCommandDisplayLength = 50

func writeHumanHistory(out io.Writer, formatter *TimeFormatter, entries []*HistoryEntry) error {
	if len(entries) == 0 {
		return writeRenderedLine(out, pterm.DefaultBasicText.Sprintln("No history entries"))
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
		timeStr := formatter.Format(e.Time)
		cmd := e.Command
		if len(cmd) > maxCommandDisplayLength {
			cmd = cmd[:maxCommandDisplayLength-3] + "..."
		}
		rows = append(rows, []string{timeStr, status, intent, cmd})
	}
	return renderTable(out, rows)
}

//nolint:gocognit // Rendering diff output combines formatting and color decisions.
func writeHumanTaskVersionDiff(out io.Writer, formatter *TimeFormatter, versions []*store.TaskVersion) error {
	if len(versions) == 0 {
		return writeRenderedLine(out, pterm.DefaultBasicText.Sprintln("No task history entries"))
	}

	for i, current := range versions {
		header := fmt.Sprintf("Version %d  %s", current.VersionID, formatter.Format(current.UpdatedAt))
		header = pterm.ThemeDefault.PrimaryStyle.Sprint(header)
		if err := writeRenderedLine(out, header+"\n"); err != nil {
			return err
		}

		var prev *store.TaskVersion
		if i+1 < len(versions) {
			prev = versions[i+1]
		}
		changes := diffTaskVersion(prev, current)
		for _, change := range changes {
			prefix := "~"
			line := fmt.Sprintf("%s %s: %s -> %s", prefix, change.Field, emptyDash(change.Old), emptyDash(change.New))
			if change.Type == changeTypeAdd {
				prefix = "+"
				line = fmt.Sprintf("%s %s: %s", prefix, change.Field, emptyDash(change.New))
			}
			if change.Type == changeTypeRemove {
				prefix = "-"
				line = fmt.Sprintf("%s %s: %s", prefix, change.Field, emptyDash(change.Old))
			}

			var colored string
			switch prefix {
			case "+":
				colored = pterm.ThemeDefault.SuccessMessageStyle.Sprint(line)
			case "-":
				colored = pterm.ThemeDefault.ErrorMessageStyle.Sprint(line)
			default:
				colored = pterm.ThemeDefault.WarningMessageStyle.Sprint(line)
			}
			if err := writeRenderedLine(out, "  "+colored+"\n"); err != nil {
				return err
			}
		}

		if i < len(versions)-1 {
			if _, err := fmt.Fprintln(out); err != nil {
				return err
			}
		}
	}

	return nil
}
