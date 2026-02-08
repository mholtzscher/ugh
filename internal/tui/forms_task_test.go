//nolint:testpackage // Tests validate internal task form behavior directly.
package tui

import "testing"

func TestTaskFormCommitNotesFromTextarea(t *testing.T) {
	form := startAddTaskForm(40).withField(taskFormFieldNotes)
	form.notes.SetValue("line one\nline two")

	updated := form.commitInput()
	if updated.values.notes != "line one\nline two" {
		t.Fatalf("unexpected notes value: %q", updated.values.notes)
	}
}
