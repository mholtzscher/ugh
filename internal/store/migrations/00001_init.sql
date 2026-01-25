-- +goose Up
CREATE TABLE tasks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  done INTEGER NOT NULL DEFAULT 0,
  priority TEXT,
  completion_date TEXT,
  creation_date TEXT,
  description TEXT NOT NULL DEFAULT "",
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);

CREATE TABLE task_projects (
  task_id INTEGER NOT NULL,
  name TEXT NOT NULL,
  FOREIGN KEY(task_id) REFERENCES tasks(id) ON DELETE CASCADE
);

CREATE TABLE task_contexts (
  task_id INTEGER NOT NULL,
  name TEXT NOT NULL,
  FOREIGN KEY(task_id) REFERENCES tasks(id) ON DELETE CASCADE
);

CREATE TABLE task_meta (
  task_id INTEGER NOT NULL,
  key TEXT NOT NULL,
  value TEXT NOT NULL,
  FOREIGN KEY(task_id) REFERENCES tasks(id) ON DELETE CASCADE
);

CREATE TABLE task_unknown (
  task_id INTEGER NOT NULL,
  ordinal INTEGER NOT NULL,
  token TEXT NOT NULL,
  FOREIGN KEY(task_id) REFERENCES tasks(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE task_unknown;
DROP TABLE task_meta;
DROP TABLE task_contexts;
DROP TABLE task_projects;
DROP TABLE tasks;
