-- +goose Up

PRAGMA foreign_keys=OFF;

DROP TABLE IF EXISTS tasks_current;
DROP VIEW IF EXISTS tasks_current;
DROP TABLE IF EXISTS task_versions;

DROP TABLE IF EXISTS task_meta;
DROP TABLE IF EXISTS task_context_links;
DROP TABLE IF EXISTS task_project_links;
DROP TABLE IF EXISTS contexts;
DROP TABLE IF EXISTS projects;

DROP TABLE IF EXISTS tasks;

CREATE TABLE tasks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
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

PRAGMA foreign_keys=ON;

-- +goose Down

SELECT 1;
