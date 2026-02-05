-- +goose Up

PRAGMA foreign_keys=OFF;

CREATE TABLE tasks_new (
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

INSERT INTO tasks_new (
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
  priority,
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

CREATE TABLE tasks_legacy (
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
  updated_at INTEGER NOT NULL,
  state TEXT NOT NULL DEFAULT 'inbox',
  prev_state TEXT
);

INSERT INTO tasks_legacy (
  id,
  done,
  status,
  priority,
  title,
  notes,
  due_on,
  defer_until,
  waiting_for,
  completed_at,
  created_at,
  updated_at,
  state,
  prev_state
)
SELECT
  id,
  CASE WHEN state = 'done' THEN 1 ELSE 0 END,
  CASE
    WHEN state = 'now' THEN 'next'
    WHEN state = 'later' THEN 'someday'
    WHEN state = 'done' THEN CASE
      WHEN prev_state = 'now' THEN 'next'
      WHEN prev_state = 'later' THEN 'someday'
      WHEN prev_state IS NULL OR prev_state = '' THEN 'inbox'
      ELSE prev_state
    END
    ELSE state
  END,
  priority,
  title,
  notes,
  due_on,
  NULL,
  waiting_for,
  completed_at,
  created_at,
  updated_at,
  state,
  prev_state
FROM tasks;

DROP TABLE tasks;
ALTER TABLE tasks_legacy RENAME TO tasks;

CREATE INDEX IF NOT EXISTS idx_tasks_done_status ON tasks(done, status);
CREATE INDEX IF NOT EXISTS idx_tasks_due_on ON tasks(due_on);
CREATE INDEX IF NOT EXISTS idx_tasks_defer_until ON tasks(defer_until);

PRAGMA foreign_keys=ON;
