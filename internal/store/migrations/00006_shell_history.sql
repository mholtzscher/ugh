-- +goose Up

CREATE TABLE shell_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    command TEXT NOT NULL,
    success BOOLEAN NOT NULL DEFAULT 0,
    result_summary TEXT,
    intent TEXT
);

CREATE INDEX idx_shell_history_timestamp ON shell_history(timestamp DESC);

-- +goose Down

DROP TABLE IF EXISTS shell_history;
