-- name: InsertTask :one
INSERT INTO tasks (
  done,
  priority,
  completion_date,
  creation_date,
  description,
  created_at,
  updated_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
)
RETURNING id, done, priority, completion_date, creation_date, CAST(description AS TEXT) AS description, created_at, updated_at;

-- name: UpdateTask :one
UPDATE tasks
SET done = ?,
  priority = ?,
  completion_date = ?,
  creation_date = ?,
  description = ?,
  updated_at = ?
WHERE id = ?
RETURNING id, done, priority, completion_date, creation_date, CAST(description AS TEXT) AS description, created_at, updated_at;

-- name: GetTask :one
SELECT id, done, priority, completion_date, creation_date, CAST(description AS TEXT) AS description, created_at, updated_at
FROM tasks
WHERE id = ?;

-- name: ListTasks :many
SELECT t.id, t.done, t.priority, t.completion_date, t.creation_date, CAST(t.description AS TEXT) AS description, t.created_at, t.updated_at
FROM tasks t
WHERE (sqlc.narg('status') IS NULL OR t.done = sqlc.narg('status'))
  AND (sqlc.narg('project') IS NULL OR EXISTS (
    SELECT 1 FROM task_projects p WHERE p.task_id = t.id AND p.name = sqlc.narg('project')
  ))
  AND (sqlc.narg('context') IS NULL OR EXISTS (
    SELECT 1 FROM task_contexts c WHERE c.task_id = t.id AND c.name = sqlc.narg('context')
  ))
  AND (sqlc.narg('priority') IS NULL OR t.priority = sqlc.narg('priority'))
  AND (sqlc.narg('search') IS NULL OR (
    t.description LIKE '%' || sqlc.narg('search') || '%'
    OR EXISTS (SELECT 1 FROM task_projects p WHERE p.task_id = t.id AND p.name LIKE '%' || sqlc.narg('search') || '%')
    OR EXISTS (SELECT 1 FROM task_contexts c WHERE c.task_id = t.id AND c.name LIKE '%' || sqlc.narg('search') || '%')
    OR EXISTS (SELECT 1 FROM task_meta m WHERE m.task_id = t.id AND (m.key LIKE '%' || sqlc.narg('search') || '%' OR m.value LIKE '%' || sqlc.narg('search') || '%'))
  ))
ORDER BY CASE WHEN t.done = 1 THEN 1 ELSE 0 END, t.priority IS NULL, t.priority ASC, t.created_at DESC;

-- name: SetDone :execrows
UPDATE tasks
SET done = ?, completion_date = ?, updated_at = ?
WHERE id IN (sqlc.slice('ids'));

-- name: DeleteTasks :execrows
DELETE FROM tasks
WHERE id IN (sqlc.slice('ids'));

-- name: InsertProject :exec
INSERT INTO task_projects (task_id, name) VALUES (?, ?);

-- name: InsertContext :exec
INSERT INTO task_contexts (task_id, name) VALUES (?, ?);

-- name: InsertMeta :exec
INSERT INTO task_meta (task_id, key, value) VALUES (?, ?, ?);

-- name: InsertUnknown :exec
INSERT INTO task_unknown (task_id, ordinal, token) VALUES (?, ?, ?);

-- name: DeleteTokens :exec
DELETE FROM task_projects WHERE task_id = ?;

-- name: DeleteContexts :exec
DELETE FROM task_contexts WHERE task_id = ?;

-- name: DeleteMeta :exec
DELETE FROM task_meta WHERE task_id = ?;

-- name: DeleteUnknown :exec
DELETE FROM task_unknown WHERE task_id = ?;

-- name: ListProjects :many
SELECT task_id, name FROM task_projects WHERE task_id IN (sqlc.slice('ids')) ORDER BY task_id;

-- name: ListContexts :many
SELECT task_id, name FROM task_contexts WHERE task_id IN (sqlc.slice('ids')) ORDER BY task_id;

-- name: ListMeta :many
SELECT task_id, key, value FROM task_meta WHERE task_id IN (sqlc.slice('ids')) ORDER BY task_id;

-- name: ListUnknown :many
SELECT task_id, ordinal, token FROM task_unknown WHERE task_id IN (sqlc.slice('ids')) ORDER BY task_id, ordinal;

-- name: ListProjectCounts :many
SELECT tp.name, COUNT(t.id) AS count
FROM task_projects tp
JOIN tasks t ON tp.task_id = t.id
WHERE (sqlc.narg('status') IS NULL OR t.done = sqlc.narg('status'))
GROUP BY tp.name
ORDER BY tp.name ASC;

-- name: ListContextCounts :many
SELECT tc.name, COUNT(t.id) AS count
FROM task_contexts tc
JOIN tasks t ON tc.task_id = t.id
WHERE (sqlc.narg('status') IS NULL OR t.done = sqlc.narg('status'))
GROUP BY tc.name
ORDER BY tc.name ASC;
