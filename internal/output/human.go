package output

import (
	"encoding/json"
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

const maxTaskEventSummaryDisplayLength = 60

const maxTaskEventChangesDisplayLength = 100

const maxTaskEventValueDisplayLength = 24

const (
	taskEventFieldProjects      = "projects"
	taskEventFieldContexts      = "contexts"
	taskEventFieldMeta          = "meta"
	taskEventDiffLineMultiplier = 2
)

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

func writeHumanTaskEvents(
	out io.Writer,
	noColor bool,
	entries []*TaskEventEntry,
	options TaskEventRenderOptions,
) error {
	view := options.View
	if view == "" {
		view = TaskEventViewTimeline
	}

	switch view {
	case TaskEventViewTimeline:
		return writeHumanTaskEventsTimeline(out, noColor, entries, options.Verbose)
	case TaskEventViewTable:
		return writeHumanTaskEventsTable(out, noColor, entries, options.Verbose)
	case TaskEventViewCompact:
		return writeHumanTaskEventsCompact(out, noColor, entries, options.Verbose)
	case TaskEventViewDiff:
		return writeHumanTaskEventsDiff(out, noColor, entries, options.Verbose)
	default:
		return writeHumanTaskEventsTable(out, noColor, entries, options.Verbose)
	}
}

func writeHumanTaskEventsTable(out io.Writer, noColor bool, entries []*TaskEventEntry, verbose bool) error {
	if len(entries) == 0 {
		return writeRenderedLine(out, noColor, pterm.DefaultBasicText.Sprintln("No task history entries"))
	}

	rows := pterm.TableData{{"Time", "Kind", "Origin", "Command", "Summary", "Changes"}}
	for _, e := range entries {
		origin := e.Origin
		if origin == "" {
			origin = "-"
		}
		command := e.ShellCommand
		if command == "" {
			command = "-"
		}
		if !verbose && len(command) > maxCommandDisplayLength {
			command = command[:maxCommandDisplayLength-3] + "..."
		}
		summary := e.Summary
		if summary == "" {
			summary = "-"
		}
		if !verbose && len(summary) > maxTaskEventSummaryDisplayLength {
			summary = summary[:maxTaskEventSummaryDisplayLength-3] + "..."
		}
		changes := summarizeTaskEventChanges(e.ChangesJSON, verbose)
		if !verbose && len(changes) > maxTaskEventChangesDisplayLength {
			changes = changes[:maxTaskEventChangesDisplayLength-3] + "..."
		}
		rows = append(rows, []string{
			e.Time.Format("2006-01-02 15:04"),
			colorTaskEventKind(e.Kind, e.Kind),
			origin,
			command,
			summary,
			changes,
		})
	}

	return renderTable(out, noColor, rows)
}

func writeHumanTaskEventsCompact(out io.Writer, noColor bool, entries []*TaskEventEntry, verbose bool) error {
	if len(entries) == 0 {
		return writeBasicLine(out, noColor, "No task history entries")
	}

	for _, e := range entries {
		summary := emptyDash(e.Summary)
		changes := summarizeTaskEventChanges(e.ChangesJSON, verbose)
		command := e.ShellCommand
		if !verbose && len(command) > maxCommandDisplayLength {
			command = command[:maxCommandDisplayLength-3] + "..."
		}

		line := fmt.Sprintf(
			"%s %s %s %s",
			e.Time.Format("2006-01-02 15:04"),
			colorTaskEventKind(e.Kind, taskEventKindIcon(e.Kind)),
			e.Kind,
			summary,
		)
		if changes != "-" {
			line += " | " + changes
		}
		if command != "" {
			line += " (" + command + ")"
		}

		if err := writeBasicLine(out, noColor, line); err != nil {
			return err
		}
	}

	return nil
}

//nolint:gocognit // Timeline formatting intentionally keeps rendering logic together.
func writeHumanTaskEventsTimeline(out io.Writer, noColor bool, entries []*TaskEventEntry, verbose bool) error {
	if len(entries) == 0 {
		return writeBasicLine(out, noColor, "No task history entries")
	}

	for i, e := range entries {
		header := fmt.Sprintf(
			"%s %s %s",
			e.Time.Format("2006-01-02 15:04"),
			colorTaskEventKind(e.Kind, taskEventKindIcon(e.Kind)),
			emptyDash(e.Summary),
		)
		if e.Origin != "" {
			header += " [" + e.Origin + "]"
		}
		if err := writeBasicLine(out, noColor, header); err != nil {
			return err
		}

		command := e.ShellCommand
		if command != "" {
			if !verbose && len(command) > maxTaskEventChangesDisplayLength {
				command = command[:maxTaskEventChangesDisplayLength-3] + "..."
			}
			if err := writeBasicLine(out, noColor, "  "+pterm.LightBlue("cmd:")+" "+command); err != nil {
				return err
			}
		}

		details := taskEventDetailLines(e.ChangesJSON, verbose)
		for _, detail := range details {
			if err := writeBasicLine(out, noColor, "  "+detail); err != nil {
				return err
			}
		}

		if i < len(entries)-1 {
			if err := writeBasicLine(out, noColor, "-"); err != nil {
				return err
			}
		}
	}

	return nil
}

//nolint:gocognit // Diff formatting intentionally keeps rendering logic together.
func writeHumanTaskEventsDiff(out io.Writer, noColor bool, entries []*TaskEventEntry, verbose bool) error {
	if len(entries) == 0 {
		return writeBasicLine(out, noColor, "No task history entries")
	}

	for i, e := range entries {
		header := fmt.Sprintf(
			"%s %s %s",
			e.Time.Format("2006-01-02 15:04"),
			colorTaskEventKind(e.Kind, taskEventKindIcon(e.Kind)),
			emptyDash(e.Summary),
		)
		if err := writeBasicLine(out, noColor, header); err != nil {
			return err
		}

		if e.ShellCommand != "" {
			command := e.ShellCommand
			if !verbose && len(command) > maxTaskEventChangesDisplayLength {
				command = command[:maxTaskEventChangesDisplayLength-3] + "..."
			}
			if err := writeBasicLine(out, noColor, "  "+pterm.LightBlue(">")+" "+command); err != nil {
				return err
			}
		}

		diffLines := taskEventDiffLines(e.ChangesJSON, verbose)
		for _, line := range diffLines {
			if err := writeBasicLine(out, noColor, "  "+colorTaskEventDiffLine(line)); err != nil {
				return err
			}
		}

		if i < len(entries)-1 {
			if err := writeBasicLine(out, noColor, "-"); err != nil {
				return err
			}
		}
	}

	return nil
}

func writeBasicLine(out io.Writer, noColor bool, line string) error {
	rendered := pterm.DefaultBasicText.Sprint(line) + "\n"
	return writeRenderedLine(out, noColor, rendered)
}

func summarizeTaskEventChanges(changesJSON string, verbose bool) string {
	if strings.TrimSpace(changesJSON) == "" {
		return "-"
	}

	var changes map[string]any
	if err := json.Unmarshal([]byte(changesJSON), &changes); err != nil {
		return "changed"
	}

	if _, ok := changes["after"]; ok && len(changes) == 1 {
		return "created"
	}
	if _, ok := changes["before"]; ok && len(changes) == 1 {
		return "deleted"
	}

	keys := make([]string, 0, len(changes))
	for key := range changes {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		part := summarizeTaskEventFieldChange(key, changes[key], verbose)
		if part != "" {
			parts = append(parts, part)
		}
	}

	if len(parts) == 0 {
		return "changed"
	}

	return strings.Join(parts, "; ")
}

func summarizeTaskEventFieldChange(key string, value any, verbose bool) string {
	change, ok := value.(map[string]any)
	if !ok {
		return key
	}

	switch key {
	case taskEventFieldProjects, taskEventFieldContexts:
		return summarizeTaskEventListChange(key, change)
	case taskEventFieldMeta:
		return summarizeTaskEventMetaChange(change)
	default:
		from, hasFrom := change["from"]
		to, hasTo := change["to"]
		if hasFrom && hasTo {
			return fmt.Sprintf(
				"%s: %s -> %s",
				key,
				formatTaskEventValue(from, taskEventValueMaxLen(verbose)),
				formatTaskEventValue(to, taskEventValueMaxLen(verbose)),
			)
		}
	}

	return key
}

func summarizeTaskEventListChange(key string, change map[string]any) string {
	added := anySliceToStrings(change["added"])
	removed := anySliceToStrings(change["removed"])

	parts := make([]string, 0, len(added)+len(removed))
	for _, value := range added {
		parts = append(parts, "+"+value)
	}
	for _, value := range removed {
		parts = append(parts, "-"+value)
	}

	if len(parts) == 0 {
		return key
	}

	return key + ":" + strings.Join(parts, ",")
}

func summarizeTaskEventMetaChange(change map[string]any) string {
	parts := make([]string, 0)

	for _, key := range mapKeysFromAny(change["added"]) {
		parts = append(parts, "+"+key)
	}
	for _, key := range mapKeysFromAny(change["updated"]) {
		parts = append(parts, "~"+key)
	}
	for _, key := range anySliceToStrings(change["removed"]) {
		parts = append(parts, "-"+key)
	}

	if len(parts) == 0 {
		return taskEventFieldMeta
	}

	return "meta:" + strings.Join(parts, ",")
}

func anySliceToStrings(value any) []string {
	slice, ok := value.([]any)
	if !ok {
		return nil
	}
	result := make([]string, 0, len(slice))
	for _, item := range slice {
		result = append(result, fmt.Sprintf("%v", item))
	}
	sort.Strings(result)
	return result
}

func mapKeysFromAny(value any) []string {
	m, ok := value.(map[string]any)
	if !ok {
		return nil
	}
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func formatTaskEventValue(value any, maxLen int) string {
	switch typed := value.(type) {
	case string:
		if typed == "" {
			return `""`
		}
		if maxLen > 0 && len(typed) > maxLen {
			return strconv.Quote(typed[:maxLen-3] + "...")
		}
		return strconv.Quote(typed)
	default:
		return fmt.Sprintf("%v", value)
	}
}

func taskEventValueMaxLen(verbose bool) int {
	if verbose {
		return 0
	}
	return maxTaskEventValueDisplayLength
}

func taskEventKindIcon(kind string) string {
	switch kind {
	case store.TaskEventKindCreate:
		return "+"
	case store.TaskEventKindDone:
		return "✓"
	case store.TaskEventKindUndo:
		return "↺"
	case store.TaskEventKindDelete:
		return "✗"
	default:
		return "~"
	}
}

func colorTaskEventKind(kind string, value string) string {
	switch kind {
	case store.TaskEventKindCreate:
		return pterm.LightGreen(value)
	case store.TaskEventKindDone:
		return pterm.LightCyan(value)
	case store.TaskEventKindUndo:
		return pterm.LightYellow(value)
	case store.TaskEventKindDelete:
		return pterm.LightRed(value)
	default:
		return pterm.LightMagenta(value)
	}
}

func colorTaskEventDiffLine(line string) string {
	if strings.HasPrefix(line, "+ ") {
		return pterm.LightGreen(line)
	}
	if strings.HasPrefix(line, "- ") {
		return pterm.LightRed(line)
	}
	if strings.HasPrefix(line, "~ ") {
		return pterm.LightYellow(line)
	}
	return line
}

func taskEventDetailLines(changesJSON string, verbose bool) []string {
	changes, ok := parseTaskEventChanges(changesJSON)
	if !ok {
		return []string{"changes: -"}
	}

	if after, exists := changes["after"]; exists && len(changes) == 1 {
		if !verbose {
			return []string{"created"}
		}
		return snapshotLines("+", after, verbose)
	}
	if before, exists := changes["before"]; exists && len(changes) == 1 {
		if !verbose {
			return []string{"deleted"}
		}
		return snapshotLines("-", before, verbose)
	}

	keys := make([]string, 0, len(changes))
	for key := range changes {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	maxLen := taskEventValueMaxLen(verbose)
	lines := make([]string, 0, len(keys))
	for _, key := range keys {
		change, isMap := changes[key].(map[string]any)
		if !isMap {
			lines = append(lines, key)
			continue
		}

		switch key {
		case taskEventFieldProjects, taskEventFieldContexts:
			lines = append(lines, summarizeTaskEventListChange(key, change))
		case taskEventFieldMeta:
			lines = append(lines, summarizeTaskEventMetaChange(change))
		default:
			from, hasFrom := change["from"]
			to, hasTo := change["to"]
			if hasFrom && hasTo {
				lines = append(lines, fmt.Sprintf(
					"%s: %s -> %s",
					key,
					formatTaskEventValue(from, maxLen),
					formatTaskEventValue(to, maxLen),
				))
			}
		}
	}

	if len(lines) == 0 {
		return []string{"changed"}
	}

	return lines
}

//nolint:gocognit // Diff extraction handles multiple structured change shapes.
func taskEventDiffLines(changesJSON string, verbose bool) []string {
	changes, ok := parseTaskEventChanges(changesJSON)
	if !ok {
		return []string{"~ changed"}
	}

	maxLen := taskEventValueMaxLen(verbose)

	if after, exists := changes["after"]; exists && len(changes) == 1 {
		return snapshotLines("+", after, verbose)
	}
	if before, exists := changes["before"]; exists && len(changes) == 1 {
		return snapshotLines("-", before, verbose)
	}

	keys := make([]string, 0, len(changes))
	for key := range changes {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	lines := make([]string, 0, len(keys)*taskEventDiffLineMultiplier)
	for _, key := range keys {
		change, isMap := changes[key].(map[string]any)
		if !isMap {
			continue
		}

		switch key {
		case taskEventFieldProjects, taskEventFieldContexts:
			for _, value := range anySliceToStrings(change["removed"]) {
				lines = append(lines, "- "+key+":"+value)
			}
			for _, value := range anySliceToStrings(change["added"]) {
				lines = append(lines, "+ "+key+":"+value)
			}
		case taskEventFieldMeta:
			for _, metaKey := range anySliceToStrings(change["removed"]) {
				lines = append(lines, "- meta."+metaKey)
			}
			for _, metaKey := range mapKeysFromAny(change["added"]) {
				lines = append(lines, "+ meta."+metaKey)
			}
			if updatedMap, updatedOK := change["updated"].(map[string]any); updatedOK {
				updatedKeys := mapKeysFromAny(updatedMap)
				for _, metaKey := range updatedKeys {
					if transition, transitionOK := updatedMap[metaKey].(map[string]any); transitionOK {
						lines = append(lines, fmt.Sprintf(
							"- meta.%s: %s",
							metaKey,
							formatTaskEventValue(transition["from"], maxLen),
						))
						lines = append(lines, fmt.Sprintf(
							"+ meta.%s: %s",
							metaKey,
							formatTaskEventValue(transition["to"], maxLen),
						))
					}
				}
			}
		default:
			from, hasFrom := change["from"]
			to, hasTo := change["to"]
			if hasFrom && hasTo {
				lines = append(lines, fmt.Sprintf("- %s: %s", key, formatTaskEventValue(from, maxLen)))
				lines = append(lines, fmt.Sprintf("+ %s: %s", key, formatTaskEventValue(to, maxLen)))
			}
		}
	}

	if len(lines) == 0 {
		return []string{"~ changed"}
	}

	return lines
}

func parseTaskEventChanges(changesJSON string) (map[string]any, bool) {
	if strings.TrimSpace(changesJSON) == "" {
		return nil, false
	}
	var changes map[string]any
	if err := json.Unmarshal([]byte(changesJSON), &changes); err != nil {
		return nil, false
	}
	return changes, true
}

func snapshotLines(prefix string, snapshotAny any, verbose bool) []string {
	snapshot, ok := snapshotAny.(map[string]any)
	if !ok {
		return []string{prefix + " changed"}
	}

	maxLen := taskEventValueMaxLen(verbose)
	keys := make([]string, 0, len(snapshot))
	for key := range snapshot {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	lines := make([]string, 0, len(keys))
	for _, key := range keys {
		lines = append(lines, fmt.Sprintf("%s %s: %s", prefix, key, formatTaskEventValue(snapshot[key], maxLen)))
	}
	return lines
}
