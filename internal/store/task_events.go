package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"sort"
	"time"

	"github.com/mholtzscher/ugh/internal/store/sqlc"
)

const defaultTaskEventLimit = 50

func (s *Store) ListTaskEvents(ctx context.Context, taskID int64, limit int64) ([]*TaskEvent, error) {
	if taskID == 0 {
		return nil, errors.New("task id is required")
	}
	if limit <= 0 {
		limit = defaultTaskEventLimit
	}

	rows, err := s.queries.ListTaskEventsByTaskID(ctx, sqlc.ListTaskEventsByTaskIDParams{TaskID: taskID, Limit: limit})
	if err != nil {
		return nil, fmt.Errorf("list task events: %w", err)
	}

	result := make([]*TaskEvent, 0, len(rows))
	for _, row := range rows {
		result = append(result, &TaskEvent{
			ID:             row.ID,
			TaskID:         row.TaskID,
			Timestamp:      row.Timestamp,
			Kind:           row.Kind,
			Summary:        row.Summary.String,
			ChangesJSON:    row.ChangesJson.String,
			Origin:         row.Origin.String,
			ShellHistoryID: parseInt64Ptr(row.ShellHistoryID),
			ShellCommand:   row.ShellCommand.String,
		})
	}

	return result, nil
}

func (s *Store) recordTaskEvent(
	ctx context.Context,
	taskID int64,
	kind string,
	summary string,
	changes map[string]any,
) error {
	if taskID == 0 {
		return nil
	}

	changesJSON, err := encodeTaskChanges(changes)
	if err != nil {
		return err
	}

	origin := auditOriginFromContext(ctx)
	if origin == "" {
		origin = "cli"
	}

	shellHistoryID := auditShellHistoryIDFromContext(ctx)
	params := sqlc.InsertTaskEventParams{
		TaskID:         taskID,
		Timestamp:      time.Now().UTC().Unix(),
		Kind:           kind,
		Summary:        nullString(summary),
		ChangesJson:    nullString(changesJSON),
		Origin:         nullString(origin),
		ShellHistoryID: nullInt64(shellHistoryID),
	}

	_, err = s.queries.InsertTaskEvent(ctx, params)
	if err != nil {
		return fmt.Errorf("insert task event: %w", err)
	}
	return nil
}

func encodeTaskChanges(changes map[string]any) (string, error) {
	if len(changes) == 0 {
		return "", nil
	}
	payload, err := json.Marshal(changes)
	if err != nil {
		return "", fmt.Errorf("marshal task event changes: %w", err)
	}
	return string(payload), nil
}

func buildCreateChanges(after *Task) map[string]any {
	if after == nil {
		return nil
	}
	return map[string]any{"after": taskSnapshot(after)}
}

func buildDeleteChanges(before *Task) map[string]any {
	if before == nil {
		return nil
	}
	return map[string]any{"before": taskSnapshot(before)}
}

func buildTaskDiff(before, after *Task) map[string]any {
	if before == nil || after == nil {
		return nil
	}

	changes := map[string]any{}
	addStringChange(changes, "state", string(before.State), string(after.State))
	addStringChange(changes, "prev_state", statePtrString(before.PrevState), statePtrString(after.PrevState))
	addStringChange(changes, "title", before.Title, after.Title)
	addStringChange(changes, "notes", before.Notes, after.Notes)
	addStringChange(changes, "due_on", datePtrString(before.DueOn), datePtrString(after.DueOn))
	addStringChange(changes, "waiting_for", before.WaitingFor, after.WaitingFor)
	addStringChange(
		changes,
		"completed_at",
		dateTimePtrString(before.CompletedAt),
		dateTimePtrString(after.CompletedAt),
	)

	if projectChanges := listChanges(before.Projects, after.Projects); len(projectChanges) > 0 {
		changes["projects"] = projectChanges
	}
	if contextChanges := listChanges(before.Contexts, after.Contexts); len(contextChanges) > 0 {
		changes["contexts"] = contextChanges
	}
	if metaChanges := metaChanges(before.Meta, after.Meta); len(metaChanges) > 0 {
		changes["meta"] = metaChanges
	}

	return changes
}

func taskSnapshot(task *Task) map[string]any {
	if task == nil {
		return nil
	}

	metaCopy := map[string]string{}
	maps.Copy(metaCopy, task.Meta)

	return map[string]any{
		"id":           task.ID,
		"state":        string(task.State),
		"prev_state":   statePtrString(task.PrevState),
		"title":        task.Title,
		"notes":        task.Notes,
		"due_on":       datePtrString(task.DueOn),
		"waiting_for":  task.WaitingFor,
		"completed_at": dateTimePtrString(task.CompletedAt),
		"projects":     sortedStrings(task.Projects),
		"contexts":     sortedStrings(task.Contexts),
		"meta":         metaCopy,
	}
}

func addStringChange(changes map[string]any, key, before, after string) {
	if before == after {
		return
	}
	changes[key] = map[string]any{"from": before, "to": after}
}

func listChanges(before, after []string) map[string]any {
	beforeSet := make(map[string]bool, len(before))
	afterSet := make(map[string]bool, len(after))
	for _, value := range before {
		beforeSet[value] = true
	}
	for _, value := range after {
		afterSet[value] = true
	}

	added := make([]string, 0)
	for value := range afterSet {
		if !beforeSet[value] {
			added = append(added, value)
		}
	}

	removed := make([]string, 0)
	for value := range beforeSet {
		if !afterSet[value] {
			removed = append(removed, value)
		}
	}

	sort.Strings(added)
	sort.Strings(removed)

	changes := map[string]any{}
	if len(added) > 0 {
		changes["added"] = added
	}
	if len(removed) > 0 {
		changes["removed"] = removed
	}

	return changes
}

func metaChanges(before, after map[string]string) map[string]any {
	changes := map[string]any{}

	added := map[string]string{}
	removed := make([]string, 0)
	updated := map[string]any{}

	for key, value := range after {
		beforeValue, ok := before[key]
		if !ok {
			added[key] = value
			continue
		}
		if beforeValue != value {
			updated[key] = map[string]any{"from": beforeValue, "to": value}
		}
	}

	for key := range before {
		if _, ok := after[key]; !ok {
			removed = append(removed, key)
		}
	}

	sort.Strings(removed)

	if len(added) > 0 {
		changes["added"] = added
	}
	if len(updated) > 0 {
		changes["updated"] = updated
	}
	if len(removed) > 0 {
		changes["removed"] = removed
	}

	return changes
}

func sortedStrings(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}
	result := append([]string(nil), values...)
	sort.Strings(result)
	return result
}

func statePtrString(value *State) string {
	if value == nil {
		return ""
	}
	return string(*value)
}

func datePtrString(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.UTC().Format("2006-01-02")
}

func dateTimePtrString(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.UTC().Format(time.RFC3339)
}

func parseInt64Ptr(value sql.NullInt64) *int64 {
	if !value.Valid {
		return nil
	}
	result := value.Int64
	return &result
}

func nullInt64(value *int64) sql.NullInt64 {
	if value == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *value, Valid: true}
}

func (s *Store) getExistingTasksByID(ctx context.Context, ids []int64) (map[int64]*Task, error) {
	result := make(map[int64]*Task, len(ids))
	for _, id := range ids {
		task, err := s.GetTask(ctx, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			return nil, fmt.Errorf("get task %d: %w", id, err)
		}
		result[id] = task
	}
	return result, nil
}
