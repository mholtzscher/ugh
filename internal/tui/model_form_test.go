//nolint:testpackage // Tests validate internal form-key handling directly.
package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
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
