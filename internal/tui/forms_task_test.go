//nolint:testpackage // Tests validate internal task form behavior directly.
package tui

import (
	"strings"
	"testing"
)

func TestTaskFormCommitNotesFromTextarea(t *testing.T) {
	form := startAddTaskForm(40).withField(taskFormFieldNotes)
	form.notes.SetValue("line one\nline two")

	updated := form.commitInput()
	if updated.values.notes != "line one\nline two" {
		t.Fatalf("unexpected notes value: %q", updated.values.notes)
	}
}

func TestTaskFormRenderShowsActiveInputInline(t *testing.T) {
	form := startAddTaskForm(40)
	form.values.notes = "keep me"
	form, _ = form.startEditing()

	view := form.render(newStyles(SelectTheme(""), true))
	if strings.Contains(view, "title: -") {
		t.Fatalf("expected title summary to be replaced by inline input, got %q", view)
	}
	if strings.Index(view, "title: ") > strings.Index(view, "notes: keep me") {
		t.Fatalf("expected inline title editor to appear before notes line, got %q", view)
	}
}

func TestTaskFormRenderShowsNotesInputInline(t *testing.T) {
	form := startAddTaskForm(40).withField(taskFormFieldNotes)
	form, _ = form.startEditing()

	view := form.render(newStyles(SelectTheme(""), true))
	if strings.Contains(view, "notes: -") {
		t.Fatalf("expected notes summary to be replaced by inline editor, got %q", view)
	}
	if !strings.Contains(view, "notes: ") {
		t.Fatalf("expected inline notes editor prompt, got %q", view)
	}
}
