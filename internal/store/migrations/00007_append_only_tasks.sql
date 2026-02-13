-- +goose Up

PRAGMA foreign_keys=OFF;

ALTER TABLE tasks RENAME TO tasks_legacy;

CREATE TABLE tasks (
  id INTEGER PRIMARY KEY,
  created_at INTEGER NOT NULL
);

CREATE TABLE task_versions (
  version_id INTEGER PRIMARY KEY AUTOINCREMENT,
  task_id INTEGER NOT NULL,
  state TEXT NOT NULL,
  prev_state TEXT,
  title TEXT NOT NULL DEFAULT '',
  notes TEXT NOT NULL DEFAULT '',
  due_on TEXT,
  waiting_for TEXT,
  completed_at INTEGER,
  updated_at INTEGER NOT NULL,
  deleted INTEGER NOT NULL DEFAULT 0,
  projects_json TEXT NOT NULL DEFAULT '[]',
  contexts_json TEXT NOT NULL DEFAULT '[]',
  meta_json TEXT NOT NULL DEFAULT '{}',
  FOREIGN KEY(task_id) REFERENCES tasks(id) ON DELETE CASCADE,
  CHECK (json_valid(projects_json)),
  CHECK (json_valid(contexts_json)),
  CHECK (json_valid(meta_json))
);

INSERT INTO tasks (id, created_at)
SELECT id, created_at
FROM tasks_legacy;

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
)
SELECT
  t.id,
  t.state,
  t.prev_state,
  t.title,
  t.notes,
  t.due_on,
  t.waiting_for,
  t.completed_at,
  t.updated_at,
  0,
  COALESCE(
    (
      SELECT json_group_array(p.name)
      FROM task_project_links tpl
      JOIN projects p ON p.id = tpl.project_id
      WHERE tpl.task_id = t.id
      ORDER BY p.name
    ),
    '[]'
  ),
  COALESCE(
    (
      SELECT json_group_array(c.name)
      FROM task_context_links tcl
      JOIN contexts c ON c.id = tcl.context_id
      WHERE tcl.task_id = t.id
      ORDER BY c.name
    ),
    '[]'
  ),
  COALESCE(
    (
      SELECT json_group_object(m.key, m.value)
      FROM task_meta m
      WHERE m.task_id = t.id
    ),
    '{}'
  )
FROM tasks_legacy t;

CREATE INDEX idx_task_versions_task_latest ON task_versions(task_id, version_id DESC);
CREATE INDEX idx_task_versions_state ON task_versions(state);
CREATE INDEX idx_task_versions_due_on ON task_versions(due_on);

CREATE TABLE tasks_current (
  id INTEGER PRIMARY KEY,
  state TEXT NOT NULL,
  prev_state TEXT,
  title TEXT NOT NULL DEFAULT '',
  notes TEXT NOT NULL DEFAULT '',
  due_on TEXT,
  waiting_for TEXT,
  completed_at INTEGER,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL,
  projects_json TEXT NOT NULL DEFAULT '[]',
  contexts_json TEXT NOT NULL DEFAULT '[]',
  meta_json TEXT NOT NULL DEFAULT '{}',
  version_id INTEGER NOT NULL,
  CHECK (json_valid(projects_json)),
  CHECK (json_valid(contexts_json)),
  CHECK (json_valid(meta_json)),
  FOREIGN KEY(id) REFERENCES tasks(id) ON DELETE CASCADE
);

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
)
SELECT
  t.id,
  tv.state,
  tv.prev_state,
  tv.title,
  tv.notes,
  tv.due_on,
  tv.waiting_for,
  tv.completed_at,
  t.created_at,
  tv.updated_at,
  tv.projects_json,
  tv.contexts_json,
  tv.meta_json,
  tv.version_id
FROM tasks t
JOIN task_versions tv ON tv.task_id = t.id
WHERE tv.deleted = 0;

DROP TABLE task_meta;
DROP TABLE task_context_links;
DROP TABLE task_project_links;
DROP TABLE contexts;
DROP TABLE projects;
DROP TABLE tasks_legacy;

PRAGMA foreign_keys=ON;

-- +goose Down

SELECT 1;
