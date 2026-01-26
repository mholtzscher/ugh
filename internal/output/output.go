package output

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mholtzscher/ugh/internal/store"
	"github.com/mholtzscher/ugh/internal/todotxt"
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
	_, err := fmt.Fprintln(w.Out, todoLine(task))
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
		if _, err := fmt.Fprintln(w.Out, todoLine(task)); err != nil {
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
	ID             int64             `json:"id"`
	Done           bool              `json:"done"`
	Priority       string            `json:"priority,omitempty"`
	CompletionDate string            `json:"completionDate,omitempty"`
	CreationDate   string            `json:"creationDate,omitempty"`
	Description    string            `json:"description"`
	Projects       []string          `json:"projects"`
	Contexts       []string          `json:"contexts"`
	Meta           map[string]string `json:"meta"`
	Unknown        []string          `json:"unknown"`
	TodoTxt        string            `json:"todoTxt"`
	CreatedAt      string            `json:"createdAt"`
	UpdatedAt      string            `json:"updatedAt"`
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
	unknown := task.Unknown
	if unknown == nil {
		unknown = []string{}
	}
	meta := task.Meta
	if meta == nil {
		meta = map[string]string{}
	}
	return TaskJSON{
		ID:             task.ID,
		Done:           task.Done,
		Priority:       task.Priority,
		CompletionDate: formatDate(task.CompletionDate),
		CreationDate:   formatDate(task.CreationDate),
		Description:    task.Description,
		Projects:       projects,
		Contexts:       contexts,
		Meta:           meta,
		Unknown:        unknown,
		TodoTxt:        todoLine(task),
		CreatedAt:      formatDateTime(task.CreatedAt),
		UpdatedAt:      formatDateTime(task.UpdatedAt),
	}
}

func todoLine(task *store.Task) string {
	parsed := todotxt.Parsed{
		Done:           task.Done,
		Priority:       task.Priority,
		CompletionDate: task.CompletionDate,
		CreationDate:   task.CreationDate,
		Description:    task.Description,
		Projects:       task.Projects,
		Contexts:       task.Contexts,
		Meta:           task.Meta,
		Unknown:        task.Unknown,
	}
	return todotxt.Format(parsed)
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

func formatDateTime(val time.Time) string {
	if val.IsZero() {
		return ""
	}
	return val.UTC().Format(time.RFC3339)
}
