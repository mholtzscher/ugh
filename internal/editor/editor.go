package editor

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/mholtzscher/ugh/internal/domain"
	"github.com/mholtzscher/ugh/internal/store"
)

//go:embed task.schema.json
var taskSchemaJSON []byte

type TaskTOML struct {
	Title      string            `toml:"title"`
	Notes      string            `toml:"notes,omitempty"`
	State      string            `toml:"state"`
	DueOn      string            `toml:"due_on,omitempty"`
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
		State:      string(task.State),
		DueOn:      formatDay(task.DueOn),
		WaitingFor: task.WaitingFor,
		Projects:   projects,
		Contexts:   contexts,
		Meta:       meta,
	}
}

const taskSchemaFileName = "ugh-task.schema.json"

func taskTOMLHeader(taskID int64) string {
	return fmt.Sprintf(`# Task %d - Edit and save to apply changes
# Lines starting with # are ignored
#
# Fields:
#   title        - The action title (required)
#   notes        - Optional notes
#   state        - %s
#   due_on       - %s
#   waiting_for  - Optional string
#   projects     - List of project names
#   contexts     - List of context names
#   meta         - Key-value pairs

`, taskID, domain.TaskStatesUsage, domain.DateTextYYYYMMDD)
}

//nolint:funlen
func Edit(task *store.Task) (*TaskTOML, bool, error) {
	taskTOML := TaskToTOML(task)

	var buf bytes.Buffer

	tmpDir, err := makeEditTempDir()
	if err != nil {
		return nil, false, fmt.Errorf("create temp dir: %w", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	useSchemaHeader := false
	schemaRef := ""
	if len(taskSchemaJSON) > 0 {
		schemaPath := filepath.Join(tmpDir, taskSchemaFileName)
		err = os.WriteFile(schemaPath, taskSchemaJSON, 0o600)
		if err == nil {
			useSchemaHeader = true
			_ = os.WriteFile(filepath.Join(tmpDir, "taplo.toml"), fmt.Appendf(nil, `include = ["*.toml"]

[schema]
path = "./%s"
enabled = true
`, taskSchemaFileName), 0o600)
			if runtime.GOOS == "windows" {
				schemaRef = "file:///" + filepath.ToSlash(schemaPath)
			} else {
				schemaRef = "file://" + schemaPath
			}
		}
	}

	header := taskTOMLHeader(task.ID)
	if useSchemaHeader {
		header = fmt.Sprintf("#:schema %s\n%s", schemaRef, header)
	}
	buf.WriteString(header)
	err = toml.NewEncoder(&buf).Encode(taskTOML)
	if err != nil {
		return nil, false, fmt.Errorf("encode task to TOML: %w", err)
	}
	original := buf.String()

	tmpFile, err := os.CreateTemp(tmpDir, "ugh-edit-*.toml")
	if err != nil {
		return nil, false, fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	_, err = tmpFile.WriteString(original)
	if err != nil {
		_ = tmpFile.Close()
		return nil, false, fmt.Errorf("write temp file: %w", err)
	}
	err = tmpFile.Close()
	if err != nil {
		return nil, false, fmt.Errorf("close temp file: %w", err)
	}

	editor := getEditor()
	cmd := exec.CommandContext(context.Background(), editor, tmpPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return nil, false, fmt.Errorf("run editor: %w", err)
	}

	edited, err := os.ReadFile(tmpPath)
	if err != nil {
		return nil, false, fmt.Errorf("read edited file: %w", err)
	}

	if string(edited) == original {
		return nil, false, nil
	}

	var result TaskTOML
	_, err = toml.Decode(string(edited), &result)
	if err != nil {
		return nil, false, fmt.Errorf("parse edited TOML: %w", err)
	}

	err = validate(&result)
	if err != nil {
		return nil, false, err
	}

	return &result, true, nil
}

func makeEditTempDir() (string, error) {
	if wd, err := os.Getwd(); err == nil && wd != "" {
		if dir, mkdirErr := os.MkdirTemp(wd, "ugh-edit-"); mkdirErr == nil {
			return dir, nil
		}
	}
	return os.MkdirTemp("", "ugh-edit-")
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

	t.State = strings.ToLower(strings.TrimSpace(t.State))
	if t.State == "" {
		t.State = domain.TaskStateInbox
	}
	if !domain.IsTaskState(t.State) {
		return domain.InvalidStateMustBeError(t.State)
	}

	t.DueOn = strings.TrimSpace(t.DueOn)
	if t.DueOn != "" {
		if _, err := time.Parse(domain.DateLayoutYYYYMMDD, t.DueOn); err != nil {
			return domain.InvalidDueOnFormatError(t.DueOn)
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
