package output

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/mholtzscher/ugh/internal/store"
	"golang.org/x/term"
)

type Writer struct {
	Out     io.Writer
	JSON    bool
	NoColor bool
	TTY     bool
}

func NewWriter(jsonMode bool, noColor bool) Writer {
	return Writer{
		Out:     os.Stdout,
		JSON:    jsonMode,
		NoColor: noColor,
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
		return writeHumanTask(w.Out, task)
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
		return writeHumanList(w.Out, tasks)
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
		return writeHumanSummary(w.Out, summary)
	}
	return writePlainSummary(w.Out, summary)
}

type TaskJSON struct {
	ID          int64             `json:"id"`
	Done        bool              `json:"done"`
	Status      string            `json:"status"`
	Priority    string            `json:"priority,omitempty"`
	Title       string            `json:"title"`
	Notes       string            `json:"notes,omitempty"`
	DueOn       string            `json:"dueOn,omitempty"`
	DeferUntil  string            `json:"deferUntil,omitempty"`
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
		Done:        task.Done,
		Status:      string(task.Status),
		Priority:    task.Priority,
		Title:       task.Title,
		Notes:       task.Notes,
		DueOn:       formatDate(task.DueOn),
		DeferUntil:  formatDate(task.DeferUntil),
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
	deferUntil := formatDate(task.DeferUntil)
	fields := []string{
		fmt.Sprintf("%d", task.ID),
		string(task.Status),
		due,
		deferUntil,
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
