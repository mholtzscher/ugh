-- +goose Up

-- This migration intentionally drops the existing legacy schema.
-- It is a clean break to a GTD-first model.

DROP TABLE IF EXISTS task_unknown;
DROP TABLE IF EXISTS task_meta;
DROP TABLE IF EXISTS task_contexts;
DROP TABLE IF EXISTS task_projects;
DROP TABLE IF EXISTS tasks;

DROP TABLE IF EXISTS task_context_links;
DROP TABLE IF EXISTS task_project_links;
DROP TABLE IF EXISTS contexts;
DROP TABLE IF EXISTS projects;

CREATE TABLE tasks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  done INTEGER NOT NULL DEFAULT 0,
  status TEXT NOT NULL DEFAULT 'inbox',
  priority TEXT,
  title TEXT NOT NULL DEFAULT '',
  notes TEXT NOT NULL DEFAULT '',
  due_on TEXT,
  defer_until TEXT,
  waiting_for TEXT,
  completed_at INTEGER,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);

CREATE TABLE projects (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE,
  notes TEXT NOT NULL DEFAULT '',
  archived_at INTEGER,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);

CREATE TABLE contexts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);

CREATE TABLE task_project_links (
  task_id INTEGER NOT NULL,
  project_id INTEGER NOT NULL,
  FOREIGN KEY(task_id) REFERENCES tasks(id) ON DELETE CASCADE,
  FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE,
  UNIQUE(task_id, project_id)
);

CREATE TABLE task_context_links (
  task_id INTEGER NOT NULL,
  context_id INTEGER NOT NULL,
  FOREIGN KEY(task_id) REFERENCES tasks(id) ON DELETE CASCADE,
  FOREIGN KEY(context_id) REFERENCES contexts(id) ON DELETE CASCADE,
  UNIQUE(task_id, context_id)
);

CREATE TABLE task_meta (
  task_id INTEGER NOT NULL,
  key TEXT NOT NULL,
  value TEXT NOT NULL,
  FOREIGN KEY(task_id) REFERENCES tasks(id) ON DELETE CASCADE
);

CREATE INDEX idx_tasks_done_status ON tasks(done, status);
CREATE INDEX idx_tasks_due_on ON tasks(due_on);
CREATE INDEX idx_tasks_defer_until ON tasks(defer_until);

CREATE INDEX idx_projects_name ON projects(name);
CREATE INDEX idx_contexts_name ON contexts(name);

CREATE INDEX idx_task_project_links_task ON task_project_links(task_id);
CREATE INDEX idx_task_context_links_task ON task_context_links(task_id);

CREATE INDEX idx_task_meta_task ON task_meta(task_id);

-- +goose Down

DROP TABLE task_meta;
DROP TABLE task_context_links;
DROP TABLE task_project_links;
DROP TABLE contexts;
DROP TABLE projects;
DROP TABLE tasks;
