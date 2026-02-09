package output

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pterm/pterm"
	"golang.org/x/term"

	"github.com/mholtzscher/ugh/internal/store"
)

type Writer struct {
	Out     io.Writer
	JSON    bool
	NoColor bool
	TTY     bool
}

type KeyValue struct {
	Key   string
	Value string
}

func NewWriter(jsonMode bool, noColor bool) Writer {
	return Writer{
		Out:     os.Stdout,
		JSON:    jsonMode,
		NoColor: noColor || os.Getenv("NO_COLOR") != "",
		TTY:     term.IsTerminal(int(os.Stdout.Fd())),
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

	if w.TTY {
		return writeHumanTask(w.Out, w.NoColor, task)
	}
	_, err := fmt.Fprintln(w.Out, plainLine(task))
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

	if w.TTY {
		return writeHumanList(w.Out, w.NoColor, tasks)
	}
	for _, task := range tasks {
		if _, err := fmt.Fprintln(w.Out, plainLine(task)); err != nil {
			return err
		}
	}
	return nil
}

func (w Writer) WriteTags(tags []store.NameCount) error {
	if w.JSON {
		return writeJSON(w.Out, tags)
	}

	if w.TTY {
		return writeHumanTags(w.Out, w.NoColor, tags)
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
	if w.TTY {
		return writeHumanTagsWithCounts(w.Out, w.NoColor, tags)
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
	if w.TTY {
		return writeHumanSummary(w.Out, w.NoColor, summary)
	}
	return writePlainSummary(w.Out, summary)
}

func (w Writer) WriteKeyValues(rows []KeyValue) error {
	if w.TTY {
		return writeHumanKeyValues(w.Out, w.NoColor, rows)
	}
	for _, row := range rows {
		if _, err := fmt.Fprintf(w.Out, "%s:\t%s\n", row.Key, row.Value); err != nil {
			return err
		}
	}
	return nil
}

func (w Writer) WriteLine(line string) error {
	if !w.TTY {
		_, err := fmt.Fprintln(w.Out, line)
		return err
	}
	formatted := pterm.DefaultBasicText.Sprintln(line)
	return writeRenderedLine(w.Out, w.NoColor, formatted)
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

func (w Writer) writePrefixLine(printer pterm.PrefixPrinter, line string) error {
	if !w.TTY {
		_, err := fmt.Fprintln(w.Out, line)
		return err
	}
	formatted := printer.Sprintln(line)
	return writeRenderedLine(w.Out, w.NoColor, formatted)
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

func plainLine(task *store.Task) string {
	if task == nil {
		return ""
	}
	due := formatDate(task.DueOn)
	fields := []string{
		strconv.FormatInt(task.ID, 10),
		string(task.State),
		due,
		task.WaitingFor,
		task.Title,
	}
	return strings.Join(fields, "\t")
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

	if w.TTY {
		return writeHumanHistory(w.Out, w.NoColor, entries)
	}

	for _, e := range entries {
		status := "✓"
		if !e.Success {
			status = "✗"
		}
		_, err := fmt.Fprintf(w.Out, "%d\t%s\t%s\t%s\t%s\n",
			e.ID,
			formatDateTime(e.Time),
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
