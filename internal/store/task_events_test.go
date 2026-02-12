//nolint:testpackage // Tests validate store internals directly.
package store

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskEventsLifecycleAndDeleteRetention(t *testing.T) {
	t.Parallel()

	ctx := WithAuditOrigin(context.Background(), "test")
	s := openTestStore(t)

	task, err := s.CreateTask(ctx, &Task{
		Title:    "first",
		Notes:    "note one",
		State:    StateInbox,
		Projects: []string{"work"},
		Meta:     map[string]string{"priority": "low"},
	})
	require.NoError(t, err, "create task error")

	due := time.Date(2026, time.March, 14, 0, 0, 0, 0, time.UTC)
	task.Title = "first updated"
	task.Notes = "note two"
	task.State = StateNow
	task.DueOn = &due
	task.Projects = []string{"home", "work"}
	task.Meta = map[string]string{"priority": "high", "owner": "me"}
	_, err = s.UpdateTask(ctx, task)
	require.NoError(t, err, "update task error")

	_, err = s.SetDone(ctx, []int64{task.ID}, true)
	require.NoError(t, err, "set done error")
	_, err = s.SetDone(ctx, []int64{task.ID}, false)
	require.NoError(t, err, "undo done error")

	_, err = s.DeleteTasks(ctx, []int64{task.ID})
	require.NoError(t, err, "delete task error")

	events, err := s.ListTaskEvents(context.Background(), task.ID, 20)
	require.NoError(t, err, "list task events error")
	require.Len(t, events, 5, "task events count mismatch")

	assert.Equal(t, TaskEventKindDelete, events[0].Kind, "latest event kind mismatch")
	assert.Equal(t, TaskEventKindUndo, events[1].Kind, "undo event kind mismatch")
	assert.Equal(t, TaskEventKindDone, events[2].Kind, "done event kind mismatch")
	assert.Equal(t, TaskEventKindUpdate, events[3].Kind, "update event kind mismatch")
	assert.Equal(t, TaskEventKindCreate, events[4].Kind, "create event kind mismatch")

	assert.Equal(t, "test", events[0].Origin, "origin should be retained")

	var updateChanges map[string]any
	err = json.Unmarshal([]byte(events[3].ChangesJSON), &updateChanges)
	require.NoError(t, err, "unmarshal update changes error")

	notesChange, ok := updateChanges["notes"].(map[string]any)
	require.True(t, ok, "notes change should be present")
	assert.Equal(t, "note one", notesChange["from"], "notes from mismatch")
	assert.Equal(t, "note two", notesChange["to"], "notes to mismatch")

	projectsChange, ok := updateChanges["projects"].(map[string]any)
	require.True(t, ok, "projects change should be present")
	added, ok := projectsChange["added"].([]any)
	require.True(t, ok, "projects added should be present")
	require.Len(t, added, 1, "projects added count mismatch")
	assert.Equal(t, "home", added[0], "projects added value mismatch")
}

func TestTaskEventsIncludeShellHistoryLink(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := openTestStore(t)

	history, err := s.RecordShellHistory(ctx, "add linked task", false, "", "")
	require.NoError(t, err, "record shell history error")

	linkedCtx := WithAuditOrigin(ctx, "shell")
	linkedCtx = WithAuditShellHistoryID(linkedCtx, history.ID)

	task, err := s.CreateTask(linkedCtx, &Task{Title: "linked"})
	require.NoError(t, err, "create task error")

	events, err := s.ListTaskEvents(ctx, task.ID, 5)
	require.NoError(t, err, "list task events error")
	require.Len(t, events, 1, "task event count mismatch")

	require.NotNil(t, events[0].ShellHistoryID, "shell history id should be present")
	assert.Equal(t, history.ID, *events[0].ShellHistoryID, "shell history id mismatch")
	assert.Equal(t, "add linked task", events[0].ShellCommand, "shell command mismatch")
	assert.Equal(t, "shell", events[0].Origin, "origin mismatch")
}
