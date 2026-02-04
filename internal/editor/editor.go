package editor

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/mholtzscher/ugh/internal/store"
)

type TaskTOML struct {
	Title      string            `toml:"title"`
	Notes      string            `toml:"notes,omitempty"`
	Status     string            `toml:"status"`
	Priority   string            `toml:"priority,omitempty"`
	Done       bool              `toml:"done"`
	DueOn      string            `toml:"due_on,omitempty"`
	DeferUntil string            `toml:"defer_until,omitempty"`
	WaitingFor string            `toml:"waiting_for,omitempty"`
	Projects   []string          `toml:"projects,omitempty"`
	Contexts   []string          `toml:"contexts,omitempty"`
	Meta       map[string]string `toml:"meta,omitempty"`
}

func TaskToTOML(task *store.Task) TaskTOML {
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

	return TaskTOML{
		Title:      task.Title,
		Notes:      task.Notes,
		Status:     string(task.Status),
		Priority:   task.Priority,
		Done:       task.Done,
		DueOn:      formatDay(task.DueOn),
		DeferUntil: formatDay(task.DeferUntil),
		WaitingFor: task.WaitingFor,
		Projects:   projects,
		Contexts:   contexts,
		Meta:       meta,
	}
}

const tomlHeader = `# Task %d - Edit and save to apply changes
# Lines starting with # are ignored
#
# Fields:
#   title        - The action title (required)
#   notes        - Optional notes
#   status       - inbox|next|waiting|someday
#   priority     - A-Z (or empty to remove)
#   done         - true/false
#   due_on       - YYYY-MM-DD
#   defer_until  - YYYY-MM-DD
#   waiting_for  - Optional string
#   projects     - List of project names
#   contexts     - List of context names
#   meta         - Key-value pairs

`

func Edit(task *store.Task) (*TaskTOML, error) {
	taskTOML := TaskToTOML(task)

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf(tomlHeader, task.ID))
	if err := toml.NewEncoder(&buf).Encode(taskTOML); err != nil {
		return nil, fmt.Errorf("encode task to TOML: %w", err)
	}
	original := buf.String()

	tmpFile, err := os.CreateTemp("", "ugh-edit-*.toml")
	if err != nil {
		return nil, fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer func() { _ = os.Remove(tmpPath) }()

	if _, err := tmpFile.WriteString(original); err != nil {
		_ = tmpFile.Close()
		return nil, fmt.Errorf("write temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return nil, fmt.Errorf("close temp file: %w", err)
	}

	editor := getEditor()
	cmd := exec.Command(editor, tmpPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("run editor: %w", err)
	}

	edited, err := os.ReadFile(tmpPath)
	if err != nil {
		return nil, fmt.Errorf("read edited file: %w", err)
	}

	if string(edited) == original {
		return nil, nil
	}

	var result TaskTOML
	if _, err := toml.Decode(string(edited), &result); err != nil {
		return nil, fmt.Errorf("parse edited TOML: %w", err)
	}

	if err := validate(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func getEditor() string {
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	for _, editor := range []string{"vim", "vi", "nano"} {
		if _, err := exec.LookPath(editor); err == nil {
			return editor
		}
	}
	return "vi"
}

func validate(t *TaskTOML) error {
	t.Title = strings.TrimSpace(t.Title)
	if t.Title == "" {
		return errors.New("title cannot be empty")
	}

	t.Status = strings.ToLower(strings.TrimSpace(t.Status))
	if t.Status == "" {
		t.Status = "inbox"
	}
	switch t.Status {
	case "inbox", "next", "waiting", "someday":
		// ok
	default:
		return fmt.Errorf("invalid status %q: must be inbox|next|waiting|someday", t.Status)
	}

	t.Priority = strings.ToUpper(strings.TrimSpace(t.Priority))
	if t.Priority != "" {
		if len(t.Priority) != 1 || t.Priority[0] < 'A' || t.Priority[0] > 'Z' {
			return fmt.Errorf("invalid priority %q: must be A-Z", t.Priority)
		}
	}

	t.DueOn = strings.TrimSpace(t.DueOn)
	if t.DueOn != "" {
		if _, err := time.Parse("2006-01-02", t.DueOn); err != nil {
			return fmt.Errorf("invalid due_on %q: expected YYYY-MM-DD", t.DueOn)
		}
	}
	t.DeferUntil = strings.TrimSpace(t.DeferUntil)
	if t.DeferUntil != "" {
		if _, err := time.Parse("2006-01-02", t.DeferUntil); err != nil {
			return fmt.Errorf("invalid defer_until %q: expected YYYY-MM-DD", t.DeferUntil)
		}
	}
	t.WaitingFor = strings.TrimSpace(t.WaitingFor)

	t.Projects = cleanTags(t.Projects)
	t.Contexts = cleanTags(t.Contexts)

	return nil
}

func formatDay(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.UTC().Format("2006-01-02")
}

func cleanTags(tags []string) []string {
	if len(tags) == 0 {
		return tags
	}
	result := make([]string, 0, len(tags))
	seen := make(map[string]bool)
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		// Remove leading + or @ if user accidentally added them
		tag = strings.TrimPrefix(tag, "+")
		tag = strings.TrimPrefix(tag, "@")
		if tag != "" && !seen[tag] {
			result = append(result, tag)
			seen[tag] = true
		}
	}
	return result
}
