//nolint:testpackage // Tests validate internal command mode behavior directly.
package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestHandleCommandInputEscCancels(t *testing.T) {
	m := newModel(nil, Options{})
	m.commandMode = true

	updatedModel, _ := m.handleCommandInput(tea.KeyMsg{Type: tea.KeyEsc})
	updated, ok := updatedModel.(model)
	if !ok {
		t.Fatalf("unexpected model type %T", updatedModel)
	}

	if updated.commandMode {
		t.Fatal("expected command mode disabled after esc")
	}
	if updated.status != "command cancelled" {
		t.Fatalf("unexpected status: %q", updated.status)
	}
}

func TestSubmitCommandInputAppliesFilters(t *testing.T) {
	m := newModel(nil, Options{})
	m.commandMode = true
	m.commandInput.SetValue("find state:now and project:work and due:today")

	updatedModel, _ := m.submitCommandInput()
	updated, ok := updatedModel.(model)
	if !ok {
		t.Fatalf("unexpected model type %T", updatedModel)
	}

	if updated.filters.state != "now" {
		t.Fatalf("state = %q, want now", updated.filters.state)
	}
	if updated.filters.project != "work" {
		t.Fatalf("project = %q, want work", updated.filters.project)
	}
	if !updated.filters.dueOnly {
		t.Fatal("dueOnly = false, want true")
	}
	if updated.tabSelected != 1 {
		t.Fatalf("tabSelected = %d, want 1 (Now)", updated.tabSelected)
	}
}

func TestSubmitCommandInputRequiresSelectedTarget(t *testing.T) {
	m := newModel(nil, Options{})
	m.commandMode = true
	m.commandInput.SetValue("set selected state:now")

	updatedModel, _ := m.submitCommandInput()
	updated, ok := updatedModel.(model)
	if !ok {
		t.Fatalf("unexpected model type %T", updatedModel)
	}

	if updated.errText == "" {
		t.Fatal("expected parse/compile error text")
	}
	if !strings.Contains(updated.errText, "SelectedTaskID") {
		t.Fatalf("unexpected error text: %q", updated.errText)
	}
}
