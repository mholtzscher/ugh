-- name: InsertTaskEvent :one
INSERT INTO task_events (
  task_id,
  timestamp,
  kind,
  summary,
  changes_json,
  origin,
  shell_history_id
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
)
RETURNING
  id,
  task_id,
  timestamp,
  kind,
  summary,
  changes_json,
  origin,
  shell_history_id;

-- name: ListTaskEventsByTaskID :many
SELECT
  te.id,
  te.task_id,
  te.timestamp,
  te.kind,
  te.summary,
  te.changes_json,
  te.origin,
  te.shell_history_id,
  sh.command AS shell_command
FROM task_events te
LEFT JOIN shell_history sh ON sh.id = te.shell_history_id
WHERE te.task_id = ?
ORDER BY te.timestamp DESC, te.id DESC
LIMIT ?;
