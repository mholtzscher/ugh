-- +goose Up

PRAGMA foreign_keys=OFF;

CREATE TABLE tasks_new (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  state TEXT NOT NULL DEFAULT 'inbox',
  prev_state TEXT,
  title TEXT NOT NULL DEFAULT '',
  notes TEXT NOT NULL DEFAULT '',
  due_on TEXT,
  waiting_for TEXT,
  completed_at INTEGER,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);

INSERT INTO tasks_new (
  id,
  state,
  prev_state,
  title,
  notes,
  due_on,
  waiting_for,
  completed_at,
  created_at,
  updated_at
)
SELECT
  id,
  state,
  prev_state,
  title,
  notes,
  due_on,
  waiting_for,
  completed_at,
  created_at,
  updated_at
FROM tasks;

DROP TABLE tasks;
ALTER TABLE tasks_new RENAME TO tasks;

CREATE INDEX IF NOT EXISTS idx_tasks_state ON tasks(state);
CREATE INDEX IF NOT EXISTS idx_tasks_due_on ON tasks(due_on);

PRAGMA foreign_keys=ON;

-- +goose Down

PRAGMA foreign_keys=OFF;

CREATE TABLE tasks_with_priority (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  state TEXT NOT NULL DEFAULT 'inbox',
  prev_state TEXT,
  priority TEXT,
  title TEXT NOT NULL DEFAULT '',
  notes TEXT NOT NULL DEFAULT '',
  due_on TEXT,
  waiting_for TEXT,
  completed_at INTEGER,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);

INSERT INTO tasks_with_priority (
  id,
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
)
SELECT
  id,
  state,
  prev_state,
  NULL,
  title,
  notes,
  due_on,
  waiting_for,
  completed_at,
  created_at,
  updated_at
FROM tasks;

DROP TABLE tasks;
ALTER TABLE tasks_with_priority RENAME TO tasks;

CREATE INDEX IF NOT EXISTS idx_tasks_state ON tasks(state);
CREATE INDEX IF NOT EXISTS idx_tasks_due_on ON tasks(due_on);

PRAGMA foreign_keys=ON;
