-- +goose Up

PRAGMA foreign_keys=OFF;

DROP VIEW IF EXISTS tasks_current;

CREATE TABLE IF NOT EXISTS tasks_current (
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

PRAGMA foreign_keys=ON;

-- +goose Down

DROP TABLE IF EXISTS tasks_current;

CREATE VIEW tasks_current AS
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
  tv.deleted,
  tv.version_id
FROM tasks t
JOIN task_versions tv ON tv.task_id = t.id
WHERE tv.version_id = (
  SELECT tv2.version_id
  FROM task_versions tv2
  WHERE tv2.task_id = t.id
  ORDER BY tv2.version_id DESC
  LIMIT 1
);
