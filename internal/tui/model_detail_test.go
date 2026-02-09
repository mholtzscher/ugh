//nolint:testpackage // Tests validate internal detail viewport behavior directly.
package tui

import (
	"strings"
	"testing"

	"github.com/mholtzscher/ugh/internal/store"
)

func TestScrollDetailDownAndUp(t *testing.T) {
	m := newModel(nil, Options{})
	m.tasks = []*store.Task{{
		ID:    1,
		State: store.StateNow,
		Title: "scroll me",
		Notes: strings.Repeat("line\n", 64),
	}}
	m.selected = 0
	m.detail.Width = 40
	m.detail.Height = 6
	m.detail.SetContent(m.renderTaskDetailContent())

	before := m.detail.YOffset
	afterDownModel, _ := m.scrollDetailDown()
	afterDown := afterDownModel.(model)
	if afterDown.detail.YOffset <= before {
		t.Fatalf("expected detail offset to increase, before=%d after=%d", before, afterDown.detail.YOffset)
	}

	afterUpModel, _ := afterDown.scrollDetailUp()
	afterUp := afterUpModel.(model)
	if afterUp.detail.YOffset >= afterDown.detail.YOffset {
		t.Fatalf(
			"expected detail offset to decrease, before=%d after=%d",
			afterDown.detail.YOffset,
			afterUp.detail.YOffset,
		)
	}
}

func TestScrollDetailNoSelectedTaskNoop(t *testing.T) {
	m := newModel(nil, Options{})
	m.detail.Width = 40
	m.detail.Height = 6
	m.detail.SetContent("a\nb\nc\nd\ne\nf\ng")

	before := m.detail.YOffset
	updatedModel, _ := m.scrollDetailDown()
	updated := updatedModel.(model)
	if updated.detail.YOffset != before {
		t.Fatalf("expected no offset change without selected task, before=%d after=%d", before, updated.detail.YOffset)
	}
}

func TestMoveSelectionResetsDetailOffset(t *testing.T) {
	m := newModel(nil, Options{})
	m.tasks = []*store.Task{
		{ID: 1, State: store.StateNow, Title: "one", Notes: strings.Repeat("line\n", 64)},
		{ID: 2, State: store.StateNow, Title: "two", Notes: strings.Repeat("line\n", 64)},
	}
	m.selected = 0
	m.detail.Width = 40
	m.detail.Height = 6
	m.detail.SetContent(m.renderTaskDetailContent())
	m.detail.HalfPageDown()

	updatedModel, _ := m.moveCurrentSelection(1)
	updated := updatedModel.(model)
	if updated.detail.YOffset != 0 {
		t.Fatalf("expected detail offset reset on selection change, got %d", updated.detail.YOffset)
	}
}

func TestApplyTasksLoadedResetsDetailOffsetWhenSelectedTaskChanges(t *testing.T) {
	m := newModel(nil, Options{})
	m.tasks = []*store.Task{{ID: 10, State: store.StateNow, Title: "old"}}
	m.selected = 0
	m.detail.Width = 40
	m.detail.Height = 6
	m.detail.SetContent(strings.Repeat("line\n", 64))
	m.detail.HalfPageDown()

	updated := m.applyTasksLoaded(tasksLoadedMsg{tasks: []*store.Task{{ID: 20, State: store.StateNow, Title: "new"}}})
	if updated.detail.YOffset != 0 {
		t.Fatalf("expected detail offset reset after task switch, got %d", updated.detail.YOffset)
	}
}

func TestRenderTaskDetailContentShowsInlineFormWhenActive(t *testing.T) {
	m := newModel(nil, Options{})
	m.taskForm = startAddTaskForm(searchInputWidth)
	m.detail.Width = 40

	content := m.renderTaskDetailContent()
	if !strings.Contains(content, "Add Task") {
		t.Fatalf("expected inline form content, got %q", content)
	}
}

func TestViewKeepsTaskListVisibleWhenTaskFormActive(t *testing.T) {
	m := newModel(nil, Options{})
	m.viewportW = 120
	m.viewportH = 32
	m.layout = calculateLayout(m.viewportW, m.viewportH)
	m.tasks = []*store.Task{{ID: 1, State: store.StateInbox, Title: "first"}}
	m.selected = 0
	m.taskForm = startAddTaskForm(searchInputWidth)

	view := m.View()
	if !strings.Contains(view, "first") {
		t.Fatalf("expected task list to remain visible, got %q", view)
	}
}

func TestFocusedPaneModeDefaultsToTasks(t *testing.T) {
	m := newModel(nil, Options{})
	if mode := m.focusedPaneMode(); mode != paneFocusTasks {
		t.Fatalf("expected tasks pane focus, got %v", mode)
	}
}

func TestFocusedPaneModeUsesDetailForFormNavigation(t *testing.T) {
	m := newModel(nil, Options{})
	m.taskForm = startAddTaskForm(searchInputWidth)
	if mode := m.focusedPaneMode(); mode != paneFocusDetail {
		t.Fatalf("expected detail pane focus, got %v", mode)
	}
}

func TestFocusedPaneModeUsesEditColorInFormEditing(t *testing.T) {
	m := newModel(nil, Options{})
	m.taskForm = startAddTaskForm(searchInputWidth)
	m.taskForm, _ = m.taskForm.startEditing()
	if mode := m.focusedPaneMode(); mode != paneFocusDetailEdit {
		t.Fatalf("expected edit pane focus, got %v", mode)
	}
}
