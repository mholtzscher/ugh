-- name: InsertTask :one
INSERT INTO tasks (
  state,
  prev_state,
  priority,
  title,
  notes,
  due_on,
  waiting_for,
  completed_at,
  created_at,
  updated_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING
  id,
  state,
  prev_state,
  priority,
  CAST(title AS TEXT) AS title,
  CAST(notes AS TEXT) AS notes,
  due_on,
  waiting_for,
  completed_at,
  created_at,
  updated_at;

-- name: UpdateTask :one
UPDATE tasks
SET state = ?,
  prev_state = ?,
  priority = ?,
  title = ?,
  notes = ?,
  due_on = ?,
  waiting_for = ?,
  completed_at = ?,
  updated_at = ?
WHERE id = ?
RETURNING
  id,
  state,
  prev_state,
  priority,
  CAST(title AS TEXT) AS title,
  CAST(notes AS TEXT) AS notes,
  due_on,
  waiting_for,
  completed_at,
  created_at,
  updated_at;

-- name: GetTask :one
SELECT
  id,
  state,
  prev_state,
  priority,
  CAST(title AS TEXT) AS title,
  CAST(notes AS TEXT) AS notes,
  due_on,
  waiting_for,
  completed_at,
  created_at,
  updated_at
FROM tasks
WHERE id = ?;

-- name: ListTasks :many
SELECT
  t.id,
  t.state,
  t.prev_state,
  t.priority,
  CAST(t.title AS TEXT) AS title,
  CAST(t.notes AS TEXT) AS notes,
  t.due_on,
  t.waiting_for,
  t.completed_at,
  t.created_at,
  t.updated_at
FROM tasks t
WHERE (? = 0 OR t.state != 'done')
  AND (? IS NULL OR t.state = ?)
  AND (? IS NULL OR EXISTS (
    SELECT 1
    FROM task_project_links tpl
    JOIN projects p ON p.id = tpl.project_id
    WHERE tpl.task_id = t.id AND p.name = ?
  ))
  AND (? IS NULL OR EXISTS (
    SELECT 1
    FROM task_context_links tcl
    JOIN contexts c ON c.id = tcl.context_id
    WHERE tcl.task_id = t.id AND c.name = ?
  ))
  AND (? IS NULL OR t.priority = ?)
  AND (? IS NULL OR (
    t.title LIKE '%' || ? || '%'
    OR t.notes LIKE '%' || ? || '%'
    OR EXISTS (
      SELECT 1
      FROM task_project_links tpl
      JOIN projects p ON p.id = tpl.project_id
      WHERE tpl.task_id = t.id AND p.name LIKE '%' || ? || '%'
    )
    OR EXISTS (
      SELECT 1
      FROM task_context_links tcl
      JOIN contexts c ON c.id = tcl.context_id
      WHERE tcl.task_id = t.id AND c.name LIKE '%' || ? || '%'
    )
    OR EXISTS (
      SELECT 1
      FROM task_meta m
      WHERE m.task_id = t.id AND (
        m.key LIKE '%' || ? || '%'
        OR m.value LIKE '%' || ? || '%'
      )
    )
  ))
  AND (? = 0 OR (t.due_on IS NOT NULL AND t.due_on != ''))
ORDER BY
  CASE WHEN t.state = 'done' THEN 1 ELSE 0 END,
  CASE WHEN t.due_on IS NULL OR t.due_on = '' THEN 1 ELSE 0 END,
  t.due_on ASC,
  t.updated_at DESC;

-- name: CompleteTasks :execrows
UPDATE tasks
SET
  prev_state = CASE WHEN state != 'done' THEN state ELSE prev_state END,
  state = 'done',
  completed_at = ?,
  updated_at = ?
WHERE id IN (sqlc.slice('ids'));

-- name: ReopenTasks :execrows
UPDATE tasks
SET
  state = COALESCE(prev_state, 'inbox'),
  prev_state = NULL,
  completed_at = NULL,
  updated_at = ?
WHERE id IN (sqlc.slice('ids'));

-- name: DeleteTasks :execrows
DELETE FROM tasks
WHERE id IN (sqlc.slice('ids'));

-- name: EnsureProject :one
INSERT INTO projects (
  name,
  notes,
  created_at,
  updated_at
) VALUES (
  ?, '', ?, ?
)
ON CONFLICT(name) DO UPDATE SET
  updated_at = excluded.updated_at
RETURNING id;

-- name: EnsureContext :one
INSERT INTO contexts (
  name,
  created_at,
  updated_at
) VALUES (
  ?, ?, ?
)
ON CONFLICT(name) DO UPDATE SET
  updated_at = excluded.updated_at
RETURNING id;

-- name: InsertTaskProjectLink :exec
INSERT INTO task_project_links (task_id, project_id) VALUES (?, ?);

-- name: InsertTaskContextLink :exec
INSERT INTO task_context_links (task_id, context_id) VALUES (?, ?);

-- name: InsertMeta :exec
INSERT INTO task_meta (task_id, key, value) VALUES (?, ?, ?);

-- name: DeleteTaskProjectLinks :exec
DELETE FROM task_project_links WHERE task_id = ?;

-- name: DeleteTaskContextLinks :exec
DELETE FROM task_context_links WHERE task_id = ?;

-- name: DeleteMeta :exec
DELETE FROM task_meta WHERE task_id = ?;

-- name: ListProjectsForTasks :many
SELECT tpl.task_id, p.name
FROM task_project_links tpl
JOIN projects p ON p.id = tpl.project_id
WHERE tpl.task_id IN (sqlc.slice('ids'))
ORDER BY tpl.task_id, p.name;

-- name: ListContextsForTasks :many
SELECT tcl.task_id, c.name
FROM task_context_links tcl
JOIN contexts c ON c.id = tcl.context_id
WHERE tcl.task_id IN (sqlc.slice('ids'))
ORDER BY tcl.task_id, c.name;

-- name: ListMeta :many
SELECT task_id, key, value
FROM task_meta
WHERE task_id IN (sqlc.slice('ids'))
ORDER BY task_id;

-- name: ListProjectCounts :many
SELECT p.name, COUNT(t.id) AS count
FROM projects p
JOIN task_project_links tpl ON tpl.project_id = p.id
JOIN tasks t ON t.id = tpl.task_id
WHERE (? = 0 OR t.state = 'done')
  AND (? = 0 OR t.state != 'done')
GROUP BY p.name
ORDER BY p.name ASC;

-- name: ListContextCounts :many
SELECT c.name, COUNT(t.id) AS count
FROM contexts c
JOIN task_context_links tcl ON tcl.context_id = c.id
JOIN tasks t ON t.id = tcl.task_id
WHERE (? = 0 OR t.state = 'done')
  AND (? = 0 OR t.state != 'done')
GROUP BY c.name
ORDER BY c.name ASC;
