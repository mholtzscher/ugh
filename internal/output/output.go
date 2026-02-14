package output

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pterm/pterm"
	"golang.org/x/term"

	"github.com/mholtzscher/ugh/internal/config"
	"github.com/mholtzscher/ugh/internal/nlp"
	"github.com/mholtzscher/ugh/internal/store"
)

type Writer struct {
	Out       io.Writer
	JSON      bool
	TTY       bool
	formatter *TimeFormatter
}

type KeyValue struct {
	Key   string
	Value string
}

type ContextStatus struct {
	SelectedID *int64  `json:"selectedId,omitempty"`
	LastIDs    []int64 `json:"lastIds,omitempty"`
	Project    string  `json:"project,omitempty"`
	Context    string  `json:"context,omitempty"`
}

type ViewHelp struct {
	Entries []ViewHelpEntry `json:"entries"`
	Usage   string          `json:"usage,omitempty"`
}

type ViewHelpEntry struct {
	Label       string `json:"label"`
	Description string `json:"description"`
}

func NewWriter(jsonMode bool, displayCfg config.Display) Writer {
	return Writer{
		Out:       os.Stdout,
		JSON:      jsonMode,
		TTY:       term.IsTerminal(int(os.Stdout.Fd())),
		formatter: NewTimeFormatter(displayCfg),
	}
}

func (w Writer) WriteTask(task *store.Task) error {
	if task == nil {
		return errors.New("task is nil")
	}
	if w.JSON {
		payload := toTaskJSON(task)
		return writeJSON(w.Out, payload)
	}

	if w.isHumanMode() {
		return w.writeHumanTask(task)
	}
	_, err := fmt.Fprintln(w.Out, w.plainLine(task))
	return err
}

func (w Writer) WriteTasks(tasks []*store.Task) error {
	if w.JSON {
		payload := make([]TaskJSON, 0, len(tasks))
		for _, task := range tasks {
			payload = append(payload, toTaskJSON(task))
		}
		return writeJSON(w.Out, payload)
	}

	if w.isHumanMode() {
		return w.writeHumanList(tasks)
	}
	for _, task := range tasks {
		if _, err := fmt.Fprintln(w.Out, w.plainLine(task)); err != nil {
			return err
		}
	}
	return nil
}

func (w Writer) WriteTags(tags []store.NameCount) error {
	if w.JSON {
		return writeJSON(w.Out, tags)
	}

	if w.isHumanMode() {
		return writeHumanTags(w.Out, tags)
	}
	for _, tag := range tags {
		if _, err := fmt.Fprintln(w.Out, tag.Name); err != nil {
			return err
		}
	}
	return nil
}

func (w Writer) WriteTagsWithCounts(tags []store.NameCount) error {
	if w.JSON {
		return writeJSON(w.Out, tags)
	}
	if w.isHumanMode() {
		return writeHumanTagsWithCounts(w.Out, tags)
	}

	for _, tag := range tags {
		if _, err := fmt.Fprintf(w.Out, "%s\t%d\n", tag.Name, tag.Count); err != nil {
			return err
		}
	}
	return nil
}

func (w Writer) WriteSummary(summary any) error {
	if w.JSON {
		return writeJSON(w.Out, summary)
	}
	if w.isHumanMode() {
		return writeHumanSummary(w.Out, summary)
	}
	return writePlainSummary(w.Out, summary)
}

func (w Writer) WriteKeyValues(rows []KeyValue) error {
	if w.isHumanMode() {
		return writeHumanKeyValues(w.Out, rows)
	}
	for _, row := range rows {
		if _, err := fmt.Fprintf(w.Out, "%s:\t%s\n", row.Key, row.Value); err != nil {
			return err
		}
	}
	return nil
}

func (w Writer) WriteInfoBlock(title string, rows []KeyValue) error {
	if w.isHumanMode() {
		return writeHumanInfoBlock(w.Out, title, rows)
	}
	return writePlainInfoBlock(w.Out, title, rows)
}

func (w Writer) WriteContextStatus(status ContextStatus) error {
	if w.JSON {
		return writeJSON(w.Out, status)
	}
	rows := []KeyValue{
		{Key: formatContextLabel("Selected", w.TTY), Value: formatContextID(status.SelectedID, w.TTY)},
		{Key: formatContextLabel("Last", w.TTY), Value: formatContextIDs(status.LastIDs, w.TTY)},
		{Key: formatContextLabel("Project", w.TTY), Value: formatContextTag(status.Project, "#", w.TTY)},
		{Key: formatContextLabel("Context", w.TTY), Value: formatContextTag(status.Context, "@", w.TTY)},
	}
	return w.WriteInfoBlock("Current Context:", rows)
}

func (w Writer) WriteViewHelp(help ViewHelp) error {
	if w.JSON {
		return writeJSON(w.Out, help)
	}
	rows := make([]KeyValue, 0, len(help.Entries)+1)
	for _, entry := range help.Entries {
		rows = append(rows, KeyValue{
			Key:   formatViewLabel(entry.Label, w.TTY),
			Value: formatViewDescription(entry.Description, w.TTY),
		})
	}
	if help.Usage != "" {
		rows = append(rows, KeyValue{
			Key:   formatContextLabel("Usage", w.TTY),
			Value: formatViewUsage(help.Usage, w.TTY),
		})
	}
	return w.WriteInfoBlock("Available Views:", rows)
}

func (w Writer) WriteLine(line string) error {
	if !w.isHumanMode() {
		_, err := fmt.Fprintln(w.Out, line)
		return err
	}
	formatted := pterm.DefaultBasicText.Sprintln(line)
	return writeRenderedLine(w.Out, formatted)
}

func (w Writer) WriteInfo(line string) error {
	return w.writePrefixLine(pterm.Info, line)
}

func (w Writer) WriteSuccess(line string) error {
	return w.writePrefixLine(pterm.Success, line)
}

func (w Writer) WriteWarning(line string) error {
	return w.writePrefixLine(pterm.Warning, line)
}

func (w Writer) WriteError(line string) error {
	return w.writePrefixLine(pterm.Error, line)
}

func (w Writer) WriteErr(err error) error {
	if err == nil {
		return nil
	}

	var diagnosticErr nlp.DiagnosticError
	if errors.As(err, &diagnosticErr) {
		diagnostics := diagnosticErr.Diagnostics()
		if len(diagnostics) > 0 {
			line := err.Error()
			if diagnostics[0].Hint != "" {
				line += " (hint: " + diagnostics[0].Hint + ")"
			}
			return w.WriteError(line)
		}
	}

	return w.WriteError(err.Error())
}

func (w Writer) writePrefixLine(printer pterm.PrefixPrinter, line string) error {
	if !w.isHumanMode() {
		prefix := ""
		if printer == pterm.Error {
			prefix = "Error: "
		}
		_, err := fmt.Fprintln(w.Out, prefix+line)
		return err
	}
	formatted := printer.Sprintln(line)
	return writeRenderedLine(w.Out, formatted)
}

type TaskJSON struct {
	ID          int64             `json:"id"`
	State       string            `json:"state"`
	Title       string            `json:"title"`
	Notes       string            `json:"notes,omitempty"`
	DueOn       string            `json:"dueOn,omitempty"`
	WaitingFor  string            `json:"waitingFor,omitempty"`
	CompletedAt string            `json:"completedAt,omitempty"`
	Projects    []string          `json:"projects"`
	Contexts    []string          `json:"contexts"`
	Meta        map[string]string `json:"meta"`
	CreatedAt   string            `json:"createdAt"`
	UpdatedAt   string            `json:"updatedAt"`
}

func toTaskJSON(task *store.Task) TaskJSON {
	projects := task.Projects
	if projects == nil {
		projects = []string{}
	}
	contexts := task.Contexts
	if contexts == nil {
		contexts = []string{}
	}
	meta := task.Meta
	if meta == nil {
		meta = map[string]string{}
	}
	return TaskJSON{
		ID:          task.ID,
		State:       string(task.State),
		Title:       task.Title,
		Notes:       task.Notes,
		DueOn:       formatDate(task.DueOn),
		WaitingFor:  task.WaitingFor,
		CompletedAt: formatDateTimePtr(task.CompletedAt),
		Projects:    projects,
		Contexts:    contexts,
		Meta:        meta,
		CreatedAt:   formatDateTime(task.CreatedAt),
		UpdatedAt:   formatDateTime(task.UpdatedAt),
	}
}

func (w Writer) plainLine(task *store.Task) string {
	if task == nil {
		return ""
	}
	due := w.formatDateWithFormatter(task.DueOn)
	fields := []string{
		strconv.FormatInt(task.ID, 10),
		string(task.State),
		due,
		task.WaitingFor,
		task.Title,
	}
	return strings.Join(fields, "\t")
}

const contextNoneValue = "none"

func formatContextLabel(label string, tty bool) string {
	if !tty {
		return label
	}
	return pterm.ThemeDefault.SecondaryStyle.Sprint(label)
}

func formatContextID(id *int64, tty bool) string {
	if id == nil {
		return contextNoneValue
	}
	value := fmt.Sprintf("#%d", *id)
	if !tty {
		return value
	}
	return pterm.ThemeDefault.PrimaryStyle.Sprint(value)
}

func formatContextIDs(ids []int64, tty bool) string {
	if len(ids) == 0 {
		return contextNoneValue
	}
	parts := make([]string, len(ids))
	for i, id := range ids {
		value := fmt.Sprintf("#%d", id)
		if tty {
			value = pterm.ThemeDefault.PrimaryStyle.Sprint(value)
		}
		parts[i] = value
	}
	return strings.Join(parts, ", ")
}

func formatContextTag(value, prefix string, tty bool) string {
	if value == "" {
		return contextNoneValue
	}
	formatted := prefix + value
	if !tty {
		return formatted
	}
	if prefix == "@" {
		return pterm.ThemeDefault.SuccessMessageStyle.Sprint(formatted)
	}
	return pterm.ThemeDefault.PrimaryStyle.Sprint(formatted)
}

func formatViewLabel(label string, tty bool) string {
	if !tty {
		return label
	}
	return pterm.ThemeDefault.SuccessMessageStyle.Sprint(label)
}

func formatViewDescription(description string, tty bool) string {
	if !tty {
		return description
	}
	return pterm.ThemeDefault.DefaultText.Sprint(description)
}

func formatViewUsage(usage string, tty bool) string {
	if !tty {
		return usage
	}
	return pterm.ThemeDefault.InfoMessageStyle.Sprint(usage)
}

func (w Writer) isHumanMode() bool {
	return w.TTY
}

func writeJSON(out io.Writer, payload any) error {
	enc := json.NewEncoder(out)
	return enc.Encode(payload)
}

func formatDate(val *time.Time) string {
	if val == nil {
		return ""
	}
	return val.Format("2006-01-02")
}

func formatDateTimePtr(val *time.Time) string {
	if val == nil {
		return ""
	}
	return formatDateTime(*val)
}

func formatDateTime(val time.Time) string {
	if val.IsZero() {
		return ""
	}
	return val.UTC().Format(time.RFC3339)
}

func (w Writer) formatDateWithFormatter(val *time.Time) string {
	if val == nil {
		return ""
	}
	formatter := w.formatter
	if formatter == nil {
		layout := "2006-01-02 15:04"
		return val.UTC().Format(layout)
	}
	year, month, day := val.UTC().Date()
	date := time.Date(year, month, day, 0, 0, 0, 0, formatter.location)
	return date.Format(formatter.layout)
}

// HistoryJSON represents a shell history entry for JSON output.
type HistoryJSON struct {
	ID      int64  `json:"id"`
	Time    string `json:"time"`
	Command string `json:"command"`
	Success bool   `json:"success"`
	Summary string `json:"summary,omitempty"`
	Intent  string `json:"intent,omitempty"`
}

// HistoryEntry represents a shell history entry for display.
type HistoryEntry struct {
	ID      int64
	Time    time.Time
	Command string
	Success bool
	Summary string
	Intent  string
}

func (w Writer) WriteHistory(entries []*HistoryEntry) error {
	if w.JSON {
		payload := make([]HistoryJSON, 0, len(entries))
		for _, e := range entries {
			payload = append(payload, HistoryJSON{
				ID:      e.ID,
				Time:    formatDateTime(e.Time),
				Command: e.Command,
				Success: e.Success,
				Summary: e.Summary,
				Intent:  e.Intent,
			})
		}
		return writeJSON(w.Out, payload)
	}

	if w.isHumanMode() {
		return w.writeHumanHistory(entries)
	}

	for _, e := range entries {
		status := "✓"
		if !e.Success {
			status = "✗"
		}
		_, err := fmt.Fprintf(w.Out, "%d\t%s\t%s\t%s\t%s\n",
			e.ID,
			w.formatter.Format(e.Time),
			status,
			e.Intent,
			e.Command,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

type TaskVersionChangeJSON struct {
	Type  string `json:"type"`
	Field string `json:"field"`
	Old   string `json:"old,omitempty"`
	New   string `json:"new,omitempty"`
}

type TaskVersionDiffJSON struct {
	VersionID int64                   `json:"versionId"`
	UpdatedAt string                  `json:"updatedAt"`
	Deleted   bool                    `json:"deleted"`
	Changes   []TaskVersionChangeJSON `json:"changes"`
}

type TaskVersionChange struct {
	Type  string
	Field string
	Old   string
	New   string
}

const (
	changeTypeAdd    = "add"
	changeTypeRemove = "remove"
)

//nolint:gocognit // Handles JSON and human rendering paths for history diff output.
func (w Writer) WriteTaskVersionDiff(versions []*store.TaskVersion) error {
	if w.JSON {
		payload := make([]TaskVersionDiffJSON, 0, len(versions))
		for i, current := range versions {
			var prev *store.TaskVersion
			if i+1 < len(versions) {
				prev = versions[i+1]
			}
			changes := diffTaskVersion(prev, current)
			jsonChanges := make([]TaskVersionChangeJSON, 0, len(changes))
			for _, change := range changes {
				jsonChanges = append(jsonChanges, TaskVersionChangeJSON(change))
			}
			payload = append(payload, TaskVersionDiffJSON{
				VersionID: current.VersionID,
				UpdatedAt: formatDateTime(current.UpdatedAt),
				Deleted:   current.Deleted,
				Changes:   jsonChanges,
			})
		}
		return writeJSON(w.Out, payload)
	}

	if w.isHumanMode() {
		return w.writeHumanTaskVersionDiff(versions)
	}

	for i, current := range versions {
		if _, err := fmt.Fprintf(
			w.Out,
			"version %d %s\n",
			current.VersionID,
			w.formatter.Format(current.UpdatedAt),
		); err != nil {
			return err
		}
		var prev *store.TaskVersion
		if i+1 < len(versions) {
			prev = versions[i+1]
		}
		for _, change := range diffTaskVersion(prev, current) {
			prefix := "~"
			if change.Type == changeTypeAdd {
				prefix = "+"
			}
			if change.Type == changeTypeRemove {
				prefix = "-"
			}
			if _, err := fmt.Fprintf(
				w.Out,
				"  %s %s: %s -> %s\n",
				prefix,
				change.Field,
				change.Old,
				change.New,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func diffTaskVersion(prev *store.TaskVersion, current *store.TaskVersion) []TaskVersionChange {
	changes := make([]TaskVersionChange, 0)
	old := &store.TaskVersion{}
	if prev != nil {
		old = prev
	}

	appendScalarChange(&changes, "state", string(old.State), string(current.State))
	appendScalarChange(&changes, "title", old.Title, current.Title)
	appendScalarChange(&changes, "notes", old.Notes, current.Notes)
	appendScalarChange(&changes, "due", formatDate(old.DueOn), formatDate(current.DueOn))
	appendScalarChange(&changes, "waiting_for", old.WaitingFor, current.WaitingFor)
	appendScalarChange(&changes, "deleted", strconv.FormatBool(old.Deleted), strconv.FormatBool(current.Deleted))

	diffListChange(&changes, "project", old.Projects, current.Projects)
	diffListChange(&changes, "context", old.Contexts, current.Contexts)
	diffMetaChange(&changes, old.Meta, current.Meta)

	if len(changes) == 0 {
		changes = append(changes, TaskVersionChange{Type: "none", Field: "snapshot", New: "no visible field changes"})
	}

	return changes
}

func appendScalarChange(changes *[]TaskVersionChange, field, oldVal, newVal string) {
	if oldVal == newVal {
		return
	}
	typeName := "change"
	if oldVal == "" && newVal != "" {
		typeName = changeTypeAdd
	}
	if oldVal != "" && newVal == "" {
		typeName = changeTypeRemove
	}
	*changes = append(*changes, TaskVersionChange{Type: typeName, Field: field, Old: oldVal, New: newVal})
}

func diffListChange(changes *[]TaskVersionChange, field string, oldVals, newVals []string) {
	oldSet := make(map[string]struct{}, len(oldVals))
	for _, value := range oldVals {
		oldSet[value] = struct{}{}
	}
	newSet := make(map[string]struct{}, len(newVals))
	for _, value := range newVals {
		newSet[value] = struct{}{}
	}

	added := make([]string, 0)
	removed := make([]string, 0)
	for value := range newSet {
		if _, ok := oldSet[value]; !ok {
			added = append(added, value)
		}
	}
	for value := range oldSet {
		if _, ok := newSet[value]; !ok {
			removed = append(removed, value)
		}
	}
	sort.Strings(added)
	sort.Strings(removed)

	for _, value := range added {
		*changes = append(*changes, TaskVersionChange{Type: changeTypeAdd, Field: field, New: value})
	}
	for _, value := range removed {
		*changes = append(*changes, TaskVersionChange{Type: changeTypeRemove, Field: field, Old: value})
	}
}

func diffMetaChange(changes *[]TaskVersionChange, oldMeta, newMeta map[string]string) {
	if oldMeta == nil {
		oldMeta = map[string]string{}
	}
	if newMeta == nil {
		newMeta = map[string]string{}
	}

	keys := make([]string, 0, len(oldMeta)+len(newMeta))
	seen := map[string]struct{}{}
	for key := range oldMeta {
		seen[key] = struct{}{}
		keys = append(keys, key)
	}
	for key := range newMeta {
		if _, ok := seen[key]; ok {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		appendScalarChange(changes, "meta."+key, oldMeta[key], newMeta[key])
	}
}
