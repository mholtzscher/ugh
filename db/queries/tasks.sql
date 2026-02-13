-- name: InsertTaskIdentity :execresult
INSERT INTO tasks (created_at) VALUES (?);

-- name: InsertTaskVersion :one
INSERT INTO task_versions (
  task_id,
  state,
  prev_state,
  title,
  notes,
  due_on,
  waiting_for,
  completed_at,
  updated_at,
  deleted,
  projects_json,
  contexts_json,
  meta_json
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING version_id;

-- name: UpsertTaskCurrent :exec
INSERT INTO tasks_current (
  id,
  state,
  prev_state,
  title,
  notes,
  due_on,
  waiting_for,
  completed_at,
  created_at,
  updated_at,
  projects_json,
  contexts_json,
  meta_json,
  version_id
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
)
ON CONFLICT(id) DO UPDATE SET
  state = excluded.state,
  prev_state = excluded.prev_state,
  title = excluded.title,
  notes = excluded.notes,
  due_on = excluded.due_on,
  waiting_for = excluded.waiting_for,
  completed_at = excluded.completed_at,
  created_at = excluded.created_at,
  updated_at = excluded.updated_at,
  projects_json = excluded.projects_json,
  contexts_json = excluded.contexts_json,
  meta_json = excluded.meta_json,
  version_id = excluded.version_id;

-- name: DeleteTaskCurrent :exec
DELETE FROM tasks_current
WHERE id = ?;

-- name: GetTask :one
SELECT
  id,
  state,
  prev_state,
  CAST(title AS TEXT) AS title,
  CAST(notes AS TEXT) AS notes,
  due_on,
  waiting_for,
  completed_at,
  created_at,
  updated_at,
  projects_json,
  contexts_json,
  meta_json,
  version_id
FROM tasks_current
WHERE id = ?;

-- name: ListTaskVersions :many
SELECT
  version_id,
  task_id,
  state,
  prev_state,
  title,
  notes,
  due_on,
  waiting_for,
  completed_at,
  updated_at,
  deleted,
  projects_json,
  contexts_json,
  meta_json
FROM task_versions
WHERE task_id = ?
ORDER BY version_id DESC
LIMIT ?;
