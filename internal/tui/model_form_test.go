//nolint:testpackage // Tests validate internal form-key handling directly.
package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/mholtzscher/ugh/internal/store"
)

func TestTaskFormNotesEnterDoesNotAdvanceField(t *testing.T) {
	m := newModel(nil, Options{})
	m.taskForm = startAddTaskForm(searchInputWidth).withField(taskFormFieldNotes)

	updatedModel, _ := m.handleTaskFormInput(tea.KeyMsg{Type: tea.KeyEnter})
	updated, ok := updatedModel.(model)
	if !ok {
		t.Fatalf("unexpected model type %T", updatedModel)
	}

	if updated.taskForm.field != taskFormFieldNotes {
		t.Fatalf("expected notes field to remain active, got %v", updated.taskForm.field)
	}
	if !updated.taskForm.editing {
		t.Fatal("expected notes field to enter editing mode")
	}
}

func TestTaskFormNotesTabAdvancesField(t *testing.T) {
	m := newModel(nil, Options{})
	m.taskForm = startAddTaskForm(searchInputWidth).withField(taskFormFieldNotes)

	updatedModel, _ := m.handleTaskFormInput(tea.KeyMsg{Type: tea.KeyTab})
	updated, ok := updatedModel.(model)
	if !ok {
		t.Fatalf("unexpected model type %T", updatedModel)
	}

	if updated.taskForm.field != taskFormFieldDue {
		t.Fatalf("expected due field after tab, got %v", updated.taskForm.field)
	}
}

func TestTaskFormEnterStartsEditingAndAdvances(t *testing.T) {
	m := newModel(nil, Options{})
	m.taskForm = startAddTaskForm(searchInputWidth).withField(taskFormFieldTitle)

	startedModel, _ := m.handleTaskFormInput(tea.KeyMsg{Type: tea.KeyEnter})
	started, ok := startedModel.(model)
	if !ok {
		t.Fatalf("unexpected model type %T", startedModel)
	}
	if !started.taskForm.editing {
		t.Fatal("expected form to enter editing mode")
	}

	started.taskForm.input.SetValue("new title")
	nextModel, _ := started.handleTaskFormInput(tea.KeyMsg{Type: tea.KeyEnter})
	next, ok := nextModel.(model)
	if !ok {
		t.Fatalf("unexpected model type %T", nextModel)
	}
	if !next.taskForm.editing {
		t.Fatal("expected form to keep editing mode on next field")
	}
	if next.taskForm.field != taskFormFieldNotes {
		t.Fatalf("expected notes field after enter, got %v", next.taskForm.field)
	}
	if next.taskForm.values.title != "new title" {
		t.Fatalf("expected title commit, got %q", next.taskForm.values.title)
	}
}

func TestTaskFormTypingWorksAfterEnteringEditMode(t *testing.T) {
	m := newModel(nil, Options{})
	m.taskForm = startAddTaskForm(searchInputWidth).withField(taskFormFieldTitle)

	startedModel, _ := m.handleTaskFormInput(tea.KeyMsg{Type: tea.KeyEnter})
	started, ok := startedModel.(model)
	if !ok {
		t.Fatalf("unexpected model type %T", startedModel)
	}

	typedModel, _ := started.handleTaskFormInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	typed, ok := typedModel.(model)
	if !ok {
		t.Fatalf("unexpected model type %T", typedModel)
	}
	if typed.taskForm.input.Value() != "h" {
		t.Fatalf("expected typed value, got %q", typed.taskForm.input.Value())
	}
}

func TestTaskFormEscStopsEditingBeforeCancel(t *testing.T) {
	m := newModel(nil, Options{})
	m.taskForm = startAddTaskForm(searchInputWidth)

	startedModel, _ := m.handleTaskFormInput(tea.KeyMsg{Type: tea.KeyEnter})
	started, ok := startedModel.(model)
	if !ok {
		t.Fatalf("unexpected model type %T", startedModel)
	}
	if !started.taskForm.editing {
		t.Fatal("expected form to enter editing mode")
	}

	stoppedModel, _ := started.handleTaskFormInput(tea.KeyMsg{Type: tea.KeyEsc})
	stopped, ok := stoppedModel.(model)
	if !ok {
		t.Fatalf("unexpected model type %T", stoppedModel)
	}
	if !stopped.taskForm.active() {
		t.Fatal("expected form to remain active after first esc")
	}
	if stopped.taskForm.editing {
		t.Fatal("expected form to leave editing mode after first esc")
	}

	cancelledModel, _ := stopped.handleTaskFormInput(tea.KeyMsg{Type: tea.KeyEsc})
	cancelled, ok := cancelledModel.(model)
	if !ok {
		t.Fatalf("unexpected model type %T", cancelledModel)
	}
	if cancelled.taskForm.active() {
		t.Fatal("expected form to cancel after second esc")
	}
}

func TestTaskFormJKNavigatesInNavigationMode(t *testing.T) {
	m := newModel(nil, Options{})
	m.taskForm = startAddTaskForm(searchInputWidth).withField(taskFormFieldTitle)

	downModel, _ := m.handleTaskFormInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	down, ok := downModel.(model)
	if !ok {
		t.Fatalf("unexpected model type %T", downModel)
	}
	if down.taskForm.field != taskFormFieldNotes {
		t.Fatalf("expected notes field after j, got %v", down.taskForm.field)
	}

	upModel, _ := down.handleTaskFormInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	up, ok := upModel.(model)
	if !ok {
		t.Fatalf("unexpected model type %T", upModel)
	}
	if up.taskForm.field != taskFormFieldTitle {
		t.Fatalf("expected title field after k, got %v", up.taskForm.field)
	}
}

func TestTaskFormCtrlSSavesWhenTitlePresent(t *testing.T) {
	m := newModel(nil, Options{})
	m.taskForm = startAddTaskForm(searchInputWidth)
	m.taskForm.input.SetValue("ship inline editor")

	updatedModel, _ := m.handleTaskFormInput(tea.KeyMsg{Type: tea.KeyCtrlS})
	updated, ok := updatedModel.(model)
	if !ok {
		t.Fatalf("unexpected model type %T", updatedModel)
	}

	if updated.taskForm.active() {
		t.Fatal("expected form to close after save")
	}
	if !updated.loading {
		t.Fatal("expected save to start async action")
	}
}

func TestTaskFormCtrlSRequiresTitle(t *testing.T) {
	m := newModel(nil, Options{})
	m.taskForm = startAddTaskForm(searchInputWidth)

	updatedModel, _ := m.handleTaskFormInput(tea.KeyMsg{Type: tea.KeyCtrlS})
	updated, ok := updatedModel.(model)
	if !ok {
		t.Fatalf("unexpected model type %T", updatedModel)
	}

	if !updated.taskForm.active() {
		t.Fatal("expected form to remain active")
	}
	if updated.errText != statusTitleRequired {
		t.Fatalf("unexpected error text: %q", updated.errText)
	}
}

func TestStartTaskEditFormBeginsInNavigationMode(t *testing.T) {
	m := newModel(nil, Options{})
	m.tasks = []*store.Task{{ID: 1, State: store.StateInbox, Title: "first"}}
	m.selected = 0

	updatedModel, _ := m.startTaskEditForm()
	updated, ok := updatedModel.(model)
	if !ok {
		t.Fatalf("unexpected model type %T", updatedModel)
	}

	if updated.taskForm.editing {
		t.Fatal("expected form to start in navigation mode")
	}
}
