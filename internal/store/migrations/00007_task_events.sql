-- +goose Up

CREATE TABLE task_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id INTEGER NOT NULL,
    timestamp INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    kind TEXT NOT NULL,
    summary TEXT,
    changes_json TEXT,
    origin TEXT,
    shell_history_id INTEGER,
    FOREIGN KEY(shell_history_id) REFERENCES shell_history(id) ON DELETE SET NULL
);

CREATE INDEX idx_task_events_task_time ON task_events(task_id, timestamp DESC, id DESC);

-- +goose Down

DROP TABLE IF EXISTS task_events;
