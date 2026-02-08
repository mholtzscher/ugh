package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mholtzscher/ugh/internal/service"
	"github.com/mholtzscher/ugh/internal/store"
)

type taskFormMode int

const (
	taskFormNone taskFormMode = iota
	taskFormAdd
	taskFormEdit
)

type taskFormField int

const (
	taskFormFieldTitle taskFormField = iota
	taskFormFieldNotes
	taskFormFieldDue
	taskFormFieldWaitingFor
	taskFormFieldProjects
	taskFormFieldContexts
)

const taskFormFieldCount = 6

const (
	formPlaceholderRequired = "required"
	formPlaceholderOptional = "optional"
	formPlaceholderCSV      = "comma,separated"
	formNotesHeight         = 4
)

type taskFormValues struct {
	title      string
	notes      string
	due        string
	waitingFor string
	projects   string
	contexts   string
}

type taskFormState struct {
	mode      taskFormMode
	field     taskFormField
	taskID    int64
	taskState string
	values    taskFormValues
	input     textinput.Model
	notes     textarea.Model
}

func inactiveTaskForm(width int) taskFormState {
	input := textinput.New()
	input.Prompt = "> "
	input.CharLimit = 1024
	input.Width = width

	notes := textarea.New()
	notes.Prompt = "> "
	notes.Placeholder = formPlaceholderOptional
	notes.CharLimit = 4096
	notes.ShowLineNumbers = false
	notes.SetWidth(width)
	notes.SetHeight(formNotesHeight)

	return taskFormState{mode: taskFormNone, input: input, notes: notes}
}

func startAddTaskForm(width int) taskFormState {
	form := inactiveTaskForm(width)
	form.mode = taskFormAdd
	form.taskState = string(store.StateInbox)
	return form.withField(taskFormFieldTitle)
}

func startEditTaskForm(task *store.Task, width int) taskFormState {
	form := inactiveTaskForm(width)
	form.mode = taskFormEdit
	if task == nil {
		return form
	}

	form.taskID = task.ID
	form.taskState = string(task.State)
	form.values = taskFormValues{
		title:      task.Title,
		notes:      task.Notes,
		due:        dueText(task),
		waitingFor: task.WaitingFor,
		projects:   strings.Join(task.Projects, ","),
		contexts:   strings.Join(task.Contexts, ","),
	}
	if form.values.due == "-" {
		form.values.due = ""
	}

	return form.withField(taskFormFieldTitle)
}

func (f taskFormState) active() bool {
	return f.mode != taskFormNone
}

func (f taskFormState) withWidth(width int) taskFormState {
	f.input.Width = width
	f.notes.SetWidth(width)
	return f
}

func (f taskFormState) withField(field taskFormField) taskFormState {
	f.field = field
	if field == taskFormFieldNotes {
		f.notes.Prompt = taskFormFieldLabel(field) + ": "
		f.notes.Placeholder = taskFormFieldPlaceholder(field)
		f.notes.SetValue(f.valueForField(field))
		f.notes.CursorEnd()
		return f
	}

	f.input.Prompt = taskFormFieldLabel(field) + ": "
	f.input.Placeholder = taskFormFieldPlaceholder(field)
	f.input.SetValue(f.valueForField(field))
	f.input.CursorEnd()
	return f
}

func (f taskFormState) focus() tea.Cmd {
	if f.field == taskFormFieldNotes {
		return f.notes.Focus()
	}
	return f.input.Focus()
}

func (f taskFormState) update(msg tea.Msg) (taskFormState, tea.Cmd) {
	if f.field == taskFormFieldNotes {
		var cmd tea.Cmd
		f.notes, cmd = f.notes.Update(msg)
		return f, cmd
	}

	var cmd tea.Cmd
	f.input, cmd = f.input.Update(msg)
	return f, cmd
}

func (f taskFormState) activeInputValue() string {
	if f.field == taskFormFieldNotes {
		return f.notes.Value()
	}
	return f.input.Value()
}

func (f taskFormState) inputView(styleSet styles) string {
	if f.field == taskFormFieldNotes {
		notes := f.notes
		notes.FocusedStyle.Prompt = styleSet.key
		notes.FocusedStyle.Placeholder = styleSet.muted
		notes.FocusedStyle.Text = lipgloss.NewStyle()
		notes.BlurredStyle.Prompt = styleSet.key
		notes.BlurredStyle.Placeholder = styleSet.muted
		notes.BlurredStyle.Text = styleSet.muted
		return notes.View()
	}

	input := f.input
	input.PromptStyle = styleSet.key
	input.PlaceholderStyle = styleSet.muted
	input.TextStyle = lipgloss.NewStyle()
	return input.View()
}

func (f taskFormState) valueForField(field taskFormField) string {
	switch field {
	case taskFormFieldTitle:
		return f.values.title
	case taskFormFieldNotes:
		return f.values.notes
	case taskFormFieldDue:
		return f.values.due
	case taskFormFieldWaitingFor:
		return f.values.waitingFor
	case taskFormFieldProjects:
		return f.values.projects
	case taskFormFieldContexts:
		return f.values.contexts
	default:
		return ""
	}
}

func (f taskFormState) commitInput() taskFormState {
	value := strings.TrimSpace(f.activeInputValue())
	switch f.field {
	case taskFormFieldTitle:
		f.values.title = value
	case taskFormFieldNotes:
		f.values.notes = value
	case taskFormFieldDue:
		f.values.due = value
	case taskFormFieldWaitingFor:
		f.values.waitingFor = value
	case taskFormFieldProjects:
		f.values.projects = value
	case taskFormFieldContexts:
		f.values.contexts = value
	}
	return f
}

func (f taskFormState) nextField() (taskFormState, bool) {
	if f.field >= taskFormFieldCount-1 {
		return f, true
	}
	return f.withField(f.field + 1), false
}

func (f taskFormState) previousField() taskFormState {
	if f.field <= 0 {
		return f.withField(taskFormFieldTitle)
	}
	return f.withField(f.field - 1)
}

func (f taskFormState) modeTitle() string {
	if f.mode == taskFormEdit {
		return "Edit Task"
	}
	return "Add Task"
}

func (f taskFormState) render(styles styles) string {
	lines := []string{styles.title.Render(f.modeTitle())}
	for field := range taskFormFieldCount {
		item := taskFormField(field)
		label := taskFormFieldLabel(item)
		value := f.valueForField(item)
		if value == "" {
			value = "-"
		}
		line := label + ": " + value
		if item == f.field {
			line = styles.selected.Render(line)
		}
		lines = append(lines, line)
	}
	lines = append(lines, "", f.inputView(styles))
	hint := "enter/tab: next/save  shift+tab: previous  esc: cancel"
	if f.field == taskFormFieldNotes {
		hint = "enter: newline  tab: next  shift+tab: previous  esc: cancel"
	}
	lines = append(lines, styles.muted.Render(hint))
	return strings.Join(lines, "\n")
}

func (f taskFormState) createRequest() service.CreateTaskRequest {
	return service.CreateTaskRequest{
		Title:      f.values.title,
		Notes:      f.values.notes,
		State:      f.taskState,
		Projects:   splitCSV(f.values.projects),
		Contexts:   splitCSV(f.values.contexts),
		DueOn:      f.values.due,
		WaitingFor: f.values.waitingFor,
	}
}

func (f taskFormState) fullUpdateRequest() service.FullUpdateTaskRequest {
	return service.FullUpdateTaskRequest{
		ID:         f.taskID,
		Title:      f.values.title,
		Notes:      f.values.notes,
		State:      f.taskState,
		Projects:   splitCSV(f.values.projects),
		Contexts:   splitCSV(f.values.contexts),
		Meta:       map[string]string{},
		DueOn:      f.values.due,
		WaitingFor: f.values.waitingFor,
	}
}

func taskFormFieldLabel(field taskFormField) string {
	switch field {
	case taskFormFieldTitle:
		return "title"
	case taskFormFieldNotes:
		return "notes"
	case taskFormFieldDue:
		return "due (YYYY-MM-DD)"
	case taskFormFieldWaitingFor:
		return "waiting-for"
	case taskFormFieldProjects:
		return "projects (csv)"
	case taskFormFieldContexts:
		return "contexts (csv)"
	default:
		return "field"
	}
}

func taskFormFieldPlaceholder(field taskFormField) string {
	switch field {
	case taskFormFieldTitle:
		return formPlaceholderRequired
	case taskFormFieldDue:
		return formPlaceholderOptional
	case taskFormFieldNotes, taskFormFieldWaitingFor:
		return formPlaceholderOptional
	case taskFormFieldProjects, taskFormFieldContexts:
		return formPlaceholderCSV
	default:
		return formPlaceholderOptional
	}
}

func splitCSV(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	seen := map[string]bool{}
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" || seen[trimmed] {
			continue
		}
		seen[trimmed] = true
		result = append(result, trimmed)
	}
	return result
}
